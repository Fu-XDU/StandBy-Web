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
      <defs>
        <linearGradient
          id="terminator-shade"
          gradientUnits="userSpaceOnUse"
          x1="0"
          :y1="shadeY.north"
          x2="0"
          :y2="shadeY.south"
        >
          <stop offset="0" stop-color="#d2d2d2" />
          <stop offset="0.4" stop-color="#c4c4c4" />
          <stop offset="0.58" stop-color="#7a7a7a" />
          <stop offset="0.78" stop-color="#3a3a3a" />
          <stop offset="1" stop-color="#1c1c1c" />
        </linearGradient>
        <!-- iPad 的 WebKit 用渐变描边时，会在路径最上沿多画一条水平线。渐变改铺在矩形上，线条只用实色遮罩。 -->
        <mask id="terminator-mask" maskUnits="userSpaceOnUse" x="0" y="0" :width="width" :height="height">
          <path
            class="terminator"
            :d="terminatorPath"
            fill="none"
            stroke="#fff"
            :stroke-width="terminatorWidth"
          />
        </mask>
        <clipPath
          v-for="dot in cutDots"
          :id="`dot-cut-${dot.row}-${dot.col}`"
          :key="`clip-${dot.row}-${dot.col}`"
        >
          <polygon :points="dot.clip" />
        </clipPath>
      </defs>
      <circle
        v-for="dot in solidDots"
        :key="`dot-${dot.row}-${dot.col}`"
        :cx="dot.x"
        :cy="dot.y"
        :r="dotRadius"
        :class="dot.off ? 'dot-off' : dot.day ? 'dot-day' : 'dot-night'"
        @pointerdown.stop
        @click.stop="onMarkerClick(dot.row, dot.col)"
      />
      <circle
        v-for="dot in cutDots"
        :key="`cut-night-${dot.row}-${dot.col}`"
        class="dot-night"
        :cx="dot.x"
        :cy="dot.y"
        :r="dotRadius"
        @pointerdown.stop
        @click.stop="onMarkerClick(dot.row, dot.col)"
      />
      <circle
        v-for="dot in cutDots"
        :key="`cut-day-${dot.row}-${dot.col}`"
        class="dot-day"
        :cx="dot.x"
        :cy="dot.y"
        :r="dotRadius"
        :clip-path="`url(#dot-cut-${dot.row}-${dot.col})`"
        @pointerdown.stop
        @click.stop="onMarkerClick(dot.row, dot.col)"
      />
      <rect
        class="terminator-shade"
        x="0"
        y="0"
        :width="width"
        :height="height"
        fill="url(#terminator-shade)"
        mask="url(#terminator-mask)"
      />
      <circle
        v-for="(city, index) in cities"
        :key="`city-${index}`"
        class="city"
        :cx="city.x"
        :cy="city.y"
        :r="cityRadius"
        @pointerdown.stop
        @click.stop="onMarkerClick(city.row, city.col)"
      />
      <Transition name="select-pop">
        <g
          v-if="selectedPoint"
          :key="`${selectedPoint.row}-${selectedPoint.col}`"
          class="selected-mark"
        >
          <circle
            class="ring"
            :cx="selectedPoint.x"
            :cy="selectedPoint.y"
            :r="selectedRingRadius"
            :stroke-width="ringWidth"
          />
          <circle
            class="city"
            :cx="selectedPoint.x"
            :cy="selectedPoint.y"
            :r="cityRadius"
            @pointerdown.stop
            @click.stop="onMarkerClick(selectedPoint.row, selectedPoint.col)"
          />
        </g>
      </Transition>
    </svg>

    <button v-if="WORLD_DOT_EDIT" type="button" class="dot-edit-log" @click.stop="logWorldDots">
      输出矩阵
    </button>

    <div
      v-if="cityEditorOpen"
      class="city-editor-backdrop"
      @pointerdown.stop
      @mousedown.stop
      @touchstart.stop
      @click="saveCity"
    />
    <div class="world-readout" :class="{ editing: cityEditorOpen }" :style="{ top: `${readoutTop}px` }">
      <form
        v-if="cityEditorOpen"
        class="city-editor"
        @submit.prevent="saveCity"
        @pointerdown.stop
        @mousedown.stop
        @touchstart.stop
        @click.stop
      >
        <input
          ref="cityInput"
          v-model="cityDraft"
          class="city-input"
          type="text"
          maxlength="20"
          enterkeyhint="done"
          autocomplete="off"
          spellcheck="false"
          placeholder="城市名称"
          aria-label="城市名称"
          @keydown.enter.prevent="saveCity"
          @keydown.esc.prevent="closeCityEditor"
        />
      </form>
      <button
        type="button"
        class="world-place"
        aria-label="修改城市名称"
        @pointerdown.stop
        @mousedown.stop
        @touchstart.stop
        @click.stop="openCityEditor"
      >
        <span class="pin" aria-hidden="true" />
        <span v-if="cityName" class="world-city">{{ cityName }}</span>
        <span v-else class="world-city-blank" aria-hidden="true" />
      </button>
      <time class="world-time" :class="{ night: isNightMode }">
        <span
          v-for="(part, index) in timeParts"
          :key="index"
          :class="part.colon ? 'world-colon' : 'world-digit'"
          :style="part.gap ? { marginRight: part.gap } : undefined"
        >{{ part.char }}</span>
      </time>
    </div>
  </div>
</template>

<script lang="ts">
// [0][1] 色点两半；[2] 城市标记、地名、定位针。时间白天为白色。
export const worldColors = [
  ['#ffd18a', '#ff9f0c', '#ff9f0c', '#ffffff', '#ff9f0c', '#ffffff'],
  ['#9fd4f8', '#3b82c4', '#3b82c4', '#ffffff', '#3b82c4', '#ffffff'],
  ['#d6bcfa', '#8b5cf6', '#8b5cf6', '#ffffff', '#8b5cf6', '#ffffff'],
  ['#9ae6b4', '#1f9d62', '#1f9d62', '#ffffff', '#1f9d62', '#ffffff'],
  ['#ffc1d1', '#e2557a', '#e2557a', '#ffffff', '#e2557a', '#ffffff'],
  ['#b6f3f5', '#14919b', '#14919b', '#ffffff', '#14919b', '#ffffff'],
]

export const worldNightModeColors = ['rgb(148,4,7)', 'rgb(111,26,23)']
</script>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useLocalStorageSync } from '@/composables/useLocalStorageSync'
import { useFloatConfig } from '@/composables/useFloatConfig'

// 从上到下是北到南，从左到右是西到东。0 没点，1 陆地点，2 橙色城市点。
// 行数和每行字符数就是网格。点在地图宽高里均匀排开，每行字符数要一样。
const WORLD_DOTS: readonly string[] = [
  '000000000000000001111111111111110001100000000000000000000000000',
  '000000000000111111110111111111000000000000111000111110001110000',
  '000001111111111111111100111111000001111001101111111111111111111',
  '000111111111111111100110111101100001111111111111111111111111111',
  '000001111111111111001110011000000011111111111111111111111111111',
  '000011000011111111101111000000011011111111111111111111111001100',
  '000000000001111111111111100000001111111111111111111111111101000',
  '000000000000111111111111000000011211110011111111111111110000000',
  '000000000000111111112000000000011001111111111111111121101000000',
  '000000000000011111110000000000001110020111111111111111012000000',
  '000000000000001110010000000000011111111111111111111111000000000',
  '000000000000000110111100000000111111111111201111111110000000000',
  '000000000000000011100000000000111111111111100111011100000000000',
  '000000000000000000111110000000111111111111000110011101100000000',
  '000000000000000000001111100000011111111111000010011011100000000',
  '000000000000000000011111110000000011111110000000011111001000000',
  '000000000000000000011111111100000011111100000000000111001110000',
  '000000000000000000011111111100000001111101000000000000010110000',
  '000000000000000000001111111000000001111111000000000000111100000',
  '000000000000000000000111111000000001111010000000000011111110000',
  '000000000000000000000111110000000001111000000000000011111110000',
  '000000000000000000000111100000000000110000000000000011111110001',
  '000000000000000000000111000000000000000000000000000000001100001',
  '000000000000000000001110000000000000000000000000000000000000011',
  '000000000000000000001100000000000000000000000000000000000000010',
  '000000000000000000000100000000000000000000000000000000000000000',
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
// 晨昏线整体南移的纬度。加大后北极更早进入极夜，圆拱往下。
const TERMINATOR_SOUTH_SHIFT = 11
// 南端圆拱额外北抬的纬度。只抬赤道以南，北极极夜的位置不变。
const TERMINATOR_SOUTH_LIFT = 20
// 晨昏线东西偏移的经度。正数往东（地图右侧），负数往西。
const TERMINATOR_LON_SHIFT = 14.5
// 白天东西宽度的比例。1 是当前宽度，小于 1 变窄，大于 1 变宽。
const TERMINATOR_DAY_SCALE = 0.94

const dayHalfWidth = (hour: number) =>
  hour - (1 - TERMINATOR_DAY_SCALE) * 90 * Math.sin((hour * Math.PI) / 180)

const placeLat = (raw: number) => {
  const shifted = raw - TERMINATOR_SOUTH_SHIFT
  if (shifted >= 0) return shifted
  const t = Math.min(1, -shifted / 90)
  return shifted + TERMINATOR_SOUTH_LIFT * t
}

// 和 placeLat 互逆，让点的昼夜判断落在抬高后的线上。
const solarLat = (lat: number) => {
  if (lat >= 0) return lat + TERMINATOR_SOUTH_SHIFT
  const scale = 1 - TERMINATOR_SOUTH_LIFT / 90
  return lat / scale + TERMINATOR_SOUTH_SHIFT
}

const root = ref<HTMLElement | null>(null)
const cityName = ref('深圳市')
const cityEditorOpen = ref(false)
const cityDraft = ref('')
const cityInput = ref<HTMLInputElement | null>(null)
useLocalStorageSync('worldCityName', cityName)

const openCityEditor = () => {
  cityDraft.value = cityName.value
  cityEditorOpen.value = true
  nextTick(() => {
    cityInput.value?.focus()
    cityInput.value?.select()
  })
}

const closeCityEditor = () => {
  cityEditorOpen.value = false
}

const saveCity = () => {
  cityName.value = cityDraft.value.trim().slice(0, 20)
  cityEditorOpen.value = false
}
const { worldSelected } = useFloatConfig()
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

// SF Pro Text 各字形左右留白占字号的比例。1 的右侧很宽，4 几乎贴边，统一字距会把 4 和 8 挤在一起。
const DIGIT_SIDE: Record<string, [number, number]> = {
  '0': [0.035, 0.04],
  '1': [0.035, 0.09],
  '2': [0.045, 0.05],
  '3': [0.035, 0.04],
  '4': [0.015, 0.02],
  '5': [0.035, 0.04],
  '6': [0.035, 0.04],
  '7': [0.025, 0.045],
  '8': [0.03, 0.035],
  '9': [0.035, 0.04],
  ':': [0.055, 0.06],
}
const TIME_GAP = 0.06

const timeParts = computed(() => {
  const chars = timeText.value.split('')
  return chars.map((char, index) => {
    const next = chars[index + 1]
    const right = DIGIT_SIDE[char]?.[1] ?? 0
    const left = next ? (DIGIT_SIDE[next]?.[0] ?? 0) : 0
    const margin = next ? TIME_GAP - right - left : 0
    return {
      char,
      colon: char === ':',
      gap: next ? `${margin.toFixed(3)}em` : '',
    }
  })
})

const frame = computed(() => {
  const w = width.value
  const h = height.value
  const aspect = LON_SPAN / LAT_SPAN
  // 宽度跟屏幕对齐。高度超出时才按高度把地图缩小。
  const maxW = w
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

const readoutTop = computed(() => frame.value.top + frame.value.height * 0.69)

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
const ringWidth = computed(() => Math.max(1.5, dotRadius.value * 1.3))
const selectedRingRadius = computed(() => cityRadius.value + ringWidth.value / 2 - Math.max(1.5, ringWidth.value * 0.2))
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
  const subLon = (12 - utc) * 15 + TERMINATOR_LON_SHIFT
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
  const shifted = solarLat(lat)
  const limit = hourLimit(shifted, decl)
  if (limit == null) {
    return -Math.tan((shifted * Math.PI) / 180) * Math.tan((decl * Math.PI) / 180) < -1
  }
  let delta = lon - subLon
  delta = wrapLon(delta)
  return Math.abs(delta) <= dayHalfWidth(limit)
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
  const placed: { x: number; y: number; row: number; col: number }[] = []
  activeRows.value.forEach((row, rowIndex) => {
    for (let col = 0; col < row.length; col += 1) {
      if (row[col] !== '2') continue
      const point = cellPoint(rowIndex, col)
      placed.push({ x: point.x, y: point.y, row: rowIndex, col })
    }
  })
  return placed
})

const selectedPoint = computed(() => {
  const current = worldSelected.value
  if (!current) return null
  const row = activeRows.value[current.row]
  const mark = row?.[current.col]
  if (mark !== '1' && mark !== '2') return null
  const point = cellPoint(current.row, current.col)
  return { x: point.x, y: point.y, row: current.row, col: current.col }
})

const onMarkerClick = (row: number, col: number) => {
  if (WORLD_DOT_EDIT) {
    onDotClick(row, col)
    return
  }
  if (worldSelected.value?.row === row && worldSelected.value?.col === col) {
    worldSelected.value = null
    return
  }
  worldSelected.value = { row, col }
}

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

const terminatorGeometry = computed(() => {
  const { decl, subLon } = solar.value
  const tanDecl = Math.tan((decl * Math.PI) / 180)
  const subpaths: { x: number; y: number }[][] = []
  let current: { x: number; y: number }[] = []
  let path = ''
  let drawing = false
  let prevLon: number | null = null
  let prevLat = 0

  const line = (lon: number, lat: number) => {
    const point = project(lon, lat)
    if (!drawing) {
      if (current.length) subpaths.push(current)
      current = [point]
      path += `M${point.x.toFixed(1)} ${point.y.toFixed(1)}`
    } else {
      current.push(point)
      path += `L${point.x.toFixed(1)} ${point.y.toFixed(1)}`
    }
    drawing = true
    prevLon = lon
    prevLat = lat
  }
  const move = (lon: number, lat: number) => {
    drawing = false
    line(lon, lat)
  }
  // 越过日期变更线时先收到地图边缘，再从另一侧接着画，避免横穿整张地图。
  const segment = (lon: number, lat: number) => {
    const wrapped = wrapLon(lon)
    if (prevLon != null && Math.abs(wrapped - prevLon) > 90) {
      const leavingEast = prevLon > 0
      const midLat = (prevLat + lat) / 2
      line(leavingEast ? 180 : -180, midLat)
      move(leavingEast ? -180 : 180, midLat)
    }
    line(wrapped, lat)
  }

  // 分点时刻两条分界在南北各收成一点。用时角均匀取样，圆拱和底部那一扫才是顺的。
  if (Math.abs(tanDecl) < 1e-4) {
    for (const sign of [1, -1]) {
      drawing = false
      prevLon = null
      for (let lat = -89; lat <= 89; lat += 2) segment(subLon + sign * dayHalfWidth(90), placeLat(lat))
    }
  } else {
    const latAt = (hour: number) =>
      placeLat((Math.atan(-Math.cos((hour * Math.PI) / 180) / tanDecl) * 180) / Math.PI)
    for (let hour = 0; hour <= 180; hour += 1) segment(subLon + dayHalfWidth(hour), latAt(hour))
    for (let hour = 180; hour >= 0; hour -= 1) segment(subLon - dayHalfWidth(hour), latAt(hour))
  }
  if (current.length) subpaths.push(current)
  return { d: path, subpaths }
})

const terminatorPath = computed(() => terminatorGeometry.value.d)

const shadeY = computed(() => ({
  north: project(0, LAT_NORTH).y,
  south: project(0, LAT_SOUTH).y,
}))

const unproject = (x: number, y: number) => {
  const box = frame.value
  return {
    lon: LON_WEST + ((x - box.left) / box.width) * LON_SPAN,
    lat: LAT_NORTH - ((y - box.top) / box.height) * LAT_SPAN,
  }
}

const placedDots = computed(() => {
  const solid: typeof dots.value = []
  const cuts: (typeof dots.value[number] & { clip: string })[] = []
  const radius = dotRadius.value
  const subpaths = terminatorGeometry.value.subpaths
  for (const dot of dots.value) {
    if (dot.off || radius <= 0) {
      solid.push(dot)
      continue
    }
    let best: { dist: number; px: number; py: number; tx: number; ty: number } | null = null
    for (const stroke of subpaths) {
      for (let index = 1; index < stroke.length; index += 1) {
        const start = stroke[index - 1]
        const end = stroke[index]
        const dx = end.x - start.x
        const dy = end.y - start.y
        const len2 = dx * dx + dy * dy
        if (len2 < 0.01) continue
        const t = Math.min(1, Math.max(0, ((dot.x - start.x) * dx + (dot.y - start.y) * dy) / len2))
        const px = start.x + dx * t
        const py = start.y + dy * t
        const dist = Math.hypot(dot.x - px, dot.y - py)
        if (dist < radius && (!best || dist < best.dist)) best = { dist, px, py, tx: dx, ty: dy }
      }
    }
    if (!best) {
      solid.push(dot)
      continue
    }
    const tangent = Math.hypot(best.tx, best.ty) || 1
    const ux = best.tx / tangent
    const uy = best.ty / tangent
    let nx = dot.x - best.px
    let ny = dot.y - best.py
    const normal = Math.hypot(nx, ny)
    if (normal < 0.4) {
      nx = -uy
      ny = ux
      const probe = unproject(best.px + nx * 6, best.py + ny * 6)
      if (!isDay(probe.lat, probe.lon)) {
        nx = -nx
        ny = -ny
      }
    } else {
      nx /= normal
      ny /= normal
      if (!dot.day) {
        nx = -nx
        ny = -ny
      }
    }
    const along = radius * 4
    const across = along * 2
    const corners = [
      [best.px + ux * along, best.py + uy * along],
      [best.px - ux * along, best.py - uy * along],
      [best.px - ux * along + nx * across, best.py - uy * along + ny * across],
      [best.px + ux * along + nx * across, best.py + uy * along + ny * across],
    ]
    cuts.push({
      ...dot,
      clip: corners.map(([x, y]) => `${x.toFixed(1)},${y.toFixed(1)}`).join(' '),
    })
  }
  return { solid, cuts }
})

const solidDots = computed(() => placedDots.value.solid)
const cutDots = computed(() => placedDots.value.cuts)

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
  fill: #cccccc;
}

.dot-night {
  fill: #515151;
}

.dot-off {
  fill: #cccccc;
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
  stroke-linecap: round;
  stroke-linejoin: round;
  pointer-events: none;
}

.terminator-shade {
  pointer-events: none;
}

.city {
  fill: var(--color-0);
  cursor: pointer;
}

.ring {
  fill: none;
  stroke: #fff;
  pointer-events: none;
}

.selected-mark {
  transform-box: fill-box;
  transform-origin: center;
}

.select-pop-enter-active {
  transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.22s ease-out;
}

.select-pop-leave-active {
  transition: transform 0.2s ease-in, opacity 0.2s ease-in;
}

.select-pop-enter-from,
.select-pop-leave-to {
  opacity: 0;
  transform: scale(0.4);
}

@media (prefers-reduced-motion: reduce) {
  .select-pop-enter-active,
  .select-pop-leave-active {
    transition: none;
  }
}

.dot-day,
.dot-night {
  cursor: pointer;
}

.city-editor-backdrop {
  position: absolute;
  inset: 0;
  z-index: 2;
}

.world-readout {
  position: absolute;
  left: 4.6vw;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.world-readout.editing {
  z-index: 3;
}

.world-place {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 0.07em;
  margin: 0 0 0.08em;
  padding: 0;
  border: 0;
  background: none;
  color: var(--color-0);
  font-family: 'SF Pro Text', -apple-system, BlinkMacSystemFont, 'Helvetica Neue', sans-serif;
  font-size: 2.7vw;
  font-weight: 600;
  line-height: 1;
  letter-spacing: 0.04em;
  cursor: pointer;
}

.world-city {
  display: inline-block;
}

.world-city-blank {
  display: block;
  width: 3.2em;
  height: 1em;
}

.city-editor {
  position: absolute;
  left: 0;
  bottom: calc(100% + 0.35em);
}

.city-input {
  width: 8em;
  margin: 0;
  padding: 0.38em 0.55em;
  border: 1px solid rgba(255, 255, 255, 0.22);
  border-radius: 0.42em;
  background: rgba(28, 28, 30, 0.94);
  color: #fff;
  font-family: inherit;
  font-size: 0.92em;
  font-weight: 600;
  letter-spacing: 0.04em;
  line-height: 1.2;
  outline: none;
}

.city-input:focus {
  border-color: var(--color-0);
}

.pin {
  width: 1.05em;
  height: 1.05em;
  display: block;
  background-color: currentColor;
  -webkit-mask: url('@/assets/location-fill.svg') center / contain no-repeat;
  mask: url('@/assets/location-fill.svg') center / contain no-repeat;
}

.world-time {
  font-family: 'SF Pro Text', -apple-system, BlinkMacSystemFont, 'Helvetica Neue', sans-serif;
  font-weight: 600;
  font-size: 9vw;
  line-height: 0.85;
  letter-spacing: 0;
  color: #fff;
}

.world-digit,
.world-colon {
  display: inline-block;
}

.world-colon {
  transform: translateY(-0.1em);
}

.world-time.night {
  color: var(--color-0);
}
</style>
