<template>
  <div v-if="digits.length > 0" class="numerical-wrapper" :style="{ opacity: displayOpacity(brightness) }">
    <svg class="numerical-filter" aria-hidden="true" focusable="false">
      <filter
        id="numerical-round"
        x="-50%"
        y="-50%"
        width="200%"
        height="200%"
        color-interpolation-filters="sRGB"
      >
        <feGaussianBlur in="SourceAlpha" :stdDeviation="cornerRadius" result="blur" />
        <feComponentTransfer in="blur" result="mask">
          <feFuncA type="linear" slope="28" intercept="-14" />
        </feComponentTransfer>
        <feComposite in="SourceGraphic" in2="mask" operator="in" />
      </filter>
    </svg>
    <div class="numerical-inner">
      <time class="numerical-time">
        <template v-for="(group, groupIndex) in digitGroups" :key="groupIndex">
          <span v-if="groupIndex === 1" class="numerical-colon" aria-hidden="true"></span>
          <span class="numerical-digits">
            <span v-for="slot in group" :key="slot.index" class="numerical-digit-slot">
              <span
                v-for="layer in slot.layers"
                :key="layer.id"
                class="numerical-digit-frame"
                :class="{
                  'numerical-digit-enter': layer.phase === 'enter',
                  'numerical-digit-leave': layer.phase === 'leave',
                }"
                :aria-hidden="layer.phase === 'leave'"
                @animationend="onDigitAnimationEnd(slot.index, layer.id, $event)"
              ><span class="numerical-digit">{{ layer.char }}</span></span>
            </span>
          </span>
        </template>
      </time>
      <div class="numerical-side">
        <div class="numerical-date-row">
          <span class="numerical-day" :style="{ color: isNightMode ? 'var(--color-0)' : '#fff' }">{{ dayOfMonth }} </span><span class="numerical-weekday">周{{ weekday }}</span>
        </div>
        <div v-if="displayedStocks.length > 0" class="numerical-stocks">
          <div
            v-for="item in displayedStocks"
            :key="item.value"
            class="stock-row"
            :class="{ 'night-mode': isNightMode }"
          >
            <span class="stock-name">{{ item.name }}</span>
            <span class="stock-price">{{ item.price }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export const numericalColors = [
  ['rgb(232,118,102)', 'rgb(242,182,125)', 'rgb(232,118,102)', 'rgb(242,182,125)', 'rgb(232,118,102)', 'rgb(242,182,125)'],
  ['rgb(65,135,225)', 'rgb(135,205,248)', 'rgb(65,135,225)', 'rgb(135,205,248)', 'rgb(65,135,225)', 'rgb(135,205,248)'],
  ['rgb(148,88,228)', 'rgb(198,158,242)', 'rgb(148,88,228)', 'rgb(198,158,242)', 'rgb(148,88,228)', 'rgb(198,158,242)'],
  ['rgb(52,178,102)', 'rgb(138,225,168)', 'rgb(52,178,102)', 'rgb(138,225,168)', 'rgb(52,178,102)', 'rgb(138,225,168)'],
  ['rgb(218,178,48)', 'rgb(245,218,108)', 'rgb(218,178,48)', 'rgb(245,218,108)', 'rgb(218,178,48)', 'rgb(245,218,108)'],
  ['rgb(228,72,132)', 'rgb(242,158,182)', 'rgb(228,72,132)', 'rgb(242,158,182)', 'rgb(228,72,132)', 'rgb(242,158,182)'],
  ['rgb(42,168,188)', 'rgb(102,215,225)', 'rgb(42,168,188)', 'rgb(102,215,225)', 'rgb(42,168,188)', 'rgb(102,215,225)'],
  ['rgb(222,72,62)', 'rgb(242,158,98)', 'rgb(222,72,62)', 'rgb(242,158,98)', 'rgb(222,72,62)', 'rgb(242,158,98)'],
]

export const numericalNightModeColors = ['rgb(148,4,7)', 'rgb(111,26,23)']
</script>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { clockNow } from '@/composables/useTime'
import { useFloatConfig } from '@/composables/useFloatConfig'
import { extractStockItems, fetchGlanceMenu, type StockItem } from '@/composables/useGlanceStocks'

const props = defineProps<{
  digits: string[]
  displayOpacity: (opacity: number) => number
  colonOpacity: number
  brightness: number
  isNightMode: boolean
}>()

type DigitPhase = 'settled' | 'enter' | 'leave'

type DigitLayer = {
  id: number
  char: string
  phase: DigitPhase
}

type DigitSlot = {
  index: number
  layers: DigitLayer[]
}

const digitSlots = ref<DigitSlot[]>([0, 1, 2, 3].map((index) => ({ index, layers: [] })))
let digitLayerId = 0

const prefersReducedMotion = () => window.matchMedia('(prefers-reduced-motion: reduce)').matches

// 变化的那一位单独过渡：旧数字加速上移变浅，新数字同时从偏低处出现并减速升到原位。
const syncDigitSlots = (digits: string[]) => {
  for (let index = 0; index < 4; index += 1) {
    const char = digits[index]
    if (!char) continue
    const slot = digitSlots.value[index]
    const current = slot.layers.find((layer) => layer.phase !== 'leave')
    if (!current) {
      slot.layers.push({ id: digitLayerId++, char, phase: 'settled' })
      continue
    }
    if (current.char === char) continue
    if (prefersReducedMotion()) {
      slot.layers = [{ id: digitLayerId++, char, phase: 'settled' }]
      continue
    }
    current.phase = 'leave'
    slot.layers.push({ id: digitLayerId++, char, phase: 'enter' })
  }
}

watch(
  () => props.digits,
  (digits) => syncDigitSlots(digits),
  { immediate: true },
)

const digitGroups = computed(() => [digitSlots.value.slice(0, 2), digitSlots.value.slice(2)])

const onDigitAnimationEnd = (index: number, id: number, event: AnimationEvent) => {
  if (event.target !== event.currentTarget) return
  const slot = digitSlots.value[index]
  const layer = slot?.layers.find((item) => item.id === id)
  if (!layer) return
  if (event.animationName.includes('leave') && layer.phase === 'leave') {
    slot.layers = slot.layers.filter((item) => item.id !== id)
    return
  }
  if (event.animationName.includes('enter') && layer.phase === 'enter') {
    layer.phase = 'settled'
  }
}

const { numericalStocks } = useFloatConfig()

const now = ref(clockNow())
let dateTimer: number | undefined

const updateDate = () => {
  now.value = clockNow()
}

// 股票数据管理
const stockList = ref<StockItem[]>([])
let stockTimer: number | undefined
let isFetching = false

const displayedStocks = computed(() => {
  if (!numericalStocks.value || numericalStocks.value.length === 0) return []
  return numericalStocks.value.map((key) => {
    const found = stockList.value.find((item) => item.value === key)
    if (found) return found
    return {
      value: key,
      name: key.replace(/^stocks:|^fx:/, ''),
      price: '--',
      fullTitle: key,
    }
  })
})

const stopStockTimer = () => {
  if (stockTimer !== undefined) {
    clearTimeout(stockTimer)
    stockTimer = undefined
  }
}

const scheduleNextFetch = (delayMs: number) => {
  stopStockTimer()
  if (!numericalStocks.value || numericalStocks.value.length === 0 || document.hidden) return
  stockTimer = window.setTimeout(() => {
    void fetchStocks()
  }, delayMs)
}

const fetchStocks = async () => {
  if (!numericalStocks.value || numericalStocks.value.length === 0 || isFetching) return
  isFetching = true
  try {
    const data = await fetchGlanceMenu()
    stockList.value = extractStockItems(data.menu)
    const interval = Math.max(data.refresh_after_seconds || 3, 3) * 1000
    scheduleNextFetch(interval)
  } catch {
    // 请求失败时不打断，5 秒后重试
    scheduleNextFetch(5000)
  } finally {
    isFetching = false
  }
}

const startStockFetch = () => {
  stopStockTimer()
  if (numericalStocks.value && numericalStocks.value.length > 0 && !document.hidden) {
    void fetchStocks()
  }
}

const onVisibilityChange = () => {
  if (document.hidden) {
    stopStockTimer()
  } else {
    startStockFetch()
  }
}

watch(
  () => numericalStocks.value,
  (newVal) => {
    if (!newVal || newVal.length === 0) {
      stopStockTimer()
      stockList.value = []
    } else {
      startStockFetch()
    }
  },
  { deep: true, immediate: true },
)

onMounted(() => {
  dateTimer = window.setInterval(updateDate, 60_000)
  document.addEventListener('visibilitychange', onVisibilityChange)
  window.addEventListener('resize', onResize)
})

onBeforeUnmount(() => {
  if (dateTimer !== undefined) clearInterval(dateTimer)
  document.removeEventListener('visibilitychange', onVisibilityChange)
  window.removeEventListener('resize', onResize)
  stopStockTimer()
})

const dayOfMonth = computed(() => now.value.getDate())
const weekday = computed(() => ['日', '一', '二', '三', '四', '五', '六'][now.value.getDay()])

const viewportWidth = ref(window.innerWidth)
// 字号是 50vw。拐角圆度约为字号的 1%。取整，避免 Safari 在小数模糊半径处把笔画切开。
const cornerRadius = computed(() => Math.max(1, Math.round(viewportWidth.value * 0.005)))

const onResize = () => {
  viewportWidth.value = window.innerWidth
}
</script>

<style scoped>
@property --digit-wash {
  syntax: '<number>';
  inherits: false;
  initial-value: 0;
}

.numerical-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
}

/* 内部保持时间 + 右侧信息排布 */
.numerical-inner {
  display: flex;
  align-items: flex-start;
  justify-content: center;
  max-width: 100vw;
  padding: 0 4vw;
  margin-bottom: 8vh;
}

.numerical-filter {
  position: absolute;
  width: 0;
  height: 0;
}

.numerical-time {
  display: flex;
  align-items: stretch;
  font-family: 'AFCamberwell-One-Regular', serif;
  font-size: 50vw;
  line-height: 1;
  flex-shrink: 0;
}

.numerical-digits {
  display: flex;
  align-items: flex-start;
}

.numerical-digit-slot {
  position: relative;
}

.numerical-digit-slot + .numerical-digit-slot {
  margin-left: -0.8vw;
}

.numerical-digit-frame {
  display: block;
}

.numerical-digit {
  --digit-wash: 0;
  /* iPhone 的 WebKit 不绘制 display:block 的 background-clip:text，字会整段消失，冒号不受影响。 */
  display: inline-block;
  /* iOS 13 没有 color-mix，祖先还有 transform。透明填充配背景裁切和 SVG 滤镜时，四个字会整段消失。 */
  color: var(--color-0);
}

@supports (aspect-ratio: 1 / 1) {
  .numerical-digit-frame {
    filter: url(#numerical-round);
    /* 样式轨道的 transform 会让 iPad 的 WebKit 把滤镜从字形中间切开。单独建一层合成层避开。 */
    transform: translate3d(0, 0, 0);
    -webkit-backface-visibility: hidden;
    backface-visibility: hidden;
  }

  .numerical-digit {
    -webkit-box-decoration-break: clone;
    box-decoration-break: clone;
    background: linear-gradient(to bottom, var(--color-0), var(--color-1));
    background: linear-gradient(
      to bottom,
      color-mix(in srgb, var(--color-0) calc((1 - var(--digit-wash)) * 100%), white),
      color-mix(in srgb, var(--color-1) calc((1 - var(--digit-wash)) * 100%), white)
    );
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    color: transparent;
  }
}

.numerical-digit-leave {
  position: absolute;
  left: 0;
  top: 0;
  z-index: 0;
  animation: numerical-digit-leave 0.2s ease-in forwards;
}

.numerical-digit-enter {
  position: relative;
  z-index: 1;
  animation: numerical-digit-enter 0.7s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes numerical-digit-leave {
  from {
    transform: translate3d(0, 0, 0);
    opacity: 1;
    --digit-wash: 0;
  }
  to {
    transform: translate3d(0, -1.6vw, 0);
    opacity: 0;
    --digit-wash: 1;
  }
}

@keyframes numerical-digit-enter {
  from {
    transform: translate3d(0, 1vw, 0);
    opacity: 0.2;
  }
  to {
    transform: translate3d(0, 0, 0);
    opacity: 1;
  }
}

@media (prefers-reduced-motion: reduce) {
  .numerical-digit-enter,
  .numerical-digit-leave {
    animation: none;
  }
}

.numerical-colon {
  width: 7.5vw;
  margin: 0 1.2vw;
  align-self: stretch;
  background: linear-gradient(to bottom, var(--color-0), var(--color-1));
  -webkit-mask-image:
    radial-gradient(circle closest-side, #000 99%, transparent 100%),
    radial-gradient(circle closest-side, #000 99%, transparent 100%);
  -webkit-mask-size: 7.5vw 7.5vw, 7.5vw 7.5vw;
  -webkit-mask-repeat: no-repeat, no-repeat;
  -webkit-mask-position: center 27%, center 73%;
  mask-image:
    radial-gradient(circle closest-side, #000 99%, transparent 100%),
    radial-gradient(circle closest-side, #000 99%, transparent 100%);
  mask-size: 7.5vw 7.5vw, 7.5vw 7.5vw;
  mask-repeat: no-repeat, no-repeat;
  mask-position: center 27%, center 73%;
}

.numerical-side {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  margin-left: 4vw;
  padding-top: 4.45vw;
  width: 26vw;
  max-width: 26vw;
  box-sizing: border-box;
  overflow: hidden;
  white-space: nowrap;
}

/* 横屏手机高度小于字号，overflow 会把数字裁掉。矮屏时把字号收到视口里。 */
@media (max-height: 520px) {
  .numerical-time {
    font-size: min(50vw, 78vh);
    font-size: min(50vw, 78dvh);
  }

  .numerical-side {
    padding-top: calc(min(50vw, 78vh) * 0.089);
    padding-top: calc(min(50vw, 78dvh) * 0.089);
  }
}

.numerical-date-row {
  display: flex;
  flex-direction: row;
  align-items: baseline;
  white-space: nowrap;
}

.numerical-day {
  font-size: 5vw;
  white-space: nowrap;
  padding-right: 1vw;
}

.numerical-weekday {
  font-size: 5vw;
  color: var(--color-0);
}

.numerical-stocks {
  display: flex;
  flex-direction: column;
  gap: 0.8vw;
  margin-top: 1.8vw;
  width: 100%;
  box-sizing: border-box;
}

.stock-row {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  gap: 1.2vw;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'SF Pro Text', 'Helvetica Neue', sans-serif;
  font-size: 2.6vw;
  line-height: 1.25;
}

.stock-name {
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.stock-price {
  color: #fff;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  text-align: right;
  white-space: nowrap;
  flex-shrink: 0;
}

.stock-row.night-mode .stock-name {
  color: var(--color-0);
  opacity: 0.8;
}

.stock-row.night-mode .stock-price {
  color: var(--color-0);
}
</style>
