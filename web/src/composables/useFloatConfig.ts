import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useLocalStorageSync } from '@/composables/useLocalStorageSync'
import { usePageRemoteSync } from '@/composables/usePageRemoteSync'
import { wasOpenedViaDeviceIdLink, getOrCreateDeviceId, registerDevice } from '@/composables/useRemoteDevice'
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
const clockStyle = ref('float')
const colorIndexMap = ref<Record<string, number>>({ float: 0, numerical: 0 })
const brightness = ref(1)
const remoteControlEnabled = ref(false)

const selectedColorIndex = computed({
  get: () => colorIndexMap.value[clockStyle.value] ?? 0,
  set: (val: number) => {
    colorIndexMap.value = { ...colorIndexMap.value, [clockStyle.value]: val }
  },
})

let floatConfigInitialized = false
let applyingRemote = false
let hydrated = false
let scheduleRemotePush: (() => void) | undefined
let cancelRemotePush: (() => void) | undefined
let startRemoteSchedule: (() => void) | undefined
let stopRemoteSchedule: (() => void) | undefined

const isRemoteSyncEnabled = () => remoteControlEnabled.value

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
  clockStyle: clockStyle.value,
  colorIndexMap: { ...colorIndexMap.value },
  brightness: brightness.value,
})

function equalIntMap(a: Record<string, number>, b: Record<string, number>): boolean {
  const keysA = Object.keys(a)
  const keysB = Object.keys(b)
  if (keysA.length !== keysB.length) return false
  return keysA.every((k) => b[k] === a[k])
}

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
  if (clockStyle.value !== config.clockStyle) patch.clockStyle = config.clockStyle
  if (!equalIntMap(colorIndexMap.value, config.colorIndexMap)) {
    patch.colorIndexMap = config.colorIndexMap
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
  if (patch.clockStyle !== undefined) clockStyle.value = patch.clockStyle
  if (patch.colorIndexMap !== undefined) colorIndexMap.value = { ...patch.colorIndexMap }
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

function migrateOldSelectedColorIndex() {
  const oldRaw = localStorage.getItem('selectedColorIndex')
  if (oldRaw === null) return
  try {
    const n = JSON.parse(oldRaw)
    if (typeof n === 'number' && Number.isFinite(n)) {
      colorIndexMap.value = { ...colorIndexMap.value, float: n }
    }
  } catch { /* ignore */ }
  localStorage.removeItem('selectedColorIndex')
}

function migrateOldClockStyle() {
  const old = localStorage.getItem('standby_clock_style')
  if (old === null) return
  clockStyle.value = old
  localStorage.removeItem('standby_clock_style')
}

function initFloatConfig() {
  useLocalStorageSync({ autoNightMode, autoInvisible, invisibleDayEnable })
  useLocalStorageSync({ nightModeRange, invisibleRange }, dateRangeOptions)
  useLocalStorageSync({ invisibleDay }, { deep: true })
  useLocalStorageSync('clockStyle', clockStyle)
  useLocalStorageSync('colorIndexMap', colorIndexMap, {
    parser: (val: unknown) => {
      if (val && typeof val === 'object' && !Array.isArray(val)) return val as Record<string, number>
      return { float: 0, numerical: 0 }
    },
    deep: true,
  })
  useLocalStorageSync('brightness', brightness)
  useLocalStorageSync('remoteControlEnabled', remoteControlEnabled)

  migrateOldSelectedColorIndex()
  migrateOldClockStyle()

  if (wasOpenedViaDeviceIdLink()) {
    remoteControlEnabled.value = true
  }

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
      clockStyle,
      colorIndexMap,
      brightness,
    ],
    () => {
      if (!hydrated || applyingRemote || !remoteControlEnabled.value) return
      bumpUpdatedAt()
      scheduleRemotePush?.()
    },
    { deep: true },
  )

  watch(remoteControlEnabled, (enabled) => {
    if (!hydrated) return
    if (!enabled) {
      cancelRemotePush?.()
      stopRemoteSchedule?.()
      return
    }
    void registerDevice(getOrCreateDeviceId()).then(() => {
      startRemoteSchedule?.()
      scheduleRemotePush?.()
    }).catch((e) => {
      console.warn('[remote] register failed', e)
    })
  })

  usePageRemoteSync(FLOAT_PAGE_ID, toSnapshot, applySnapshot, readFloatUpdatedAt, {
    isSyncEnabled: isRemoteSyncEnabled,
    onPushReady: ({ schedulePush, cancelPush, startSchedule, stopSchedule }) => {
      scheduleRemotePush = schedulePush
      cancelRemotePush = cancelPush
      startRemoteSchedule = startSchedule
      stopRemoteSchedule = stopSchedule
    },
  })
}

export function triggerRemoteSync(): void {
  scheduleRemotePush?.()
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
    clockStyle,
    colorIndexMap,
    selectedColorIndex,
    brightness,
    remoteControlEnabled,
  }
}
