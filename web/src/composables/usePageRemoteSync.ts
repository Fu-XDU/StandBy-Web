import { onBeforeUnmount, onMounted } from 'vue'
import type { RemotePageId, RemoteSyncResponse } from '@/types/remote'
import { getOrCreateDeviceId, registerDevice } from '@/composables/useRemoteDevice'
import { apiUrl } from '@/utils/apiRoot'

type GetSnapshot<T> = () => T
type ApplySnapshot<T> = (config: T, updatedAt: number) => void
type GetUpdatedAt = () => number

export interface RemotePushHandlers {
  schedulePush: () => void
  cancelPush: () => void
  startSchedule: () => void
  stopSchedule: () => void
}

async function syncPage<T>(
  deviceId: string,
  pageId: RemotePageId,
  updatedAt: number,
  config: T,
): Promise<RemoteSyncResponse<T>> {
  const res = await fetch(apiUrl('/v1/remote/sync'), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ deviceId, pageId, updatedAt, config }),
  })
  if (!res.ok) {
    throw new Error(`sync failed: ${res.status}`)
  }
  return (await res.json()) as RemoteSyncResponse<T>
}

function startMinuteAlignedSchedule(task: () => void): () => void {
  task()
  const now = new Date()
  const msUntilNextMinute = (60 - now.getSeconds()) * 1000 - now.getMilliseconds()
  let alignTimeoutId: number | undefined
  let intervalId: number | undefined

  alignTimeoutId = window.setTimeout(() => {
    alignTimeoutId = undefined
    task()
    intervalId = window.setInterval(task, 60 * 1000)
  }, msUntilNextMinute)

  return () => {
    if (alignTimeoutId !== undefined) clearTimeout(alignTimeoutId)
    if (intervalId !== undefined) clearInterval(intervalId)
  }
}

export function usePageRemoteSync<T>(
  pageId: RemotePageId,
  getSnapshot: GetSnapshot<T>,
  applySnapshot: ApplySnapshot<T>,
  getUpdatedAt: GetUpdatedAt,
  options?: {
    pushDebounceMs?: number
    isSyncEnabled?: () => boolean
    onPushReady?: (handlers: RemotePushHandlers) => void
  },
): void {
  const pushDebounceMs = options?.pushDebounceMs ?? 300
  let stopSchedule: (() => void) | undefined
  let syncing = false
  let pushDebounceTimer: number | undefined
  let pendingPush = false

  const cancelPush = () => {
    if (pushDebounceTimer !== undefined) {
      clearTimeout(pushDebounceTimer)
      pushDebounceTimer = undefined
    }
    pendingPush = false
  }

  const runSync = async () => {
    if (options?.isSyncEnabled && !options.isSyncEnabled()) return
    if (syncing) {
      pendingPush = true
      return
    }
    syncing = true
    const deviceId = getOrCreateDeviceId()
    try {
      const result = await syncPage(deviceId, pageId, getUpdatedAt(), getSnapshot())
      if (result.action === 'apply_server' && result.config !== undefined) {
        applySnapshot(result.config, result.updatedAt)
      } else if (result.updatedAt > getUpdatedAt()) {
        applySnapshot(getSnapshot(), result.updatedAt)
      }
    } catch (e) {
      console.warn(`[remote:${pageId}] sync failed`, e)
    } finally {
      syncing = false
      if (pendingPush) {
        pendingPush = false
        if (!options?.isSyncEnabled || options.isSyncEnabled()) {
          void runSync()
        }
      }
    }
  }

  const scheduledSync = () => {
    if (options?.isSyncEnabled && !options.isSyncEnabled()) return
    void runSync()
  }

  const schedulePush = () => {
    if (options?.isSyncEnabled && !options.isSyncEnabled()) return
    if (pushDebounceTimer !== undefined) clearTimeout(pushDebounceTimer)
    pushDebounceTimer = window.setTimeout(() => {
      pushDebounceTimer = undefined
      void runSync()
    }, pushDebounceMs)
  }

  const startSchedule = () => {
    if (stopSchedule) return
    stopSchedule = startMinuteAlignedSchedule(scheduledSync)
  }

  const stopScheduleFn = () => {
    stopSchedule?.()
    stopSchedule = undefined
  }

  onMounted(() => {
    options?.onPushReady?.({
      schedulePush,
      cancelPush,
      startSchedule,
      stopSchedule: stopScheduleFn,
    })
    if (options?.isSyncEnabled?.()) {
      void registerDevice(getOrCreateDeviceId()).catch((e) => {
        console.warn('[remote] register failed', e)
      })
      startSchedule()
    }
  })

  onBeforeUnmount(() => {
    cancelPush()
    stopScheduleFn()
  })
}
