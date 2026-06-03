import { nextTick, onMounted, ref, watch } from 'vue'
import { useLocalStorageSync } from '@/composables/useLocalStorageSync'
import { usePageRemoteSync } from '@/composables/usePageRemoteSync'
import {
  clampOpacity,
  equalFloatConfig,
  FLOAT_PAGE_ID,
  type FloatConfigPatch,
  type FloatConfigPayload,
  readFloatUpdatedAt,
  writeFloatUpdatedAt,
} from '@/config/floatConfig'

const autoNightMode = ref(false)
const nightModeRange = ref([new Date(2025, 0, 0, 0, 0), new Date(2025, 0, 0, 6, 0)])
const autoInvisible = ref(false)
const invisibleRange = ref([new Date(2025, 0, 0, 0, 0), new Date(2025, 0, 0, 6, 0)])
const invisibleDayEnable = ref(false)
const invisibleDay = ref([false, false, false, false, false, false, false])
const selectedColorIndex = ref(0)
const brightness = ref(1)

let floatConfigInitialized = false
let applyingRemote = false
let hydrated = false

const dateRangeOptions = {
  parser: (arr: string[]) => arr.map((s) => new Date(s)),
  deep: true as const,
}

const toSnapshot = (): FloatConfigPayload => ({
  autoNightMode: autoNightMode.value,
  nightModeRange: nightModeRange.value.map((d) => d.toISOString()),
  autoInvisible: autoInvisible.value,
  invisibleRange: invisibleRange.value.map((d) => d.toISOString()),
  invisibleDayEnable: invisibleDayEnable.value,
  invisibleDay: [...invisibleDay.value],
  selectedColorIndex: selectedColorIndex.value,
  brightness: brightness.value,
})

const applyFieldByField = (config: FloatConfigPayload) => {
  const patch: FloatConfigPatch = {}
  if (autoNightMode.value !== config.autoNightMode) patch.autoNightMode = config.autoNightMode
  if (
    config.nightModeRange.length === 2 &&
    (nightModeRange.value[0]?.toISOString() !== config.nightModeRange[0] ||
      nightModeRange.value[1]?.toISOString() !== config.nightModeRange[1])
  ) {
    patch.nightModeRange = config.nightModeRange
  }
  if (autoInvisible.value !== config.autoInvisible) patch.autoInvisible = config.autoInvisible
  if (
    config.invisibleRange.length === 2 &&
    (invisibleRange.value[0]?.toISOString() !== config.invisibleRange[0] ||
      invisibleRange.value[1]?.toISOString() !== config.invisibleRange[1])
  ) {
    patch.invisibleRange = config.invisibleRange
  }
  if (invisibleDayEnable.value !== config.invisibleDayEnable) {
    patch.invisibleDayEnable = config.invisibleDayEnable
  }
  const dayChanged = config.invisibleDay.some((v, i) => invisibleDay.value[i] !== v)
  if (dayChanged) patch.invisibleDay = config.invisibleDay
  if (selectedColorIndex.value !== config.selectedColorIndex) {
    patch.selectedColorIndex = config.selectedColorIndex
  }
  const b = clampOpacity(config.brightness)
  if (Math.abs(brightness.value - b) >= 1e-6) patch.brightness = b
  if (Object.keys(patch).length === 0) return
  applyPatchToRefs(patch)
}

const applyPatchToRefs = (patch: FloatConfigPatch) => {
  if (patch.autoNightMode !== undefined) autoNightMode.value = patch.autoNightMode
  if (patch.nightModeRange !== undefined && patch.nightModeRange.length === 2) {
    nightModeRange.value = patch.nightModeRange.map((s) => new Date(s))
  }
  if (patch.autoInvisible !== undefined) autoInvisible.value = patch.autoInvisible
  if (patch.invisibleRange !== undefined && patch.invisibleRange.length === 2) {
    invisibleRange.value = patch.invisibleRange.map((s) => new Date(s))
  }
  if (patch.invisibleDayEnable !== undefined) invisibleDayEnable.value = patch.invisibleDayEnable
  if (patch.invisibleDay !== undefined) {
    for (let i = 0; i < 7; i++) invisibleDay.value[i] = patch.invisibleDay[i]
  }
  if (patch.selectedColorIndex !== undefined) selectedColorIndex.value = patch.selectedColorIndex
  if (patch.brightness !== undefined) brightness.value = clampOpacity(patch.brightness)
}

const applySnapshot = (config: FloatConfigPayload, updatedAt: number) => {
  const current = toSnapshot()
  if (equalFloatConfig(current, config) && readFloatUpdatedAt() === updatedAt) {
    return
  }
  applyingRemote = true
  try {
    if (!equalFloatConfig(current, config)) {
      applyFieldByField(config)
    }
    writeFloatUpdatedAt(updatedAt)
  } finally {
    applyingRemote = false
  }
}

const bumpUpdatedAt = () => {
  writeFloatUpdatedAt(Date.now())
}

function initFloatConfig() {
  useLocalStorageSync({ autoNightMode, autoInvisible, invisibleDayEnable })
  useLocalStorageSync({ nightModeRange, invisibleRange }, dateRangeOptions)
  useLocalStorageSync({ invisibleDay }, { deep: true })
  useLocalStorageSync('selectedColorIndex', selectedColorIndex, {
    parser: (val: unknown) => {
      if (typeof val === 'number' && Number.isFinite(val)) return val
      const n = parseInt(String(val), 10)
      return Number.isFinite(n) ? n : 0
    },
  })
  useLocalStorageSync('brightness', brightness)

  onMounted(() => {
    nextTick(() => {
      hydrated = true
    })
  })

  watch(
    [
      autoNightMode,
      nightModeRange,
      autoInvisible,
      invisibleRange,
      invisibleDayEnable,
      invisibleDay,
      selectedColorIndex,
      brightness,
    ],
    () => {
      if (!hydrated || applyingRemote) return
      bumpUpdatedAt()
    },
    { deep: true },
  )

  usePageRemoteSync(FLOAT_PAGE_ID, toSnapshot, applySnapshot, readFloatUpdatedAt)
}

export function useFloatConfig() {
  if (!floatConfigInitialized) {
    floatConfigInitialized = true
    initFloatConfig()
  }

  return {
    autoNightMode,
    nightModeRange,
    autoInvisible,
    invisibleRange,
    invisibleDayEnable,
    invisibleDay,
    selectedColorIndex,
    brightness,
  }
}
