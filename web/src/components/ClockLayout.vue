<template>
  <div style="display: flex;justify-content: center;flex-direction: column;align-items: center;padding-bottom: 110vh">
    <div>向下滑动</div>
    <div>Swipe Down</div>
  </div>

  <slot
      :digits="digits"
      :display-opacity="displayOpacity"
      :colon-opacity="colonOpacity"
      :brightness="brightness"
      :is-night-mode="isNightMode"
  />

  <div style="display: flex;flex-direction: row;margin-top: 20vh;align-items: center;">
    <div class="style-selector">
      <button
          v-for="style in clockStyles"
          :key="style.id"
          :class="['style-btn', { active: clockStyle === style.id }]"
          @click="clockStyle = style.id"
      >
        {{ style.name }}
      </button>
    </div>
  </div>
  <div style="display: flex;flex-direction: row;margin-top: 1.5vh;align-items: center;">
    <button class="toggle-button" @click="toggleExpanded">
      <span class="inner-dot"/>
    </button>
    <div class="color-bar" :class="{ expanded: isExpanded }">
      <div
          v-for="(color, index) in colors"
          :key="color[0] + color[1]"
          class="color-dot-wrapper"
          :class="{ selected: isExpanded && selectedColorIndex === index }"
          @click="setColors(index)"
      >
        <div class="color-dot"
             :style="{ backgroundImage: `linear-gradient(to right, ${color[0]} 50%, ${color[1]} 50%)` }"/>
      </div>
    </div>
  </div>
  <div class="brightness-panel">
    <div class="brightness-row">
      <span class="brightness-label">亮度</span>
      <el-slider
          v-model="brightnessPercent"
          :min="0"
          :max="100"
          :show-tooltip="true"
          class="brightness-slider"
      />
      <span class="brightness-value">{{ brightnessPercent }}%</span>
    </div>
  </div>
  <div style="display: flex; flex-direction: row; align-items: center; padding-top: 1vh">
    <div style="padding-right: 0.4vw;">定时夜间模式</div>
    <el-switch
        v-model="autoNightMode"
        class="ml-2"
        style="--el-switch-on-color: #ff4949; --el-switch-off-color: #4C4D4F"
    />
    <div style="display: flex; align-items: center; padding-left: 1.5vh" v-if="autoNightMode">
      <div style="padding-right: 0.5vw;">起止时间：</div>
      <el-time-picker
          v-model="nightModeRange[0]"
          size="default"
          arrow-control
          placeholder="开始时间"
          style="width: 16vw;"
          :clearable="false"
      />
      <div style="padding: 0 5px">To</div>
      <el-time-picker
          v-model="nightModeRange[1]"
          size="default"
          arrow-control
          placeholder="结束时间"
          style="width: 16vw;"
          :clearable="false"
      />
    </div>
  </div>
  <div style="display: flex; flex-direction: row; align-items: center; padding-top: 1vh">
    <div style="padding-right: 0.4vw;">定时关闭显示</div>
    <el-switch
        v-model="autoInvisible"
        class="ml-2"
        style="--el-switch-on-color: #ff4949; --el-switch-off-color: #4C4D4F"
    />
    <div style="display: flex; align-items: center; padding-left: 1.5vh" v-if="autoInvisible">
      <div style="padding-right: 0.5vw;">起止时间：</div>
      <el-time-picker
          v-model="invisibleRange[0]"
          size="default"
          placeholder="开始时间"
          style="width: 16vw;"
          :clearable="false"
      />
      <div style="padding: 0 5px">To</div>
      <el-time-picker
          v-model="invisibleRange[1]"
          size="default"
          placeholder="结束时间"
          style="width: 16vw;"
          :clearable="false"
      />
    </div>
  </div>
  <div style="display: flex; flex-direction: row; align-items: center; padding-top: 1vh">
    <div style="padding-right: 0.4vw;">全天关闭显示</div>
    <el-switch
        v-model="invisibleDayEnable"
        class="ml-2"
        style="--el-switch-on-color: #ff4949; --el-switch-off-color: #4C4D4F"
    />
    <div style="padding-left: 20px" v-if="invisibleDayEnable">
      <el-checkbox-button
          v-for="(_, index) in invisibleDay"
          :key="index"
          v-model="invisibleDay[index]"
          :label="'周'+weekLabels[index]"
          size="large"
      />
    </div>
  </div>
  <div class="device-id-row">
    <span class="device-id-label">设备 ID</span>
    <code class="device-id-value" :title="deviceId">{{ deviceId }}</code>
  </div>

  <div style="padding-bottom: 10vh"></div>
  <footer class="footer">
    <a :href="`${repoUrl}/commit/${commitHash}`" target="_blank" rel="noopener noreferrer">
      Commit: {{ commitHash }}
    </a>
    <p>Build time: {{ buildTime }}</p>
  </footer>
</template>

<script setup lang="ts">
import {useTime} from '@/composables/useTime'
import {computed, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import {clampOpacity, colonOpacityFromBrightness} from '@/config/floatConfig'
import {useFloatConfig} from '@/composables/useFloatConfig'
import {getOrCreateDeviceId} from '@/composables/useRemoteDevice'
import {useClockStyle} from '@/composables/useClockStyle'

const props = defineProps<{
  colors: string[][]
  nightModeColors: string[]
}>()

const deviceId = getOrCreateDeviceId()
const {clockStyle, clockStyles} = useClockStyle()

const repoUrl = __REPO_URL__
const commitHash = __COMMIT_HASH__
const buildTime = __BUILD_TIME__

const {digits} = useTime()

const isExpanded = ref(false)
const toggleExpanded = () => {
  isExpanded.value = !isExpanded.value
}

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
  selectedColorIndex,
  brightness,
} = useFloatConfig()

const weekLabels = ref(['日', '一', '二', '三', '四', '五', '六'])

const displayOpacity = (opacity: number) => (isInvisible.value ? 0 : clampOpacity(opacity))

const colonOpacity = computed(() => colonOpacityFromBrightness(brightness.value))

const brightnessPercent = computed({
  get: () => Math.round(clampOpacity(brightness.value) * 100),
  set: (percent: number) => {
    brightness.value = clampOpacity(percent / 100)
  },
})

const setColors = (index: number) => {
  selectedColorIndex.value = index
  if (isNightMode.value) return

  const root = document.documentElement
  root.style.setProperty('--color-0', props.colors[index][2])
  root.style.setProperty('--color-1', props.colors[index][3])
  root.style.setProperty('--color-2', props.colors[index][4])
  root.style.setProperty('--color-3', props.colors[index][5])
  root.style.setProperty('--color-colon', 'rgba(196, 207, 218)')
}

const timerRoutine = () => {
  if (autoNightMode.value) {
    checkTimeAndSetNightMode()
  } else {
    setToLightMode()
  }

  let invisible = false;
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
  setColors(selectedColorIndex.value)
}

const setToNightMode = () => {
  if (isNightMode.value) return
  isNightMode.value = true

  const root = document.documentElement
  root.style.setProperty('--color-0', props.nightModeColors[0])
  root.style.setProperty('--color-1', props.nightModeColors[0])
  root.style.setProperty('--color-2', props.nightModeColors[0])
  root.style.setProperty('--color-3', props.nightModeColors[0])
  root.style.setProperty('--color-colon', props.nightModeColors[1])
}

const checkTimeAndSetNightMode = () => {
  const nightModeStartAt = nightModeRange.value[0]
  const nightModeEndAt = nightModeRange.value[1]
  if (!nightModeStartAt || !nightModeEndAt) return

  const isNight = checkBetweenTime(nightModeStartAt, nightModeEndAt)
  if (isNight) {
    setToNightMode();
  } else {
    setToLightMode();
  }
}

const isDayInvisible = () => {
  return invisibleDay.value[new Date().getDay()]
}

const isInvisibleTime = () => {
  const invisibleStartAt = invisibleRange.value[0]
  const invisibleEndAt = invisibleRange.value[1]
  if (!invisibleStartAt || !invisibleEndAt) return false

  return invisibleDay.value[new Date().getDay()] || checkBetweenTime(invisibleStartAt, invisibleEndAt);
}

watch(selectedColorIndex, (index) => {
  setColors(index)
}, {immediate: true})

onMounted(() => {
  startMinuteAlignedSchedule()
  document.addEventListener('visibilitychange', onVisibilityChange)
})

onBeforeUnmount(() => {
  document.removeEventListener('visibilitychange', onVisibilityChange)
  stopScheduledTimers()
})

const checkBetweenTime = (start: Date, end: Date) => {
  const now = new Date();
  const nowSeconds = now.getHours() * 3600 + now.getMinutes() * 60 + now.getSeconds();

  const startSeconds = start.getHours() * 3600 + start.getMinutes() * 60 + start.getSeconds();
  const endSeconds = end.getHours() * 3600 + end.getMinutes() * 60 + end.getSeconds();

  return startSeconds <= endSeconds
      ? nowSeconds >= startSeconds && nowSeconds < endSeconds
      : nowSeconds >= startSeconds || nowSeconds < endSeconds
}
</script>

<style>
:root {
  --color-0: rgb(6, 93, 166);
  --color-1: rgb(57, 151, 220);
  --color-2: rgb(6, 93, 166);
  --color-3: rgb(57, 151, 220);
  --color-colon: rgba(196, 207, 218);
}

.brightness-panel {
  width: min(92vw, 520px);
  margin-top: 2vh;
  padding: 1vh 2vw;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.06);
}

.brightness-row {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 8px;
  margin-bottom: 0.6vh;
}

.brightness-label {
  width: 3.2em;
  font-size: 13px;
  color: #ccc;
  flex-shrink: 0;
}

.brightness-slider {
  flex: 1;
}

.brightness-value {
  width: 3em;
  font-size: 12px;
  color: #888;
  text-align: right;
  flex-shrink: 0;
}

.toggle-button {
  width: 5vw;
  height: 5vw;
  border-radius: 50%;
  background: rgba(90, 90, 90);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  position: relative;
  box-shadow: 0 0 6px rgba(0, 0, 0, 0.3);
}

.toggle-button .inner-dot {
  width: 2.7vw;
  height: 2.7vw;
  border-radius: 50%;
  background: rgba(255, 255, 255);
}

.color-bar {
  display: flex;
  overflow-x: auto;
  width: 0;
  transition: width 0.2s ease;
  margin-left: 15px;
  flex-direction: row;
  align-items: center;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
  padding-right: 0;
}

.color-bar.expanded {
  width: 75vw;
  padding-right: 1vw;
}

.color-bar::-webkit-scrollbar {
  display: none;
}

.color-dot-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  margin: 5px;
  border-radius: 50%;
  transition: box-shadow 0.2s ease;
}

.color-dot-wrapper.selected {
  box-shadow: 0 0 0 3px rgba(11, 132, 235);
}

.color-dot-wrapper:not(.selected) .color-dot:hover {
  transform: scale(1.1);
}

.color-dot {
  width: 3vw;
  height: 3vw;
  border-radius: 50%;
  background-clip: padding-box;
  box-shadow: inset 0 0 0 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.2s ease;
}

.device-id-row {
  display: flex;
  flex-direction: row;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.6em;
  padding-top: 1vh;
  max-width: 92vw;
}

.device-id-label {
  font-size: 13px;
  color: #aaa;
  flex-shrink: 0;
}

.device-id-value {
  font-size: 12px;
  color: #8ab4f8;
  word-break: break-all;
  user-select: all;
}

.style-selector {
  display: flex;
  gap: 8px;
  margin-right: 12px;
}

.style-btn {
  padding: 4px 14px;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: transparent;
  color: rgba(255, 255, 255, 0.5);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.style-btn.active {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
  border-color: rgba(255, 255, 255, 0.4);
}

.footer {
  text-align: center;
  font-size: 12px;
  color: #666;
  padding: 10px 0;
}
</style>
