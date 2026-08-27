<template>
  <div v-if="digits.length > 0" class="numerical-wrapper" :style="{ opacity: displayOpacity(brightness) }">
    <div class="numerical-inner">
      <time class="numerical-time">{{ digits[0] }}{{ digits[1] }}<span class="numerical-colon">:</span>{{ digits[2] }}{{ digits[3] }}</time>
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
import { useFloatConfig } from '@/composables/useFloatConfig'
import { extractStockItems, fetchGlanceMenu, type StockItem } from '@/composables/useGlanceStocks'

defineProps<{
  digits: string[]
  displayOpacity: (opacity: number) => number
  colonOpacity: number
  brightness: number
  isNightMode: boolean
}>()

const { numericalStocks } = useFloatConfig()

const now = ref(new Date())
let dateTimer: number | undefined

const updateDate = () => {
  now.value = new Date()
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
})

onBeforeUnmount(() => {
  if (dateTimer !== undefined) clearInterval(dateTimer)
  document.removeEventListener('visibilitychange', onVisibilityChange)
  stopStockTimer()
})

const dayOfMonth = computed(() => now.value.getDate())
const weekday = computed(() => ['日', '一', '二', '三', '四', '五', '六'][now.value.getDay()])
</script>

<style scoped>
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

.numerical-time {
  font-family: 'AFCamberwell-One-Regular', serif;
  font-size: 50vw;
  line-height: 1;
  letter-spacing: -0.8vw;
  background: linear-gradient(to bottom, var(--color-0), var(--color-1));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  color: transparent;
}

.numerical-colon {
  padding: 0 0.3vw;
  vertical-align: 4vw;
}

.numerical-side {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  margin-left: 5vw;
  padding-top: 9vw;
  white-space: nowrap;
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
  min-width: 24vw;
  box-sizing: border-box;
}

.stock-row {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  gap: 2vw;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'SF Pro Text', 'Helvetica Neue', sans-serif;
  font-size: 2.8vw;
  line-height: 1.25;
}

.stock-name {
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
  max-width: 20vw;
}

.stock-price {
  color: #fff;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  text-align: right;
  white-space: nowrap;
}

.stock-row.night-mode .stock-name {
  color: var(--color-0);
  opacity: 0.8;
}

.stock-row.night-mode .stock-price {
  color: var(--color-0);
}
</style>
