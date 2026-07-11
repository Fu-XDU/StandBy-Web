import { ref } from 'vue'
import { apiUrl } from '@/utils/apiRoot'

const DEVICE_ID_KEY = 'standby_device_id'
const DEVICE_ID_PARAM = 'deviceId'
const DEVICE_ID_FROM_URL_KEY = 'standby_device_id_from_url'

const deviceIdRef = ref('')

function generateDeviceId(): string {
  const bytes = new Uint8Array(16)
  crypto.getRandomValues(bytes)
  return Array.from(bytes, (b) => b.toString(16).padStart(2, '0')).join('')
}

function isValidDeviceId(id: string): boolean {
  return /^[a-f0-9]{32}$/i.test(id)
}

function loadDeviceIdFromStorage(): string {
  let id = localStorage.getItem(DEVICE_ID_KEY)
  if (!id) {
    id = generateDeviceId()
    localStorage.setItem(DEVICE_ID_KEY, id)
  }
  return id
}

/** 从 URL ?deviceId= 读取并覆盖 localStorage，须在应用挂载前调用 */
export function initDeviceIdFromUrl(): void {
  const fromUrl = new URLSearchParams(window.location.search).get(DEVICE_ID_PARAM)?.trim()
  const url = new URL(window.location.href)

  if (fromUrl) {
    if (isValidDeviceId(fromUrl)) {
      localStorage.setItem(DEVICE_ID_KEY, fromUrl)
      deviceIdRef.value = fromUrl
      sessionStorage.setItem(DEVICE_ID_FROM_URL_KEY, '1')
    }
    url.searchParams.delete(DEVICE_ID_PARAM)
    window.history.replaceState(null, '', url.toString())
  }
}

export function wasOpenedViaDeviceIdLink(): boolean {
  return sessionStorage.getItem(DEVICE_ID_FROM_URL_KEY) === '1'
}

export function getOrCreateDeviceId(): string {
  if (!deviceIdRef.value) {
    deviceIdRef.value = loadDeviceIdFromStorage()
  }
  return deviceIdRef.value
}

export function useDeviceId() {
  getOrCreateDeviceId()
  return deviceIdRef
}

export function resetDeviceId(): string {
  const id = generateDeviceId()
  localStorage.setItem(DEVICE_ID_KEY, id)
  deviceIdRef.value = id
  return id
}

export function buildRemoteControlUrl(deviceId: string): string {
  const url = new URL(window.location.href)
  url.searchParams.set(DEVICE_ID_PARAM, deviceId)
  return url.toString()
}

export async function registerDevice(deviceId: string): Promise<void> {
  const res = await fetch(apiUrl('/v1/remote/register'), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ deviceId }),
  })
  if (!res.ok) {
    throw new Error(`register failed: ${res.status}`)
  }
}
