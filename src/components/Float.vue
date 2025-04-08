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
  'blue': ['#0083F8', '#50A2EC', '#C4CFDA'],
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
  --color-1: #0083F8;
  --color-2: #50A2EC;
  --color-colon: #C4CFDA;
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
  z-index: 3;
  opacity: 0.75;
}

.digit-1 {
  color: var(--color-2);
  font-size: 41vw;
  transform: rotate(1deg);
  z-index: 2;
}

.colon {
  color: var(--color-colon);
  z-index: 2;
  transform: translateY(-3vh) scale(1.1);
  opacity: 0.98;
}

.digit-2 {
  color: var(--color-1);
  transform: rotate(5deg) scaleY(1.1);
  z-index: 1;
  opacity: 0.75;
}

.digit-3 {
  color: var(--color-2);
  transform: rotate(-2deg) scaleY(1.1);
  z-index: 0;
}
</style>