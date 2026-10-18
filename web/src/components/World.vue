<template>
  <div
    ref="root"
    class="world"
    :class="{ editing: WORLD_DOT_EDIT }"
    :style="{ opacity: displayOpacity(brightness) }"
    @pointerdown="onWorldPointerDown"
  >
    <svg
      v-if="width > 0 && height > 0"
      class="world-svg"
      :viewBox="`0 0 ${width} ${height}`"
    >
      <circle
        v-for="(dot, index) in dots"
        :key="index"
        :cx="dot.x"
        :cy="dot.y"
        :r="dotRadius"
        :class="dot.off ? 'dot-off' : dot.day ? 'dot-day' : 'dot-night'"
        @click.stop="onDotClick(dot.row, dot.col)"
      />
      <path class="terminator" :d="terminatorPath" :stroke-width="terminatorWidth" />
      <circle
        v-for="(city, index) in cities"
        :key="`city-${index}`"
        class="city"
        :cx="city.x"
        :cy="city.y"
        :r="cityRadius"
      />
    </svg>

    <button v-if="WORLD_DOT_EDIT" type="button" class="dot-edit-log" @click.stop="logWorldDots">
      输出矩阵
    </button>

    <div class="world-readout" :style="{ top: `${readoutTop}px` }">
      <div class="world-place">
        <svg class="pin" viewBox="0 0 24 24" aria-hidden="true">
          <path
            fill="currentColor"
            fill-rule="evenodd"
            d="M12 2.2c-3.7 0-6.7 2.9-6.7 6.6 0 4.9 6.7 12.6 6.7 12.6s6.7-7.7 6.7-12.6c0-3.7-3-6.6-6.7-6.6zm0 9a2.4 2.4 0 1 1 0-4.8 2.4 2.4 0 0 1 0 4.8z"
          />
        </svg>
        <span>深圳市</span>
      </div>
      <time class="world-time" :class="{ night: isNightMode }">{{ timeText }}</time>
    </div>
  </div>
</template>

<script lang="ts">
// [0][1] 色点两半；[2] 城市标记、地名、定位针。时间白天为白色。
export const worldColors = [
  ['#ffd18a', '#f29e19', '#f29e19', '#ffffff', '#f29e19', '#ffffff'],
  ['#9fd4f8', '#3b82c4', '#3b82c4', '#ffffff', '#3b82c4', '#ffffff'],
  ['#d6bcfa', '#8b5cf6', '#8b5cf6', '#ffffff', '#8b5cf6', '#ffffff'],
  ['#9ae6b4', '#1f9d62', '#1f9d62', '#ffffff', '#1f9d62', '#ffffff'],
  ['#ffc1d1', '#e2557a', '#e2557a', '#ffffff', '#e2557a', '#ffffff'],
  ['#b6f3f5', '#14919b', '#14919b', '#ffffff', '#14919b', '#ffffff'],
]

export const worldNightModeColors = ['rgb(148,4,7)', 'rgb(111,26,23)']
</script>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

// 从上到下是北到南，从左到右是西到东。0 没点，1 陆地点，2 橙色城市点。
// 行数和每行字符数就是网格。点在地图宽高里均匀排开，每行字符数要一样。
const WORLD_DOTS: readonly string[] = [
  '0000000000000000000001111111111111110001100000000000000000000000000',
  '0000000000000000111111110111111111000000000000111000111110001110000',
  '0000000001111111111111111100111111000001111001101111111111111111111',
  '0000000111111111111111100110111101100001111111111111111111111111111',
  '0000000001111111111111001110011000000011111111111111111111111111111',
  '0000000011000011111111101111000000011011111111111111111111111001100',
  '0000000000000001111111111111100000001111111111111111111111111101000',
  '0000000000000000111111111111000000011211110011111111111111110000000',
  '0000000000000000111111112000000000011001111111111111111121101000000',
  '0000000000000000011111110000000000001110020111111111111111012000000',
  '0000000000000000001110010000000000011111111111111111111111000000000',
  '0000000000000000000110111100000000111111111111201111111110000000000',
  '0000000000000000000011100000000000111111111111100111011100000000000',
  '0000000000000000000000111110000000111111111111000110011101100000000',
  '0000000000000000000000001111100000011111111111000010011011100000000',
  '0000000000000000000000011111110000000011111110000000011111001000000',
  '0000000000000000000000011111111100000011111100000000000111001110000',
  '0000000000000000000000011111111100000001111101000000000000010110000',
  '0000000000000000000000001111111000000001111111000000000000111100000',
  '0000000000000000000000000111111000000001111010000000000011111110000',
  '0000000000000000000000000111110000000001111000000000000011111110000',
  '0000000000000000000000000111100000000000110000000000000011111110001',
  '0000000000000000000000000111000000000000000000000000000000001100001',
  '0000000000000000000000001110000000000000000000000000000000000000011',
  '0000000000000000000000001100000000000000000000000000000000000000010',
  '0000000000000000000000000100000000000000000000000000000000000000000',
]

// 点阵编辑。调好后把控制台打出的矩阵贴回这里，再把这里改成 false。
const WORLD_DOT_EDIT = false

const props = defineProps<{
  digits: string[]
  displayOpacity: (opacity: number) => number
  colonOpacity: number
  brightness: number
  isNightMode: boolean
}>()

const LON_WEST = -180
const LON_SPAN = 360
const LAT_NORTH = 78
const LAT_SOUTH = -56
const LAT_SPAN = LAT_NORTH - LAT_SOUTH

const root = ref<HTMLElement | null>(null)
const dotGrid = ref<string[]>(WORLD_DOTS.map((row) => row))
const activeRows = computed(() => (WORLD_DOT_EDIT ? dotGrid.value : WORLD_DOTS))
const width = ref(0)
const height = ref(0)
const now = ref(new Date())

let resizeObserver: ResizeObserver | undefined

const timeText = computed(() => {
  if (props.digits.length < 4) return ''
  return `${props.digits[0]}${props.digits[1]}:${props.digits[2]}${props.digits[3]}`
})

const frame = computed(() => {
  const w = width.value
  const h = height.value
  const aspect = LON_SPAN / LAT_SPAN
  const maxW = w * 0.94
  const maxH = h * 0.86
  let mapW = maxW
  let mapH = mapW / aspect
  if (mapH > maxH) {
    mapH = maxH
    mapW = mapH * aspect
  }
  const left = (w - mapW) / 2
  // 时间挂在地图下沿，整组在屏幕里居中，避免地图贴底留下大片空白。
  const hang = h * 0.1
  const top = Math.max(h * 0.02, (h - mapH - hang) / 2)
  return { left, top, width: mapW, height: mapH }
})

const readoutTop = computed(() => frame.value.top + frame.value.height * 0.7)

const gridSize = computed(() => ({
  cols: Math.max(activeRows.value[0]?.length ?? 1, 1),
  rows: Math.max(activeRows.value.length, 1),
}))

const cell = computed(() => Math.min(
  frame.value.width / gridSize.value.cols,
  frame.value.height / gridSize.value.rows,
))
const dotRadius = computed(() => cell.value * 0.3)
const cityRadius = computed(() => Math.max(dotRadius.value * 2.15, frame.value.width * 0.0085))
const terminatorWidth = computed(() => Math.max(1.25, cell.value * 0.16))

const project = (lon: number, lat: number) => {
  const box = frame.value
  return {
    x: box.left + ((lon - LON_WEST) / LON_SPAN) * box.width,
    y: box.top + ((LAT_NORTH - lat) / LAT_SPAN) * box.height,
  }
}

const solar = computed(() => {
  const date = now.value
  const start = Date.UTC(date.getUTCFullYear(), 0, 0)
  const day = Math.floor((date.getTime() - start) / 86_400_000)
  const decl = -23.44 * Math.cos((2 * Math.PI * (day + 10)) / 365)
  const utc = date.getUTCHours() + date.getUTCMinutes() / 60 + date.getUTCSeconds() / 3600
  const subLon = (12 - utc) * 15
  return { decl, subLon }
})

const wrapLon = (lon: number) => {
  let value = ((lon + 180) % 360 + 360) % 360 - 180
  if (value === -180) value = 180
  return value
}

const hourLimit = (lat: number, decl: number) => {
  const cosH = -Math.tan((lat * Math.PI) / 180) * Math.tan((decl * Math.PI) / 180)
  if (cosH < -1 || cosH > 1) return null
  return (Math.acos(cosH) * 180) / Math.PI
}

const isDay = (lat: number, lon: number) => {
  const { decl, subLon } = solar.value
  const limit = hourLimit(lat, decl)
  if (limit == null) {
    return -Math.tan((lat * Math.PI) / 180) * Math.tan((decl * Math.PI) / 180) < -1
  }
  let delta = lon - subLon
  delta = wrapLon(delta)
  return Math.abs(delta) <= limit
}

const cellPoint = (rowIndex: number, col: number) => {
  const { cols, rows } = gridSize.value
  const lon = LON_WEST + ((col + 0.5) / cols) * LON_SPAN
  const lat = LAT_NORTH - ((rowIndex + 0.5) / rows) * LAT_SPAN
  return { lon, lat, ...project(lon, lat) }
}

const dots = computed(() => {
  const placed: { x: number; y: number; day: boolean; off: boolean; row: number; col: number }[] = []
  activeRows.value.forEach((row, rowIndex) => {
    for (let col = 0; col < row.length; col += 1) {
      const mark = row[col]
      if (mark === '2') continue
      if (!WORLD_DOT_EDIT && mark !== '1') continue
      const point = cellPoint(rowIndex, col)
      placed.push({
        x: point.x,
        y: point.y,
        day: mark === '1' && isDay(point.lat, point.lon),
        off: mark !== '1',
        row: rowIndex,
        col,
      })
    }
  })
  return placed
})

const cities = computed(() => {
  const placed: { x: number; y: number }[] = []
  activeRows.value.forEach((row, rowIndex) => {
    for (let col = 0; col < row.length; col += 1) {
      if (row[col] !== '2') continue
      const point = cellPoint(rowIndex, col)
      placed.push({ x: point.x, y: point.y })
    }
  })
  return placed
})

const logWorldDots = () => {
  const lines = dotGrid.value.map((row) => `  '${row}',`).join('\n')
  console.log(`const WORLD_DOTS: readonly string[] = [\n${lines}\n]`)
}

const onDotClick = (row: number, col: number) => {
  if (!WORLD_DOT_EDIT) return
  const chars = dotGrid.value[row].split('')
  const mark = chars[col]
  if (mark === '2' || mark === undefined) return
  chars[col] = mark === '1' ? '0' : '1'
  dotGrid.value[row] = chars.join('')
  dotGrid.value = [...dotGrid.value]
  logWorldDots()
}

const onWorldPointerDown = (event: PointerEvent) => {
  if (!WORLD_DOT_EDIT) return
  event.stopPropagation()
}

const terminatorPath = computed(() => {
  const { decl, subLon } = solar.value
  let path = ''
  for (const sign of [1, -1]) {
    let drawing = false
    let prevLon: number | null = null
    for (let lat = LAT_SOUTH - 16; lat <= LAT_NORTH + 16; lat += 2) {
      const limit = hourLimit(lat, decl)
      if (limit == null) {
        drawing = false
        prevLon = null
        continue
      }
      const lon = wrapLon(subLon + sign * limit)
      if (prevLon != null && Math.abs(lon - prevLon) > 120) drawing = false
      const point = project(lon, lat)
      path += `${drawing ? 'L' : 'M'}${point.x.toFixed(1)} ${point.y.toFixed(1)}`
      drawing = true
      prevLon = lon
    }
  }
  return path
})

const syncSize = () => {
  const el = root.value
  if (!el) return
  width.value = el.clientWidth
  height.value = el.clientHeight
}

watch(() => props.digits, () => {
  now.value = new Date()
})

onMounted(() => {
  now.value = new Date()
  syncSize()
  if (root.value) {
    resizeObserver = new ResizeObserver(syncSize)
    resizeObserver.observe(root.value)
  }
})

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
})
</script>

<style scoped>
.world {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: #000;
}

.world-svg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  display: block;
}

.dot-day {
  fill: #d0d0d0;
}

.dot-night {
  fill: #4c4c4c;
}

.dot-off {
  fill: #d0d0d0;
  opacity: 0.16;
}

.world.editing .world-svg {
  cursor: pointer;
}

.dot-edit-log {
  position: absolute;
  top: 4vh;
  right: 4vw;
  z-index: 2;
  padding: 0.6em 0.9em;
  border: 0;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.16);
  color: #fff;
  font: inherit;
  font-size: 14px;
  cursor: pointer;
}

.terminator {
  fill: none;
  stroke: #c6c6c6;
  stroke-linecap: round;
  stroke-linejoin: round;
  pointer-events: none;
}

.city {
  fill: var(--color-0);
  pointer-events: none;
}

.world-readout {
  position: absolute;
  left: 4.6vw;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.world-place {
  display: flex;
  align-items: center;
  gap: 0.45em;
  margin-bottom: 0.28em;
  color: var(--color-0);
  font-size: 2.7vw;
  font-weight: 700;
  line-height: 1;
  letter-spacing: 0.04em;
}

.pin {
  width: 1.05em;
  height: 1.05em;
  display: block;
}

.world-time {
  font-family: 'SF-Pro-Rounded-Black', sans-serif;
  font-size: 11vw;
  line-height: 0.85;
  letter-spacing: -0.035em;
  color: #fff;
  font-variant-numeric: tabular-nums;
}

.world-time.night {
  color: var(--color-0);
}
</style>
