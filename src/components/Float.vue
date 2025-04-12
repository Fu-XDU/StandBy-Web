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
</template>

<script setup lang="ts">
import {useTime} from '@/composables/useTime'

const {digits} = useTime()

const colors = {
  'blue': ['rgba(0, 105, 200, 0.93)', 'rgba(88, 174, 256, 0.93)', 'rgba(196, 207, 218, 0.98)'],
}

// 动态设置颜色 CSS 变量
function setColors(color1: string, color2: string, colonColor: string) {
  const root = document.documentElement
  root.style.setProperty('--color-1', color1)
  root.style.setProperty('--color-2', color2)
  root.style.setProperty('--color-colon', colonColor)
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
  transform: translateY(-3vh) scale(1.1);
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
</style>