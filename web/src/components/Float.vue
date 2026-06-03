<template>
  <div style="display: flex;justify-content: center;flex-direction: column;align-items: center;padding-bottom: 110vh">
    <div>向下滑动</div>
    <div>Swipe Down</div>
  </div>
  <time v-if="digits.length > 0" class="font">
    <span :class="['digit', 'digit-0']" :style="{ opacity: displayOpacity(brightness) }">
      {{ digits[0] }}
    </span>
    <span :class="['digit', 'digit-1']" :style="{ opacity: displayOpacity(brightness) }">
      {{ digits[1] }}
    </span>
    <span class="colon" :style="{ opacity: displayOpacity(colonOpacity) }">:</span>
    <span :class="['digit', 'digit-2']" :style="{ opacity: displayOpacity(brightness) }">
      {{ digits[2] }}
    </span>
    <span :class="['digit', 'digit-3']" :style="{ opacity: displayOpacity(brightness) }">
      {{ digits[3] }}
    </span>
  </time>

  <div style="display: flex;flex-direction: row;margin-top: 20vh;align-items: center;">
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

const deviceId = getOrCreateDeviceId()

const repoUrl = __REPO_URL__
const commitHash = __COMMIT_HASH__
const buildTime = __BUILD_TIME__

const {digits} = useTime()

const colors = [
  ['rgb(76,151,219)', 'rgba(76,151,219)', 'rgb(6,93,166)', 'rgb(57,151,220)', 'rgb(6,93,166)', 'rgb(57,151,220)'],
  ['rgb(159,133,235)', 'rgba(159,133,235)', 'rgb(109,72,223)', 'rgb(165,138,245)', 'rgb(109,72,223)', 'rgb(165,138,245)'],
  ['rgb(245,109,142)', 'rgba(245,109,142)', 'rgb(210,60,104)', 'rgb(245,109,142)', 'rgb(210,60,104)', 'rgb(245,109,142)'],
  ['rgb(215,103,175)', 'rgb(215,103,175)', 'rgb(154,51,162)', 'rgb(215,103,175)', 'rgb(154,51,162)', 'rgb(215,103,175)'],
  ['rgb(243,209,76)', 'rgba(243,209,76)', 'rgb(172,148,44)', 'rgb(243,209,76)', 'rgb(172,148,44)', 'rgb(243,209,76)'],
  ['rgb(244,128,121)', 'rgba(244,128,121)', 'rgb(244,89,72)', 'rgb(244,128,121)', 'rgb(244,89,72)', 'rgb(244,128,121)'],
  ['rgb(245,155,86)', 'rgb(245,155,86)', 'rgb(225,107,26)', 'rgb(245,155,86)', 'rgb(225,107,26)', 'rgb(245,155,86)'],
  ['rgb(110,213,134)', 'rgba(110,213,134)', 'rgb(45,130,95)', 'rgb(110,213,134)', 'rgb(45,130,95)', 'rgb(110,213,134)'],
  ['rgb(54,200,209)', 'rgb(54,200,209)', 'rgb(4,126,135)', 'rgb(54,200,209)', 'rgb(4,126,135)', 'rgb(54,200,209)'],
  ['rgb(80,156,228)', 'rgb(110,213,134)', 'rgb(5,94,179)', 'rgb(110,213,134)', 'rgb(80,156,228)', 'rgb(45,130,95)'],
  ['rgb(5,94,179)', 'rgb(245,155,86)', 'rgb(80,156,228)', 'rgb(245,155,86)', 'rgb(5,94,179)', 'rgb(225,107,26)'],
  ['rgb(109,72,223)', 'rgb(215,103,175)', 'rgb(165,138,245)', 'rgb(215,103,175)', 'rgb(109,72,223)', 'rgb(154,51,162)'],
  ['rgb(244,128,121)', 'rgb(243,209,76)', 'rgb(244,89,72)', 'rgb(243,209,76)', 'rgb(244,128,121)', 'rgb(172,148,44)'],
  ['rgb(245,109,142)', 'rgb(54,200,209)', 'rgb(210,60,104)', 'rgb(54,200,209)', 'rgb(245,109,142)', 'rgb(4,126,135)'],
]

// [Digit, Colon]
const nightModeColors = ['rgb(148,4,7)', 'rgb(111,26,23)']

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
  if (isNightMode.value) {
    return
  }

  const root = document.documentElement
  root.style.setProperty('--color-0', colors[index][2])
  root.style.setProperty('--color-1', colors[index][3])
  root.style.setProperty('--color-2', colors[index][4])
  root.style.setProperty('--color-3', colors[index][5])
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
  root.style.setProperty('--color-0', nightModeColors[0])
  root.style.setProperty('--color-1', nightModeColors[0])
  root.style.setProperty('--color-2', nightModeColors[0])
  root.style.setProperty('--color-3', nightModeColors[0])

  root.style.setProperty('--color-colon', nightModeColors[1])
}

const checkTimeAndSetNightMode = () => {
  const nightModeStartAt = nightModeRange.value[0]
  const nightModeEndAt = nightModeRange.value[1]
  if (!nightModeStartAt || !nightModeEndAt) {
    return
  }

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
  if (!invisibleStartAt || !invisibleEndAt) {
    return false
  }

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
      ? nowSeconds >= startSeconds && nowSeconds < endSeconds // 不跨天
      : nowSeconds >= startSeconds || nowSeconds < endSeconds // 跨天
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

.font {
  font-family: 'SF-Pro-Rounded-Black', serif;
  font-size: 37vw;
  display: flex;
  justify-content: center;
  align-items: center;
  max-width: 100vw;
  letter-spacing: -5vw;
  padding-right: 5vw;
}

.digit-0 {
  color: var(--color-0);
  transform: rotate(-5deg) scaleY(1.2);
  z-index: 0;
}

.digit-1 {
  color: var(--color-1);
  font-size: 41vw;
  transform: rotate(1deg) scaleY(1.1);
  z-index: 1;
  mix-blend-mode: screen;
}

.colon {
  color: var(--color-colon);
  z-index: 2;
  transform: translateY(-3vw) scale(1.1) rotate(2deg);
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

.digit-2 {
  color: var(--color-2);
  transform: rotate(5deg) scaleY(1.15);
  z-index: 0;
}

.digit-3 {
  color: var(--color-3);
  transform: rotate(-2deg) scaleY(1.15);
  z-index: 1;
  mix-blend-mode: screen;
}

/* 外部大圆 */
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

/* 内部小圆 */
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
  -webkit-overflow-scrolling: touch; /* iOS 惯性滚动 */
  scrollbar-width: none; /* Firefox 隐藏滚动条 */
  padding-right: 0;
}

.color-bar.expanded {
  width: 75vw; /* 展开后的宽度 */
  padding-right: 1vw;
}

.color-bar::-webkit-scrollbar {
  display: none; /* Chrome/Safari 隐藏滚动条 */
}

.color-dot-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px; /* 这个 padding 决定蓝环与彩色圆点的间距 */
  margin: 5px;
  border-radius: 50%;
  transition: box-shadow 0.2s ease;
}

.color-dot-wrapper.selected {
  box-shadow: 0 0 0 3px rgba(11, 132, 235); /* 蓝色外环 */
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

.footer {
  text-align: center; /* 居中 */
  font-size: 12px; /* 小一点 */
  color: #666; /* 灰色 */
  padding: 10px 0;
}
</style>