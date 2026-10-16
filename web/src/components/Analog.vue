<template>
  <div ref="root" class="analog" :style="{ opacity: displayOpacity(brightness) }">
    <svg
      v-if="width > 0 && height > 0"
      class="analog-svg"
      :viewBox="`0 0 ${width} ${height}`"
      preserveAspectRatio="xMidYMid meet"
    >
      <line
        v-for="(tick, index) in ticks"
        :key="index"
        :x1="tick.x1"
        :y1="tick.y1"
        :x2="tick.x2"
        :y2="tick.y2"
        :stroke-width="tickWidth"
        class="tick"
      />

      <g :transform="`translate(${centerX} ${centerY})`">
        <line
          class="hand"
          y1="0"
          :y2="-hourLength"
          :stroke-width="handWidth"
          stroke-linecap="round"
          :transform="`rotate(${hourAngle})`"
        />
        <line
          class="hand"
          y1="0"
          :y2="-minuteLength"
          :stroke-width="handWidth"
          stroke-linecap="round"
          :transform="`rotate(${minuteAngle})`"
        />
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
        class="date"
        :x="dateX"
        :y="centerY"
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
        <circle class="cap" :r="capRadius" />
      </g>
    </svg>
  </div>
</template>

<script lang="ts">
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
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

defineProps<{
  digits: string[]
  displayOpacity: (opacity: number) => number
  colonOpacity: number
  brightness: number
  isNightMode: boolean
}>()

// 比例来自参考截图（1024×472）：刻度外端贴着内缩矩形，时针刻度从 56% 半径拉到边缘。
const MARGIN_X = 0.0527
const MARGIN_Y = 0.0551
const HOUR_INNER = 0.56
const TICK_INSET = 0.0614
const NUM_RX = 0.627
const NUM_RY = 0.563
// 时针、分针长度跟着数字所在的椭圆走。秒针长度固定，竖直向上时针尖落在表盘上沿。
const HOUR_REACH = 0.58
const MINUTE_REACH = 0.76
const SECOND_TAIL = 0.083
const HAND_WIDTH = 0.034
const SECOND_WIDTH = 0.0052
const CAP_RADIUS = 0.018
const TICK_WIDTH = 0.0076
const NUMERAL_SIZE = 0.148
const DATE_SIZE = 0.048
const DATE_X = 0.336
const CARDINAL_GAP = 0.05

const root = ref<HTMLElement | null>(null)
const width = ref(0)
const height = ref(0)
const now = ref(new Date())

let resizeObserver: ResizeObserver | undefined
let frame = 0

const centerX = computed(() => width.value / 2)
const centerY = computed(() => height.value / 2)
const shortSide = computed(() => Math.min(width.value, height.value))

const secondTail = computed(() => shortSide.value * SECOND_TAIL)
const handWidth = computed(() => shortSide.value * HAND_WIDTH)
const secondWidth = computed(() => Math.max(1.5, shortSide.value * SECOND_WIDTH))
const capRadius = computed(() => shortSide.value * CAP_RADIUS)
const tickWidth = computed(() => shortSide.value * TICK_WIDTH)
const numeralSize = computed(() => shortSide.value * NUMERAL_SIZE)
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
}

function ellipseRadius(angleDeg: number, rx: number, ry: number): number {
  const rad = (angleDeg * Math.PI) / 180
  const denom = Math.hypot(ry * Math.sin(rad), rx * Math.cos(rad))
  return denom < 1e-6 ? 0 : (rx * ry) / denom
}

function rayLength(angleDeg: number, halfW: number, halfH: number): number {
  const rad = (angleDeg * Math.PI) / 180
  const dx = Math.abs(Math.sin(rad))
  const dy = Math.abs(Math.cos(rad))
  const tx = dx < 1e-6 ? Infinity : halfW / dx
  const ty = dy < 1e-6 ? Infinity : halfH / dy
  return Math.min(tx, ty)
}

const hourLength = computed(() => {
  const { numRx, numRy } = layout.value
  return ellipseRadius(hourAngle.value, numRx, numRy) * HOUR_REACH
})

const minuteLength = computed(() => {
  const { numRx, numRy } = layout.value
  return ellipseRadius(minuteAngle.value, numRx, numRy) * MINUTE_REACH
})

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
    })
  }
  return points
})

const numerals = computed(() => {
  const { cx, cy, numRx, numRy } = layout.value
  return [
    { label: '12', x: cx, y: cy - numRy },
    { label: '3', x: cx + numRx, y: cy },
    { label: '6', x: cx, y: cy + numRy },
    { label: '9', x: cx - numRx, y: cy },
  ]
})

const syncSize = () => {
  const el = root.value
  if (!el) return
  width.value = el.clientWidth
  height.value = el.clientHeight
}

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
  /*background: #000 url('../assets/analog-reference.png') center / 100% auto no-repeat;*/
}

.analog-svg {
  display: block;
  width: 100%;
  height: 100%;
  /* 半透明叠在参考截图上，方便逐像素对照。对齐后删掉这行和背景图。 */
  opacity: 0.5;
}

.tick {
  stroke: var(--color-0);
  stroke-linecap: butt;
  opacity: 0.64;
}

.numeral {
  fill: var(--color-0);
  font-family: 'SF-Pro-Rounded-Black', sans-serif;
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
}

.second {
  stroke: var(--color-1);
}

.cap {
  fill: var(--color-1);
}
</style>
