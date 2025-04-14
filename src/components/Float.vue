<template>
  <time v-if="digits.length > 0" class="font">
    <span :class="['digit', 'digit-0']">
      {{ digits[0] }}
    </span>
    <span :class="['digit', 'digit-1']">
      {{ digits[1] }}
    </span>
    <span class="colon">:</span>
    <span :class="['digit', 'digit-2']">
      {{ digits[2] }}
    </span>
    <span :class="['digit', 'digit-3']">
      {{ digits[3] }}
    </span>
  </time>

  <div style="display: flex;flex-direction: row;margin-top: 20vh;margin-bottom: 15vh;">
    <button class="toggle-button" @click="toggle">
      <span class="inner-dot"/>
    </button>
    <div class="color-bar" :class="{ expanded: isExpanded }">
      <div
          v-for="(color, index) in colors"
          :key="color[0] + color[1]"
          class="color-dot-wrapper"
          :class="{ selected: isExpanded && selectedColorIndex === index }"
          @click="setColors(color, index)"
      >
        <div class="color-dot"
             :style="{ backgroundImage: `linear-gradient(to right, ${color[0]} 50%, ${color[1]} 50%)` }"/>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useTime} from '@/composables/useTime'
import {ref} from 'vue'

const {digits} = useTime()

const colors = [
  ['rgba(82, 162, 237, 1)', 'rgba(82, 162, 237, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(172, 143, 255, 1)', 'rgba(172, 143, 255, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(255, 112, 147, 1)', 'rgba(255, 112, 147, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(224, 105, 182, 1)', 'rgba(224, 105, 182, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(253, 218, 77, 1)', 'rgba(253, 218, 77, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(255, 113, 125, 1)', 'rgba(255, 113, 125, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(114, 222, 138, 1)', 'rgba(114, 222, 138, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(53, 208, 217, 1)', 'rgba(53, 208, 217, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],

  ['rgba(82, 162, 237, 1)', 'rgba(114, 222, 138, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(0, 98, 186, 1)', 'rgba(255, 161, 89, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(113, 74, 232, 1)', 'rgba(224, 105, 182, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(255, 133, 125, 1)', 'rgba(253, 218, 77, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
  ['rgba(255, 132, 147, 1)', 'rgba(53, 208, 217, 1)', 'rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
]

const isExpanded = ref(false)
const toggle = () => {
  isExpanded.value = !isExpanded.value
}

const selectedColorIndex = ref(0)

const setColors = (color: Array<string>, index: number) => {
  selectedColorIndex.value = index

  const root = document.documentElement
  root.style.setProperty('--color-1', color[2])
  root.style.setProperty('--color-2', color[3])
  root.style.setProperty('--color-colon', color[4])
}
</script>

<style>
:root {
  --color-1: rgba(0, 105, 200, 0.93);
  --color-2: rgba(88, 174, 256, 0.93);
  --color-colon: rgba(196, 207, 218, 0.98);
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
  color: var(--color-1);
  transform: rotate(-5deg) scaleY(1.1);
  z-index: 0;
}

.digit-1 {
  color: var(--color-2);
  font-size: 41vw;
  transform: rotate(1deg);
  z-index: 1;
  mix-blend-mode: screen;
}

.colon {
  color: var(--color-colon);
  z-index: 2;
  transform: translateY(-3vh) scale(1.1) rotate(2deg);
  /*backdrop-filter: blur(10px);*/
}

.digit-2 {
  color: var(--color-1);
  transform: rotate(5deg) scaleY(1.1);
  z-index: 0;
}

.digit-3 {
  color: var(--color-2);
  transform: rotate(-2deg) scaleY(1.1);
  z-index: 1;
  mix-blend-mode: screen;
}

/* 外部大圆 */
.toggle-button {
  width: 5vw;
  height: 5vw;
  border-radius: 50%;
  background: rgba(90, 90, 90, 1);
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
  background: rgba(255, 255, 255, 1);
}

.color-bar {
  display: flex;
  overflow: hidden;
  width: 0;
  transition: width 0.2s ease;
  margin-left: 15px;
  flex-direction: row;
  align-items: center;
  padding-left: 5px;
}

.color-bar.expanded {
  width: 70vw; /* 展开后的宽度 */
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
  box-shadow: 0 0 0 3px rgba(11, 132, 235, 1); /* 蓝色外环 */
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
</style>