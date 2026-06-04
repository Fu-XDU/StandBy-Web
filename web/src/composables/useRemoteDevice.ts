const DEVICE_ID_KEY = 'standby_device_id'

function generateDeviceId(): string {
  const bytes = new Uint8Array(16)
  crypto.getRandomValues(bytes)
  return Array.from(bytes, (b) => b.toString(16).padStart(2, '0')).join('')
}

export function getOrCreateDeviceId(): string {
  let id = localStorage.getItem(DEVICE_ID_KEY)
  if (!id) {
    id = generateDeviceId()
    localStorage.setItem(DEVICE_ID_KEY, id)
  }
  return id
}

import { apiUrl } from '@/utils/apiRoot'

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
