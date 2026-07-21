<template>
  <div class="color-panel-overlay" @click.self="$emit('close')">
    <button class="settings-btn" @click="$emit('open-settings')" aria-label="设置">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
           stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="3"/>
        <path
            d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
      </svg>
    </button>

    <div class="panel-content" :class="{ visible: contentVisible }">
      <div class="brightness-row">
        <svg class="brightness-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
          <circle cx="12" cy="12" r="5"/>
          <line x1="12" y1="1" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <line x1="12" y1="21" x2="12" y2="23" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" stroke="currentColor" stroke-width="2"
                stroke-linecap="round"/>
          <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" stroke="currentColor" stroke-width="2"
                stroke-linecap="round"/>
          <line x1="1" y1="12" x2="3" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <line x1="21" y1="12" x2="23" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" stroke="currentColor" stroke-width="2"
                stroke-linecap="round"/>
          <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" stroke="currentColor" stroke-width="2"
                stroke-linecap="round"/>
        </svg>
        <el-slider
            v-model="brightnessPercent"
            :min="0"
            :max="100"
            :show-tooltip="false"
            class="brightness-slider"
            @input="resetAutoClose"
        />
      </div>

      <div class="color-row-wrap">
        <div class="color-row">
          <div
              v-for="(color, index) in colors"
              :key="color[0] + color[1]"
              class="color-dot-wrapper"
              :class="{ selected: selectedColorIndex === index }"
              @click="selectColor(index)"
          >
            <div class="color-dot"
                 :style="{ backgroundImage: `linear-gradient(to right, ${color[0]} 50%, ${color[1]} 50%)` }"/>
          </div>
        </div>
        <button
            class="fullscreen-btn"
            type="button"
            :aria-label="isFullscreen ? '退出全屏' : '全屏'"
            @pointerdown.stop
            @touchend.stop.prevent="onFullscreenTap"
            @click.stop.prevent="onFullscreenTap"
        >
          <!-- enter fullscreen -->
          <svg v-if="!isFullscreen" width="18" height="18" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="15 3 21 3 21 9"/>
            <polyline points="9 21 3 21 3 15"/>
            <line x1="21" y1="3" x2="14" y2="10"/>
            <line x1="3" y1="21" x2="10" y2="14"/>
          </svg>
          <!-- exit fullscreen -->
          <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="4 14 10 14 10 20"/>
            <polyline points="20 10 14 10 14 4"/>
            <line x1="14" y1="10" x2="21" y2="3"/>
            <line x1="3" y1="21" x2="10" y2="14"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, onBeforeUnmount, onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {clampOpacity} from '@/config/floatConfig'
import {useFloatConfig} from '@/composables/useFloatConfig'
import {
  exitAppFullscreen,
  isFullscreenActive,
  isStandaloneDisplay,
  onFullscreenChange,
  requestAppFullscreen,
  supportsNativeFullscreen,
} from '@/utils/fullscreen'

defineProps<{
  colors: string[][]
}>()

const emit = defineEmits<{
  close: []
  'open-settings': []
}>()

const {selectedColorIndex, brightness} = useFloatConfig()

const contentVisible = ref(false)
const isFullscreen = ref(isFullscreenActive())
let autoCloseTimer: number | undefined
let stopFullscreenListen: (() => void) | undefined
let lastFullscreenTapAt = 0

const resetAutoClose = () => {
  if (autoCloseTimer !== undefined) clearTimeout(autoCloseTimer)
  autoCloseTimer = window.setTimeout(() => emit('close'), 5000)
}

const syncFullscreenState = () => {
  isFullscreen.value = isFullscreenActive()
}

const showHomeScreenGuide = async () => {
  await ElMessageBox.alert(
    'iPhone / iPad 的浏览器（含 Chrome）无法网页全屏。\n\n请点击底部分享按钮 →「添加到主屏幕」，再从主屏幕打开，即可全屏显示时钟。',
    '如何全屏显示',
    {
      confirmButtonText: '知道了',
      customClass: 'fullscreen-guide-box',
    },
  )
}

const toggleFullscreen = async () => {
  resetAutoClose()

  if (isStandaloneDisplay()) {
    ElMessage.success('当前已是全屏显示')
    syncFullscreenState()
    return
  }

  if (!supportsNativeFullscreen()) {
    await showHomeScreenGuide()
    return
  }

  try {
    if (isFullscreenActive()) {
      await exitAppFullscreen()
    } else {
      await requestAppFullscreen()
    }
  } catch {
    await showHomeScreenGuide()
  } finally {
    syncFullscreenState()
  }
}

const onFullscreenTap = () => {
  const now = Date.now()
  if (now - lastFullscreenTapAt < 400) return
  lastFullscreenTapAt = now
  void toggleFullscreen()
}

onMounted(() => {
  requestAnimationFrame(() => {
    contentVisible.value = true
  })
  resetAutoClose()
  stopFullscreenListen = onFullscreenChange(syncFullscreenState)
})

onBeforeUnmount(() => {
  if (autoCloseTimer !== undefined) clearTimeout(autoCloseTimer)
  stopFullscreenListen?.()
})

const brightnessPercent = computed({
  get: () => Math.round(clampOpacity(brightness.value) * 100),
  set: (percent: number) => {
    brightness.value = clampOpacity(percent / 100)
  },
})

const selectColor = (index: number) => {
  selectedColorIndex.value = index
  resetAutoClose()
}
</script>

<style scoped>
.color-panel-overlay {
  position: absolute;
  inset: 0;
  z-index: 10;
}

.settings-btn {
  position: absolute;
  top: 5%;
  right: 5%;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: rgba(40, 40, 40, 0.85);
  border: 1px solid rgba(255, 255, 255, 0.15);
  color: rgba(255, 255, 255, 0.85);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

.panel-content {
  position: absolute;
  bottom: 12%;
  left: 50%;
  transform: translateX(-50%) translateY(16px);
  opacity: 0;
  transition: opacity 0.3s ease, transform 0.3s ease;
  width: min(88vw, 480px);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.panel-content.visible {
  opacity: 1;
  transform: translateX(-50%) translateY(0);
}

.brightness-row {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  background: rgba(30, 30, 30, 0.80);
  border-radius: 20px;
  padding: 10px 16px;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  box-sizing: border-box;
}

.brightness-icon {
  color: rgba(255, 255, 255, 0.7);
  flex-shrink: 0;
}

.brightness-slider {
  flex: 1;
}

.color-row-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}

.color-row {
  display: flex;
  overflow-x: auto;
  gap: 4px;
  padding: 10px 12px;
  background: rgba(30, 30, 30, 0.80);
  border-radius: 24px;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  scrollbar-width: none;
  -webkit-overflow-scrolling: touch;
  flex: 1;
  min-width: 0;
  box-sizing: border-box;
}

.color-row::-webkit-scrollbar {
  display: none;
}

.fullscreen-btn {
  flex-shrink: 0;
  width: max(40px, 8vw);
  height: max(40px, 8vw);
  border-radius: 50%;
  background: rgba(30, 30, 30, 0.80);
  border: 1px solid rgba(255, 255, 255, 0.12);
  color: rgba(255, 255, 255, 0.85);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  padding: 0;
}

.color-dot-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 3px;
  border-radius: 50%;
  flex-shrink: 0;
  transition: box-shadow 0.2s ease;
}

.color-dot-wrapper.selected {
  box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.9);
}

.color-dot {
  width: max(28px, 5.5vw);
  height: max(28px, 5.5vw);
  border-radius: 50%;
  cursor: pointer;
  box-shadow: inset 0 0 0 3px rgba(0, 0, 0, 0.12);
}
</style>
