import { onMounted, onUnmounted, type Ref } from 'vue'

export interface GestureAxisAvailability {
  up?: boolean
  down?: boolean
  left?: boolean
  right?: boolean
}

export interface GestureOptions {
  onLongPress?: () => void
  onSwipeUp?: () => void
  onSwipeDown?: () => void
  onSwipeLeft?: () => void
  onSwipeRight?: () => void
  onDragX?: (dx: number) => void
  onDragY?: (dy: number) => void
  onDragEnd?: () => void
  getAxisAvailability?: () => GestureAxisAvailability
  isGestureEnabled?: () => boolean
  longPressDelay?: number
  swipeThreshold?: number
  cancelMoveThreshold?: number
  edgeResistanceLimit?: number
  wheelThreshold?: number
  wheelCooldown?: number
  wheelIdleReset?: number
}

type LockedAxis = 'x' | 'y' | null

export function useGesture(el: Ref<HTMLElement | null>, options: GestureOptions) {
  const {
    onLongPress,
    onSwipeUp,
    onSwipeDown,
    onSwipeLeft,
    onSwipeRight,
    onDragX,
    onDragY,
    onDragEnd,
    getAxisAvailability,
    isGestureEnabled,
    longPressDelay = 500,
    swipeThreshold = 40,
    cancelMoveThreshold = 8,
    edgeResistanceLimit = 154,
    wheelThreshold = 70,
    wheelCooldown = 400,
    wheelIdleReset = 120,
  } = options

  let startX = 0
  let startY = 0
  let longPressTimer: number | undefined
  let longPressFired = false
  let dragging = false
  let lockedAxis: LockedAxis = null
  let wheelAccumY = 0
  let wheelResetTimer: number | undefined
  let wheelCooldownUntil = 0
  let mouseTracking = false

  const defaultAvailability = (): GestureAxisAvailability => ({
    up: true,
    down: true,
    left: true,
    right: true,
  })

  const applyEdgeResistance = (delta: number): number => {
    if (delta === 0) return 0
    const abs = Math.abs(delta)
    return Math.sign(delta) * edgeResistanceLimit * abs / (abs + edgeResistanceLimit * 0.55)
  }

  const clampDelta = (dx: number, dy: number) => {
    const avail = getAxisAvailability?.() ?? defaultAvailability()
    let clampedDx = dx
    let clampedDy = dy

    if (dx > 0 && !avail.right) clampedDx = applyEdgeResistance(dx)
    else if (dx < 0 && !avail.left) clampedDx = applyEdgeResistance(dx)

    if (dy > 0 && !avail.down) clampedDy = applyEdgeResistance(dy)
    else if (dy < 0 && !avail.up) clampedDy = applyEdgeResistance(dy)

    return { clampedDx, clampedDy, avail }
  }

  const cancelLongPress = () => {
    if (longPressTimer !== undefined) {
      clearTimeout(longPressTimer)
      longPressTimer = undefined
    }
  }

  const beginPointer = (x: number, y: number) => {
    startX = x
    startY = y
    longPressFired = false
    dragging = false
    lockedAxis = null

    longPressTimer = window.setTimeout(() => {
      longPressTimer = undefined
      longPressFired = true
      onLongPress?.()
    }, longPressDelay)
  }

  const onTouchStart = (e: TouchEvent) => {
    beginPointer(e.touches[0].clientX, e.touches[0].clientY)
  }

  const onTouchMove = (e: TouchEvent) => {
    const t = e.touches[0]
    const dx = t.clientX - startX
    const dy = t.clientY - startY

    if (Math.abs(dx) > cancelMoveThreshold || Math.abs(dy) > cancelMoveThreshold) {
      cancelLongPress()
    }

    if (longPressFired) return

    if (lockedAxis === null && (Math.abs(dx) > cancelMoveThreshold || Math.abs(dy) > cancelMoveThreshold)) {
      lockedAxis = Math.abs(dx) > Math.abs(dy) ? 'x' : 'y'
      dragging = true
      if (lockedAxis === 'x') onDragY?.(0)
      else onDragX?.(0)
    }

    if (lockedAxis === null) return

    if (!dragging) dragging = true

    const { clampedDx, clampedDy } = clampDelta(dx, dy)

    if (lockedAxis === 'x') {
      onDragY?.(0)
      onDragX?.(clampedDx)
    } else {
      onDragX?.(0)
      onDragY?.(clampedDy)
    }
  }

  const onTouchEnd = (e: TouchEvent) => {
    cancelLongPress()
    if (longPressFired) return

    const t = e.changedTouches[0]
    const dx = t.clientX - startX
    const dy = t.clientY - startY

    if (!dragging || lockedAxis === null) {
      onDragEnd?.()
      lockedAxis = null
      dragging = false
      return
    }

    const { avail } = clampDelta(dx, dy)

    if (lockedAxis === 'x') {
      if (Math.abs(dx) >= swipeThreshold) {
        if (dx > 0 && avail.right) onSwipeRight?.()
        else if (dx < 0 && avail.left) onSwipeLeft?.()
        else onDragEnd?.()
      } else {
        onDragEnd?.()
      }
    } else {
      if (Math.abs(dy) >= swipeThreshold) {
        if (dy > 0 && avail.down) onSwipeDown?.()
        else if (dy < 0 && avail.up) onSwipeUp?.()
        else onDragEnd?.()
      } else {
        onDragEnd?.()
      }
    }

    lockedAxis = null
    dragging = false
  }

  const onTouchCancel = () => {
    cancelLongPress()
    if (!longPressFired) {
      onDragEnd?.()
    }
    lockedAxis = null
    dragging = false
  }

  const onMouseMove = (e: MouseEvent) => {
    const dx = e.clientX - startX
    const dy = e.clientY - startY
    if (Math.abs(dx) > cancelMoveThreshold || Math.abs(dy) > cancelMoveThreshold) {
      cancelLongPress()
    }
  }

  const stopMouseTracking = () => {
    if (!mouseTracking) return
    mouseTracking = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  }

  const onMouseUp = () => {
    stopMouseTracking()
    cancelLongPress()
  }

  const onMouseDown = (e: MouseEvent) => {
    if (e.button !== 0) return
    mouseTracking = true
    beginPointer(e.clientX, e.clientY)
    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
  }

  const onContextMenu = (e: MouseEvent) => {
    if (longPressFired) e.preventDefault()
  }

  const resetWheelAccum = () => {
    wheelAccumY = 0
    if (wheelResetTimer !== undefined) {
      clearTimeout(wheelResetTimer)
      wheelResetTimer = undefined
    }
  }

  const onWheel = (e: WheelEvent) => {
    if (Math.abs(e.deltaX) > Math.abs(e.deltaY) * 1.2) return
    if (isGestureEnabled && !isGestureEnabled()) return

    e.preventDefault()

    const now = Date.now()
    if (now < wheelCooldownUntil) return

    wheelAccumY += e.deltaY

    if (wheelResetTimer !== undefined) clearTimeout(wheelResetTimer)
    wheelResetTimer = window.setTimeout(resetWheelAccum, wheelIdleReset)

    const avail = getAxisAvailability?.() ?? defaultAvailability()

    if (wheelAccumY <= -wheelThreshold) {
      resetWheelAccum()
      if (avail.down) {
        onSwipeDown?.()
        wheelCooldownUntil = now + wheelCooldown
      } else {
        onDragY?.(applyEdgeResistance(wheelThreshold))
        onDragEnd?.()
      }
    } else if (wheelAccumY >= wheelThreshold) {
      resetWheelAccum()
      if (avail.up) {
        onSwipeUp?.()
        wheelCooldownUntil = now + wheelCooldown
      } else {
        onDragY?.(applyEdgeResistance(-wheelThreshold))
        onDragEnd?.()
      }
    }
  }

  onMounted(() => {
    const target = el.value
    if (!target) return
    target.addEventListener('touchstart', onTouchStart, { passive: true })
    target.addEventListener('touchmove', onTouchMove, { passive: true })
    target.addEventListener('touchend', onTouchEnd, { passive: true })
    target.addEventListener('touchcancel', onTouchCancel, { passive: true })
    target.addEventListener('mousedown', onMouseDown)
    target.addEventListener('contextmenu', onContextMenu)
    target.addEventListener('wheel', onWheel, { passive: false })
  })

  onUnmounted(() => {
    cancelLongPress()
    resetWheelAccum()
    stopMouseTracking()
    const target = el.value
    if (!target) return
    target.removeEventListener('touchstart', onTouchStart)
    target.removeEventListener('touchmove', onTouchMove)
    target.removeEventListener('touchend', onTouchEnd)
    target.removeEventListener('touchcancel', onTouchCancel)
    target.removeEventListener('mousedown', onMouseDown)
    target.removeEventListener('contextmenu', onContextMenu)
    target.removeEventListener('wheel', onWheel)
  })
}
