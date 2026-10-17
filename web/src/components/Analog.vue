<template>
  <div ref="root" class="analog" :style="{ opacity: displayOpacity(brightness) }">
    <svg
      v-if="width > 0 && height > 0"
      class="analog-svg"
      :viewBox="`0 0 ${width} ${height}`"
      preserveAspectRatio="xMidYMid meet"
    >
      <g :transform="`translate(${centerX} ${centerY}) scale(${dialScale}) translate(${-centerX} ${-centerY})`">
      <line
        v-for="(tick, index) in ticks"
        :key="index"
        :x1="tick.x1"
        :y1="tick.y1"
        :x2="tick.x2"
        :y2="tick.y2"
        :stroke-width="tickWidth"
        :stroke-opacity="tick.opacity"
        class="tick"
      />

      <g :transform="`translate(${centerX} ${centerY})`">
        <g
          v-for="hand in hands"
          :key="hand.name"
          :transform="`rotate(${hand.angle})`"
        >
          <line
            class="hand"
            stroke-linecap="round"
            :y1="handShape.neckY1"
            :y2="handShape.neckY2"
            :stroke-width="neckWidth"
          />
          <path class="hand-fill" :d="handShape.taper" />
          <line
            class="hand"
            stroke-linecap="butt"
            :y1="handShape.bodyY1"
            :y2="hand.shaftEnd"
            :stroke-width="handWidth"
          />
          <circle class="hand-fill" :cy="hand.tipCenter" :r="handWidth / 2" />
        </g>
      </g>

      <text
        v-for="mark in numerals"
        :key="mark.label"
        :x="mark.x"
        :y="mark.y"
        :font-size="numeralSize"
        class="numeral"
        text-anchor="middle"
        dominant-baseline="central"
      >{{ mark.label }}</text>

      <text
        ref="dateEl"
        class="date"
        :x="dateX"
        :y="centerY"
        :dy="dateDy"
        :font-size="dateSize"
        text-anchor="start"
        dominant-baseline="central"
      >
        <tspan class="day">{{ dayOfMonth }}</tspan>
        <tspan class="weekday" :class="{ night: isNightMode }" dx="0.35em">周{{ weekday }}</tspan>
      </text>

      <g :transform="`translate(${centerX} ${centerY})`">
        <line
          class="second"
          :y1="secondTail"
          :y2="-secondLength"
          :stroke-width="secondWidth"
          stroke-linecap="round"
          :transform="`rotate(${secondAngle})`"
        />
        <circle class="hub" :r="hubRadius" />
        <circle class="cap" :r="capRadius" />
        <circle class="pin" :r="pinRadius" />
      </g>
      </g>
    </svg>
  </div>
</template>

<script lang="ts">
// 每个主题 6 个颜色，下标与颜色面板、ClockLayout 一致：
// [0] 色点左半  [1] 色点右半
// [2] 刻度、数字、日期、时针、分针
// [3] 秒针、中心圆点
// [4][5] 预留，Analog 暂未使用
export const analogColors = [
  ['#f4f4f5', '#e39b4e', '#ffffff', '#e39b4e', '#ffffff', '#e39b4e'],
  ['#5aa6e6', '#9fd4f8', '#8ec8f8', '#e39b4e', '#8ec8f8', '#e39b4e'],
  ['#b794f6', '#d6bcfa', '#c4b5fd', '#e39b4e', '#c4b5fd', '#e39b4e'],
  ['#3dd68c', '#9ae6b4', '#7dcea0', '#e39b4e', '#7dcea0', '#e39b4e'],
  ['#f6d365', '#ffe9a0', '#ffe08a', '#e39b4e', '#ffe08a', '#e39b4e'],
  ['#ff7a9c', '#ffc1d1', '#ff9bb3', '#e39b4e', '#ff9bb3', '#e39b4e'],
  ['#5ce1e6', '#b6f3f5', '#8eebef', '#e39b4e', '#8eebef', '#e39b4e'],
  ['#ff9f43', '#ffd0a8', '#ffb15a', '#ffffff', '#ffb15a', '#ffffff'],
]

export const analogNightModeColors = ['rgb(148,4,7)', 'rgb(111,26,23)']
</script>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

defineProps<{
  digits: string[]
  displayOpacity: (opacity: number) => number
  colonOpacity: number
  brightness: number
  isNightMode: boolean
}>()

// 参考图刻度内缩。整表按这个比例绘制，再等比放大，让左右刻度贴到屏幕边缘。
const MARGIN_X = 0.0527
const MARGIN_Y = 0.0551
const HOUR_INNER = 0.56
const TICK_INSET = 0.0614
const HOUR_TICK_OPACITY = 0.66
const MINUTE_TICK_OPACITY = 0.315
const NUM_RX = 0.627
const NUM_RY = 0.563
// 时针、分针长度为中心到表盘上沿的比例，不随角度变化。秒针竖直向上时针尖落在表盘上沿。
const HOUR_LENGTH = 0.51
const MINUTE_LENGTH = 0.89
const SECOND_TAIL = 0.083
const HAND_WIDTH = 0.034
const SECOND_WIDTH = 0.0066
// 时针、分针细颈从圆心伸出，再过渡到全宽。比例相对短边。
const NECK_WIDTH = 0.0143
const NECK_END = 0.058
const BODY_START = 0.072
// 三针原点固定在圆心。白圆盖住交点，上面是金色圆盘和黑点。秒针从下面穿过。
const HUB_RADIUS = 29.5 / 1290
const CAP_RADIUS = 0.0163
const PIN_RADIUS = 0.007
const TICK_WIDTH = 0.0076
const NUMERAL_SIZE = 0.148
// 数字在画面上的字号。电脑上表盘短边约 885px 时是 78px，按这个比例随表盘缩放。
const NUMERAL_FONT = 87 / 885
const DATE_SIZE = 0.048
const DATE_X = 0.336
const CARDINAL_GAP = 0.05
// 12 上移、6 下移。按参考高度 472px 约 5px，随表盘缩放。
const NUMERAL_SHIFT = 5 / 472

const root = ref<HTMLElement | null>(null)
const dateEl = ref<SVGTextElement | null>(null)
const dateDy = ref(0)
const width = ref(0)
const height = ref(0)
const now = ref(new Date())

let resizeObserver: ResizeObserver | undefined
let frame = 0

const centerX = computed(() => width.value / 2)
const centerY = computed(() => height.value / 2)
const dialScale = computed(() => 0.5 / (0.5 - MARGIN_X))
const shortSide = computed(() => Math.min(width.value, height.value))

const secondTail = computed(() => shortSide.value * SECOND_TAIL)
const handWidth = computed(() => shortSide.value * HAND_WIDTH)
const neckWidth = computed(() => shortSide.value * NECK_WIDTH)
const secondWidth = computed(() => Math.max(1.5, shortSide.value * SECOND_WIDTH))
const hubRadius = computed(() => shortSide.value * HUB_RADIUS)
const capRadius = computed(() => shortSide.value * CAP_RADIUS)
const pinRadius = computed(() => shortSide.value * PIN_RADIUS)

const handShape = computed(() => {
  const nr = neckWidth.value / 2
  const br = handWidth.value / 2
  const neckEnd = shortSide.value * NECK_END
  const bodyStart = shortSide.value * BODY_START
  const y0 = -neckEnd
  const y1 = -bodyStart
  const y2 = y1 - br * 0.45
  const du = br - nr
  const dv = y0 - y1
  // 细颈到全宽的肩部是圆弧，并与全宽侧边相切。
  const radius = (du * du + dv * dv) / (2 * du)
  const cx = br - radius
  const cy = y1
  const a0 = Math.atan2(y0 - cy, nr - cx)
  const a1 = Math.atan2(0, radius)
  let delta = a1 - a0
  if (delta > Math.PI) delta -= Math.PI * 2
  if (delta < -Math.PI) delta += Math.PI * 2
  const steps = 24
  const right: string[] = []
  for (let i = 0; i <= steps; i++) {
    const angle = a0 + (delta * i) / steps
    right.push(`${cx + radius * Math.cos(angle)},${cy + radius * Math.sin(angle)}`)
  }
  const left = right.map((point) => {
    const [x, y] = point.split(',')
    return `${-Number(x)},${y}`
  }).reverse()
  return {
    neckY1: 0,
    neckY2: y0,
    taper: `M ${-nr},${y0} L ${right.join(' ')} L ${br},${y2} L ${-br},${y2} L ${left.join(' ')} Z`,
    bodyY1: y1,
  }
})

const hands = computed(() => {
  const br = handWidth.value / 2
  return [
    { name: 'hour', angle: hourAngle.value, length: hourLength.value },
    { name: 'minute', angle: minuteAngle.value, length: minuteLength.value },
  ].map((hand) => {
    const visualTip = hand.length + br
    const tipCenter = -(visualTip - br)
    return {
      name: hand.name,
      angle: hand.angle,
      tipCenter,
      shaftEnd: tipCenter - br * 0.35,
    }
  })
})
const tickWidth = computed(() => shortSide.value * TICK_WIDTH)
const numeralSize = computed(() => shortSide.value * NUMERAL_FONT)
const dateSize = computed(() => shortSide.value * DATE_SIZE)

const hourAngle = computed(() => {
  const date = now.value
  const hours = date.getHours() % 12
  const minutes = date.getMinutes()
  const seconds = date.getSeconds() + date.getMilliseconds() / 1000
  return (hours + minutes / 60 + seconds / 3600) * 30
})

const minuteAngle = computed(() => {
  const date = now.value
  const seconds = date.getSeconds() + date.getMilliseconds() / 1000
  return (date.getMinutes() + seconds / 60) * 6
})

const secondAngle = computed(() => {
  const date = now.value
  return (date.getSeconds() + date.getMilliseconds() / 1000) * 6
})

const dayOfMonth = computed(() => now.value.getDate())
const weekday = computed(() => ['日', '一', '二', '三', '四', '五', '六'][now.value.getDay()])

const dateX = computed(() => centerX.value + (width.value / 2) * DATE_X)

interface Tick {
  x1: number
  y1: number
  x2: number
  y2: number
  opacity: number
}

function rayLength(angleDeg: number, halfW: number, halfH: number): number {
  const rad = (angleDeg * Math.PI) / 180
  const dx = Math.abs(Math.sin(rad))
  const dy = Math.abs(Math.cos(rad))
  const tx = dx < 1e-6 ? Infinity : halfW / dx
  const ty = dy < 1e-6 ? Infinity : halfH / dy
  return Math.min(tx, ty)
}

const hourLength = computed(() => layout.value.halfH * HOUR_LENGTH)
const minuteLength = computed(() => layout.value.halfH * MINUTE_LENGTH)

const secondLength = computed(() => Math.max(0, layout.value.halfH - secondWidth.value / 2))

const layout = computed(() => {
  const w = width.value
  const h = height.value
  const s = shortSide.value
  const cx = w / 2
  const cy = h / 2
  const halfW = cx - w * MARGIN_X
  const halfH = cy - h * MARGIN_Y
  const inset = s * TICK_INSET
  const numRx = cx * NUM_RX
  const numRy = cy * NUM_RY
  const numeral = s * NUMERAL_SIZE
  return { w, h, s, cx, cy, halfW, halfH, inset, numRx, numRy, numeral }
})

const ticks = computed<Tick[]>(() => {
  const { cx, cy, halfW, halfH, inset, numRx, numRy, numeral, s } = layout.value
  if (halfW <= inset || halfH <= inset) return []

  const glyphHalfH = numeral * 0.34
  const glyphHalfW = numeral * 0.36
  const gap = s * CARDINAL_GAP
  const points: Tick[] = []

  for (let i = 0; i < 60; i++) {
    const angle = i * 6
    const rad = (angle * Math.PI) / 180
    const dx = Math.sin(rad)
    const dy = -Math.cos(rad)
    const outer = rayLength(angle, halfW, halfH)
    let inner: number
    if (i % 15 === 0) {
      inner = i % 30 === 0 ? numRy + glyphHalfH + gap : numRx + glyphHalfW + gap
    } else if (i % 5 === 0) {
      inner = outer * HOUR_INNER
    } else {
      inner = rayLength(angle, halfW - inset, halfH - inset)
    }
    if (!(inner < outer)) continue
    points.push({
      x1: cx + dx * inner,
      y1: cy + dy * inner,
      x2: cx + dx * outer,
      y2: cy + dy * outer,
      opacity: i % 5 === 0 ? HOUR_TICK_OPACITY : MINUTE_TICK_OPACITY,
    })
  }
  return points
})

const numerals = computed(() => {
  const { cx, cy, numRx, numRy, s } = layout.value
  const shift = s * NUMERAL_SHIFT
  return [
    { label: '12', x: cx, y: cy - numRy - shift },
    { label: '3', x: cx + numRx, y: cy },
    { label: '6', x: cx, y: cy + numRy + shift },
    { label: '9', x: cx - numRx, y: cy },
  ]
})

const syncSize = () => {
  const el = root.value
  if (!el) return
  width.value = el.clientWidth
  height.value = el.clientHeight
}

// 手机浏览器常忽略 dominant-baseline，日期会贴着基线偏上。按字形包围盒拉回中线。
const alignDate = () => {
  const el = dateEl.value
  if (!el || !centerY.value) return
  const box = el.getBBox()
  if (box.height <= 0) return
  const inkCenter = box.y + box.height / 2
  const next = dateDy.value + (centerY.value - inkCenter)
  if (Math.abs(next - dateDy.value) > 0.2) dateDy.value = next
}

watch([width, height, dayOfMonth, weekday, dateSize], () => {
  nextTick(alignDate)
})

const tick = () => {
  now.value = new Date()
  frame = window.requestAnimationFrame(tick)
}

const onVisibility = () => {
  if (document.hidden) {
    if (frame) cancelAnimationFrame(frame)
    frame = 0
    return
  }
  if (!frame) frame = window.requestAnimationFrame(tick)
}

onMounted(() => {
  syncSize()
  if (root.value) {
    resizeObserver = new ResizeObserver(syncSize)
    resizeObserver.observe(root.value)
  }
  frame = window.requestAnimationFrame(tick)
  document.addEventListener('visibilitychange', onVisibility)
})

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
  if (frame) cancelAnimationFrame(frame)
  document.removeEventListener('visibilitychange', onVisibility)
})
</script>

<style scoped>
.analog {
  width: 100%;
  height: auto;
  aspect-ratio: 1024 / 472;
  flex: 0 0 auto;
}

.analog-svg {
  display: block;
  width: 100%;
  height: 100%;
  text-size-adjust: 100%;
  -webkit-text-size-adjust: 100%;
}

.tick {
  stroke: var(--color-0);
  stroke-linecap: butt;
}

.numeral {
  fill: var(--color-0);
  font-family: 'ResidenzGrotesk-SemBdExp', sans-serif;
}

.date {
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Text', 'PingFang SC', 'Helvetica Neue', sans-serif;
  font-weight: 500;
}

.day {
  fill: var(--color-0);
}

.weekday {
  fill: #ff453a;
  font-weight: 600;
}

.weekday.night {
  fill: var(--color-0);
}

.hand {
  stroke: var(--color-0);
  fill: none;
}

.hand-fill {
  fill: var(--color-0);
}

.hub {
  fill: var(--color-0);
}

.second {
  stroke: var(--color-1);
}

.cap {
  fill: var(--color-1);
}

.pin {
  fill: #000;
}
</style>
