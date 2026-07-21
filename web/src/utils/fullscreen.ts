type FullscreenDocument = Document & {
  webkitFullscreenElement?: Element | null
  webkitExitFullscreen?: () => Promise<void> | void
  webkitCancelFullScreen?: () => Promise<void> | void
}

type FullscreenElement = HTMLElement & {
  webkitRequestFullscreen?: (options?: FullscreenOptions) => Promise<void> | void
  webkitRequestFullScreen?: () => Promise<void> | void
}

type NavigatorStandalone = Navigator & {
  standalone?: boolean
}

export function isStandaloneDisplay(): boolean {
  const nav = navigator as NavigatorStandalone
  if (nav.standalone) return true
  return window.matchMedia('(display-mode: standalone)').matches
    || window.matchMedia('(display-mode: fullscreen)').matches
}

export function supportsNativeFullscreen(): boolean {
  if (isStandaloneDisplay()) return false
  const el = document.documentElement as FullscreenElement
  return typeof el.requestFullscreen === 'function'
    || typeof el.webkitRequestFullscreen === 'function'
    || typeof el.webkitRequestFullScreen === 'function'
}

export function getFullscreenElement(): Element | null {
  const doc = document as FullscreenDocument
  return document.fullscreenElement ?? doc.webkitFullscreenElement ?? null
}

export function isFullscreenActive(): boolean {
  return isStandaloneDisplay() || getFullscreenElement() != null
}

export async function requestAppFullscreen(): Promise<void> {
  if (isStandaloneDisplay()) return

  const el = document.documentElement as FullscreenElement
  if (typeof el.requestFullscreen === 'function') {
    await el.requestFullscreen({ navigationUI: 'hide' })
    return
  }
  if (typeof el.webkitRequestFullscreen === 'function') {
    await el.webkitRequestFullscreen()
    return
  }
  if (typeof el.webkitRequestFullScreen === 'function') {
    await el.webkitRequestFullScreen()
    return
  }
  throw new Error('fullscreen_unsupported')
}

export async function exitAppFullscreen(): Promise<void> {
  if (isStandaloneDisplay()) return

  const doc = document as FullscreenDocument
  if (typeof document.exitFullscreen === 'function') {
    await document.exitFullscreen()
    return
  }
  if (typeof doc.webkitExitFullscreen === 'function') {
    await doc.webkitExitFullscreen()
    return
  }
  if (typeof doc.webkitCancelFullScreen === 'function') {
    await doc.webkitCancelFullScreen()
    return
  }
  throw new Error('fullscreen_unsupported')
}

export function onFullscreenChange(handler: () => void): () => void {
  document.addEventListener('fullscreenchange', handler)
  document.addEventListener('webkitfullscreenchange', handler)
  return () => {
    document.removeEventListener('fullscreenchange', handler)
    document.removeEventListener('webkitfullscreenchange', handler)
  }
}
