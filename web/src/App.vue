<script setup lang="ts">
import {computed, markRaw, onMounted, onUnmounted, ref, useTemplateRef, watch} from 'vue'
import ClockLayout from './components/ClockLayout.vue'
import Float, {floatColors, floatNightModeColors} from './components/Float.vue'
import Numerical, {numericalColors, numericalNightModeColors} from './components/Numerical.vue'
import ColorPanel from './components/ColorPanel.vue'
import SettingsSheet from './components/SettingsSheet.vue'
import {useGesture} from '@/composables/useGesture'
import {useClockStyle, CLOCK_TYPES} from '@/composables/useClockStyle'
import {isFullscreenActive, requestAppFullscreen} from '@/utils/fullscreen'

const {clockStyle} = useClockStyle()

// ── 导航状态 ─────────────────────────────────────────────
const typeIndex = ref(0)
const styleIndex = ref(0)

// 当前类型的所有样式
const currentType = computed(() => CLOCK_TYPES[typeIndex.value] ?? CLOCK_TYPES[0])
const maxStyleIndex = computed(() => currentType.value.styles.length - 1)
const maxTypeIndex = computed(() => CLOCK_TYPES.length - 1)

// styleIndex → clockStyle ref 双向同步
watch([typeIndex, styleIndex], () => {
  const style = currentType.value.styles[styleIndex.value]
  if (style) clockStyle.value = style.id
}, {immediate: true})

watch(clockStyle, (id) => {
  for (let ti = 0; ti < CLOCK_TYPES.length; ti++) {
    const si = CLOCK_TYPES[ti].styles.findIndex(s => s.id === id)
    if (si !== -1) {
      typeIndex.value = ti
      styleIndex.value = si
      return
    }
  }
})

// ── 拖拽偏移 ─────────────────────────────────────────────
const dragX = ref(0)
const dragY = ref(0)
const transitioning = ref(false)

let transitionTimer: number | undefined

const commitTransition = () => {
  transitioning.value = true
  dragX.value = 0
  dragY.value = 0
  if (transitionTimer !== undefined) clearTimeout(transitionTimer)
  transitionTimer = window.setTimeout(() => {
    transitioning.value = false
  }, 380)
}

const resetDrag = () => {
  transitioning.value = true
  dragX.value = 0
  dragY.value = 0
  if (transitionTimer !== undefined) clearTimeout(transitionTimer)
  transitionTimer = window.setTimeout(() => {
    transitioning.value = false
  }, 280)
}

const typeTrackStyle = computed(() => ({
  transform: `translateX(calc(${-typeIndex.value * 100}vw + ${dragX.value}px))`,
  transition: transitioning.value ? 'transform 0.35s cubic-bezier(0.32,0.72,0,1)' : 'none',
}))

// 每列独立 translateY（只对当前列生效，其他列固定）
const styleTrackStyleForType = (ti: number) => computed(() => {
  const isActive = ti === typeIndex.value
  return {
    transform: `translateY(calc(${-styleIndex.value * 100}vh + ${isActive ? dragY.value : 0}px))`,
    transition: transitioning.value ? 'transform 0.35s cubic-bezier(0.32,0.72,0,1)' : 'none',
  }
})

// ── 长按 / 颜色面板 ───────────────────────────────────────
const panelVisible = ref(false)
const settingsVisible = ref(false)
const clockShrunk = ref(false)

const openPanel = () => {
  panelVisible.value = true
  clockShrunk.value = true
}

const closePanel = () => {
  panelVisible.value = false
  clockShrunk.value = false
}

const openSettings = () => {
  closePanel()
  settingsVisible.value = true
}

// ── Anti-burn ────────────────────────────────────────────
const antiBurnRef = useTemplateRef<HTMLElement>('antiBurn')

const positions = [
  'translate(0, 0)',
  'translate(1px, 0)',
  'translate(1px, 1px)',
  'translate(0, 1px)',
] as const

let burnStep = 0
let burnInterval: number | undefined

const applyAntiBurn = () => {
  const el = antiBurnRef.value
  if (!el) return
  el.style.transform = positions[burnStep]
  burnStep = (burnStep + 1) % positions.length
}

// ── 手势 ─────────────────────────────────────────────────
const rootEl = useTemplateRef<HTMLElement>('root')

useGesture(rootEl as ReturnType<typeof useTemplateRef<HTMLElement>>, {
  isGestureEnabled: () => !panelVisible.value && !settingsVisible.value,
  getAxisAvailability: () => ({
    up: styleIndex.value < maxStyleIndex.value,
    down: styleIndex.value > 0,
    left: typeIndex.value < maxTypeIndex.value,
    right: typeIndex.value > 0,
  }),
  onLongPress: () => {
    if (!panelVisible.value && !settingsVisible.value) openPanel()
  },
  onSwipeUp: () => {
    if (panelVisible.value || settingsVisible.value) return
    if (styleIndex.value < maxStyleIndex.value) {
      styleIndex.value++
      commitTransition()
    } else {
      resetDrag()
    }
  },
  onSwipeDown: () => {
    if (panelVisible.value || settingsVisible.value) return
    if (styleIndex.value > 0) {
      styleIndex.value--
      commitTransition()
    } else {
      resetDrag()
    }
  },
  onSwipeLeft: () => {
    if (panelVisible.value || settingsVisible.value) return
    if (typeIndex.value < maxTypeIndex.value) {
      typeIndex.value++
      styleIndex.value = 0
      commitTransition()
    } else {
      resetDrag()
    }
  },
  onSwipeRight: () => {
    if (panelVisible.value || settingsVisible.value) return
    if (typeIndex.value > 0) {
      typeIndex.value--
      styleIndex.value = 0
      commitTransition()
    } else {
      resetDrag()
    }
  },
  onDragX: (dx) => {
    if (panelVisible.value || settingsVisible.value) return
    dragX.value = dx
  },
  onDragY: (dy) => {
    if (panelVisible.value || settingsVisible.value) return
    dragY.value = dy
  },
  onDragEnd: () => {
    resetDrag()
  },
})

// ── 全屏 ─────────────────────────────────────────────────
let fullscreenRequested = false
const tryFullscreen = () => {
  if (fullscreenRequested || isFullscreenActive()) return
  fullscreenRequested = true
  void requestAppFullscreen().catch(() => {})
}

onMounted(() => {
  if (!window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    applyAntiBurn()
    burnInterval = window.setInterval(applyAntiBurn, 30_000)
  }
  rootEl.value?.addEventListener('touchstart', tryFullscreen, {once: true, passive: true})
})

onUnmounted(() => {
  if (burnInterval !== undefined) clearInterval(burnInterval)
  if (transitionTimer !== undefined) clearTimeout(transitionTimer)
})

// 把每个 type 的 styleTrackStyle 计算属性提前算好（避免模板里调用函数）
const styleTrackStyles = CLOCK_TYPES.map((_, ti) => styleTrackStyleForType(ti))

// currentStyle 供 ColorPanel 传 colors
const currentStyleDef = computed(() => currentType.value.styles[styleIndex.value] ?? currentType.value.styles[0])
</script>

<template>
  <meta name="apple-mobile-web-app-capable" content="yes">
  <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
  <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">

  <div ref="root" class="standby-root">
    <!-- anti-burn 外层，只负责 translate -->
    <div ref="antiBurn" class="clock-anti-burn">
      <!-- 时钟缩放层，长按时 scale(0.78) -->
      <div class="clock-scale" :class="{ shrunk: clockShrunk }">
        <!-- 类型轨道（左右滑动） -->
        <div class="type-track" :style="typeTrackStyle">
          <div
              v-for="(type, ti) in CLOCK_TYPES"
              :key="type.id"
              class="type-col"
          >
            <!-- 样式轨道（上下滑动） -->
            <div class="style-track" :style="styleTrackStyles[ti].value">
              <div
                  v-for="style in type.styles"
                  :key="style.id"
                  class="clock-cell"
              >
                <ClockLayout
                    :style-id="style.id"
                    :colors="style.colors"
                    :night-mode-colors="style.nightModeColors"
                >
                  <template #default="slotProps">
                    <component :is="style.component" v-bind="slotProps"/>
                  </template>
                </ClockLayout>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 长按后叠加的颜色面板 -->
    <ColorPanel
        v-if="panelVisible"
        :colors="currentStyleDef.colors"
        @close="closePanel"
        @open-settings="openSettings"
    />

    <!-- 设置 Sheet -->
    <SettingsSheet
        v-if="settingsVisible"
        @close="settingsVisible = false"
    />
  </div>
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  background: #000;
  overflow: hidden;
  min-height: 100dvh;
}

html::-webkit-scrollbar, body::-webkit-scrollbar {
  display: none;
}

#app {
  width: 100%;
  height: 100%;
  min-height: 100dvh;
  min-height: -webkit-fill-available;
}

.standby-root {
  position: fixed;
  inset: 0;
  background: #000;
  overflow: hidden;
  touch-action: none;
}

.clock-anti-burn {
  width: 100vw;
  height: 100vh;
  will-change: transform;
}

.clock-scale {
  width: 100%;
  height: 100%;
  transition: transform 0.35s ease;
  transform-origin: center center;
}

.clock-scale.shrunk {
  transform: scale(0.78);
}

.type-track {
  display: flex;
  flex-direction: row;
  width: 100%;
  height: 100%;
  will-change: transform;
}

.type-col {
  width: 100vw;
  height: 100vh;
  flex-shrink: 0;
  overflow: hidden;
}

.style-track {
  display: flex;
  flex-direction: column;
  will-change: transform;
}

.clock-cell {
  width: 100vw;
  height: 100vh;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

@media (prefers-reduced-motion: reduce) {
  .clock-anti-burn {
    transform: none !important;
  }
}
</style>
