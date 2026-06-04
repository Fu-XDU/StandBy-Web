export const FLOAT_PAGE_ID = 'float' as const

/** 冒号相对统一亮度的系数 */
export const COLON_BRIGHTNESS_FACTOR = 0.98

export interface FloatConfigPayload {
  autoNightMode: boolean
  nightModeRange: string[]
  autoInvisible: boolean
  invisibleRange: string[]
  invisibleDayEnable: boolean
  invisibleDay: boolean[]
  selectedColorIndex: number
  brightness: number
}

/** 远程 PATCH 时可只传需要修改的字段 */
export type FloatConfigPatch = {
  [K in keyof FloatConfigPayload]?: FloatConfigPayload[K]
}

export function clampOpacity(value: number, fallback = 1): number {
  if (!Number.isFinite(value)) return fallback
  return Math.min(1, Math.max(0, value))
}

export function colonOpacityFromBrightness(brightness: number): number {
  return clampOpacity(brightness * COLON_BRIGHTNESS_FACTOR)
}

function floatEqual(a: number, b: number): boolean {
  return Math.abs(a - b) < 1e-6
}

export const FLOAT_UPDATED_AT_KEY = 'standby_updated_at_float'

function equalStringSlice(a: string[], b: string[]): boolean {
  if (a.length !== b.length) return false
  return a.every((v, i) => v === b[i])
}

function equalBoolSlice(a: boolean[], b: boolean[]): boolean {
  if (a.length !== b.length) return false
  return a.every((v, i) => v === b[i])
}

/** 逐字段比较，避免 JSON 字符串比较 */
export function equalFloatConfig(a: FloatConfigPayload, b: FloatConfigPayload): boolean {
  return (
    a.autoNightMode === b.autoNightMode &&
    equalStringSlice(a.nightModeRange, b.nightModeRange) &&
    a.autoInvisible === b.autoInvisible &&
    equalStringSlice(a.invisibleRange, b.invisibleRange) &&
    a.invisibleDayEnable === b.invisibleDayEnable &&
    equalBoolSlice(a.invisibleDay, b.invisibleDay) &&
    a.selectedColorIndex === b.selectedColorIndex &&
    floatEqual(a.brightness, b.brightness)
  )
}

/** 将局部 patch 合并到当前快照（仅更新出现的字段） */
export function applyFloatConfigPatch(
  current: FloatConfigPayload,
  patch: FloatConfigPatch,
): FloatConfigPayload {
  const next = { ...current, invisibleDay: [...current.invisibleDay] }
  if (patch.autoNightMode !== undefined) next.autoNightMode = patch.autoNightMode
  if (patch.nightModeRange !== undefined) next.nightModeRange = [...patch.nightModeRange]
  if (patch.autoInvisible !== undefined) next.autoInvisible = patch.autoInvisible
  if (patch.invisibleRange !== undefined) next.invisibleRange = [...patch.invisibleRange]
  if (patch.invisibleDayEnable !== undefined) next.invisibleDayEnable = patch.invisibleDayEnable
  if (patch.invisibleDay !== undefined) next.invisibleDay = [...patch.invisibleDay]
  if (patch.selectedColorIndex !== undefined) next.selectedColorIndex = patch.selectedColorIndex
  if (patch.brightness !== undefined) next.brightness = clampOpacity(patch.brightness)
  return next
}

export function readFloatUpdatedAt(): number {
  const raw = localStorage.getItem(FLOAT_UPDATED_AT_KEY)
  if (raw === null) return 0
  const n = Number(raw)
  return Number.isFinite(n) && n >= 0 ? n : 0
}

export function writeFloatUpdatedAt(ts: number): void {
  localStorage.setItem(FLOAT_UPDATED_AT_KEY, String(ts))
}
