<template>
  <div class="clock-layout" :style="layoutStyle">
    <slot
        :digits="digits"
        :display-opacity="displayOpacity"
        :colon-opacity="colonOpacity"
        :brightness="brightness"
        :is-night-mode="isNightMode"
    />
  </div>
</template>

<script setup lang="ts">
import {useTime} from '@/composables/useTime'
import {computed, onBeforeUnmount, onMounted, ref} from 'vue'
import {clampOpacity, colonOpacityFromBrightness} from '@/config/floatConfig'
import {useFloatConfig} from '@/composables/useFloatConfig'

const props = defineProps<{
  styleId: string
  colors: string[][]
  nightModeColors: string[]
}>()

const {digits} = useTime()

const isNightMode = ref(false)
const isInvisible = ref(false)

let intervalId: number | undefined
let alignTimeoutId: number | undefined

const stopScheduledTimers = () => {
  if (intervalId !== undefined) {
    clearInterval(intervalId)
    intervalId = undefined
  }
  if (alignTimeoutId !== undefined) {
    clearTimeout(alignTimeoutId)
    alignTimeoutId = undefined
  }
}

const startMinuteAlignedSchedule = () => {
  stopScheduledTimers()
  timerRoutine()
  const now = new Date()
  const msUntilNextMinute = (60 - now.getSeconds()) * 1000 - now.getMilliseconds()
  alignTimeoutId = window.setTimeout(() => {
    alignTimeoutId = undefined
    timerRoutine()
    intervalId = window.setInterval(timerRoutine, 60 * 1000)
  }, msUntilNextMinute)
}

const onVisibilityChange = () => {
  if (document.hidden) {
    stopScheduledTimers()
  } else {
    startMinuteAlignedSchedule()
  }
}

const {
  autoNightMode,
  nightModeRange,
  autoInvisible,
  invisibleRange,
  invisibleDayEnable,
  invisibleDay,
  colorIndexMap,
  brightness,
} = useFloatConfig()

const displayOpacity = (opacity: number) => (isInvisible.value ? 0 : clampOpacity(opacity))

const colonOpacity = computed(() => colonOpacityFromBrightness(brightness.value))

const colorIndex = computed(() => colorIndexMap.value[props.styleId] ?? 0)

const layoutStyle = computed(() => {
  if (isNightMode.value) {
    return {
      '--color-0': props.nightModeColors[0],
      '--color-1': props.nightModeColors[0],
      '--color-2': props.nightModeColors[0],
      '--color-3': props.nightModeColors[0],
      '--color-colon': props.nightModeColors[1],
    }
  }

  const palette = props.colors[colorIndex.value]
  if (!palette) return {}

  return {
    '--color-0': palette[2],
    '--color-1': palette[3],
    '--color-2': palette[4],
    '--color-3': palette[5],
    '--color-colon': 'rgba(196, 207, 218)',
  }
})

const timerRoutine = () => {
  if (autoNightMode.value) {
    checkTimeAndSetNightMode()
  } else {
    setToLightMode()
  }

  let invisible = false
  if (autoInvisible.value) {
    invisible = isInvisibleTime()
  }
  if (invisibleDayEnable.value) {
    invisible ||= isDayInvisible()
  }
  isInvisible.value = invisible
}

const setToLightMode = () => {
  if (!isNightMode.value) return
  isNightMode.value = false
}

const setToNightMode = () => {
  if (isNightMode.value) return
  isNightMode.value = true
}

const checkTimeAndSetNightMode = () => {
  const nightModeStartAt = nightModeRange.value[0]
  const nightModeEndAt = nightModeRange.value[1]
  if (!nightModeStartAt || !nightModeEndAt) return

  const isNight = checkBetweenTime(nightModeStartAt, nightModeEndAt)
  if (isNight) {
    setToNightMode()
  } else {
    setToLightMode()
  }
}

const isDayInvisible = () => {
  return invisibleDay.value[new Date().getDay()]
}

const isInvisibleTime = () => {
  const invisibleStartAt = invisibleRange.value[0]
  const invisibleEndAt = invisibleRange.value[1]
  if (!invisibleStartAt || !invisibleEndAt) return false
  return invisibleDay.value[new Date().getDay()] || checkBetweenTime(invisibleStartAt, invisibleEndAt)
}

onMounted(() => {
  startMinuteAlignedSchedule()
  document.addEventListener('visibilitychange', onVisibilityChange)
})

onBeforeUnmount(() => {
  document.removeEventListener('visibilitychange', onVisibilityChange)
  stopScheduledTimers()
})

const checkBetweenTime = (start: Date, end: Date) => {
  const now = new Date()
  const nowSeconds = now.getHours() * 3600 + now.getMinutes() * 60 + now.getSeconds()
  const startSeconds = start.getHours() * 3600 + start.getMinutes() * 60 + start.getSeconds()
  const endSeconds = end.getHours() * 3600 + end.getMinutes() * 60 + end.getSeconds()
  return startSeconds <= endSeconds
      ? nowSeconds >= startSeconds && nowSeconds < endSeconds
      : nowSeconds >= startSeconds || nowSeconds < endSeconds
}
</script>

<style scoped>
.clock-layout {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  --color-0: rgb(6, 93, 166);
  --color-1: rgb(57, 151, 220);
  --color-2: rgb(6, 93, 166);
  --color-3: rgb(57, 151, 220);
  --color-colon: rgba(196, 207, 218);
}
</style>
