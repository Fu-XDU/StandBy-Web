<template>
  <time v-if="digits.length > 0" class="float-font">
    <span class="float-digit float-digit-0" :style="{ opacity: displayOpacity(brightness) }">
      {{ digits[0] }}
    </span>
    <span class="float-digit float-digit-1" :style="{ opacity: displayOpacity(brightness) }">
      {{ digits[1] }}
    </span>
    <span class="float-colon" :style="{ opacity: displayOpacity(colonOpacity) }">:</span>
    <span class="float-digit float-digit-2" :style="{ opacity: displayOpacity(brightness) }">
      {{ digits[2] }}
    </span>
    <span class="float-digit float-digit-3" :style="{ opacity: displayOpacity(brightness) }">
      {{ digits[3] }}
    </span>
  </time>
</template>

<script lang="ts">
export const floatColors = [
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

export const floatNightModeColors = ['rgb(148,4,7)', 'rgb(111,26,23)']
</script>

<script setup lang="ts">
defineProps<{
  digits: string[]
  displayOpacity: (opacity: number) => number
  colonOpacity: number
  brightness: number
  isNightMode: boolean
}>()
</script>

<style scoped>
.float-font {
  font-family: 'SF-Pro-Rounded-Black', serif;
  font-size: 37vw;
  display: flex;
  justify-content: center;
  align-items: center;
  max-width: 100vw;
  letter-spacing: -5vw;
  padding-right: 5vw;
}

.float-digit-0 {
  color: var(--color-0);
  transform: rotate(-5deg) scaleY(1.2);
  z-index: 0;
}

.float-digit-1 {
  color: var(--color-1);
  font-size: 41vw;
  transform: rotate(1deg) scaleY(1.1);
  z-index: 1;
  mix-blend-mode: screen;
}

.float-colon {
  color: var(--color-colon);
  z-index: 2;
  transform: translateY(-3vw) scale(1.1) rotate(2deg);
}

.float-digit-2 {
  color: var(--color-2);
  transform: rotate(5deg) scaleY(1.15);
  z-index: 0;
}

.float-digit-3 {
  color: var(--color-3);
  transform: rotate(-2deg) scaleY(1.15);
  z-index: 1;
  mix-blend-mode: screen;
}
</style>
