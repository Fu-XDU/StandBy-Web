<template>
  <time v-if="digits.length > 0" class="float-font">
    <template v-for="slot in digitSlots" :key="slot.index">
      <span v-if="slot.index === 2" class="float-colon" :style="{ opacity: displayOpacity(colonOpacity) }">:</span>
      <span class="float-digit-slot" :class="`float-slot-${slot.index}`" :style="{ opacity: displayOpacity(brightness) }">
        <span
          v-for="layer in slot.layers"
          :key="layer.id"
          class="float-digit"
          :class="[
            `float-digit-${slot.index}`,
            {
              'float-digit-enter': layer.phase === 'enter',
              'float-digit-leave': layer.phase === 'leave',
            },
          ]"
          :aria-hidden="layer.phase === 'leave'"
          @animationend="onDigitAnimationEnd(slot.index, layer.id, $event)"
        >{{ layer.char }}</span>
      </span>
    </template>
  </time>
</template>

<script lang="ts">
export const floatColors = [
  ['rgb(6,93,166)', 'rgb(57,151,220)', 'rgb(6,93,166)', 'rgb(57,151,220)', 'rgb(6,93,166)', 'rgb(57,151,220)'],
  ['rgb(109,72,223)', 'rgb(165,138,245)', 'rgb(109,72,223)', 'rgb(165,138,245)', 'rgb(109,72,223)', 'rgb(165,138,245)'],
  ['rgb(210,60,104)', 'rgb(245,109,142)', 'rgb(210,60,104)', 'rgb(245,109,142)', 'rgb(210,60,104)', 'rgb(245,109,142)'],
  ['rgb(154,51,162)', 'rgb(215,103,175)', 'rgb(154,51,162)', 'rgb(215,103,175)', 'rgb(154,51,162)', 'rgb(215,103,175)'],
  ['rgb(172,148,44)', 'rgb(243,209,76)', 'rgb(172,148,44)', 'rgb(243,209,76)', 'rgb(172,148,44)', 'rgb(243,209,76)'],
  ['rgb(244,89,72)', 'rgb(244,128,121)', 'rgb(244,89,72)', 'rgb(244,128,121)', 'rgb(244,89,72)', 'rgb(244,128,121)'],
  ['rgb(225,107,26)', 'rgb(245,155,86)', 'rgb(225,107,26)', 'rgb(245,155,86)', 'rgb(225,107,26)', 'rgb(245,155,86)'],
  ['rgb(45,130,95)', 'rgb(110,213,134)', 'rgb(45,130,95)', 'rgb(110,213,134)', 'rgb(45,130,95)', 'rgb(110,213,134)'],
  ['rgb(4,126,135)', 'rgb(54,200,209)', 'rgb(4,126,135)', 'rgb(54,200,209)', 'rgb(4,126,135)', 'rgb(54,200,209)'],
  ['rgb(5,94,179)', 'rgb(110,213,134)', 'rgb(5,94,179)', 'rgb(110,213,134)', 'rgb(80,156,228)', 'rgb(45,130,95)'],
  ['rgb(80,156,228)', 'rgb(245,155,86)', 'rgb(80,156,228)', 'rgb(245,155,86)', 'rgb(5,94,179)', 'rgb(225,107,26)'],
  ['rgb(165,138,245)', 'rgb(215,103,175)', 'rgb(165,138,245)', 'rgb(215,103,175)', 'rgb(109,72,223)', 'rgb(154,51,162)'],
  ['rgb(244,89,72)', 'rgb(243,209,76)', 'rgb(244,89,72)', 'rgb(243,209,76)', 'rgb(244,128,121)', 'rgb(172,148,44)'],
  ['rgb(210,60,104)', 'rgb(54,200,209)', 'rgb(210,60,104)', 'rgb(54,200,209)', 'rgb(245,109,142)', 'rgb(4,126,135)'],
]

export const floatNightModeColors = ['rgb(148,4,7)', 'rgb(111,26,23)']
</script>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  digits: string[]
  displayOpacity: (opacity: number) => number
  colonOpacity: number
  brightness: number
  isNightMode: boolean
}>()

type DigitPhase = 'settled' | 'enter' | 'leave'

type DigitLayer = {
  id: number
  char: string
  phase: DigitPhase
}

type DigitSlot = {
  index: number
  layers: DigitLayer[]
}

const digitSlots = ref<DigitSlot[]>([0, 1, 2, 3].map((index) => ({ index, layers: [] })))
let digitLayerId = 0

const prefersReducedMotion = () => window.matchMedia('(prefers-reduced-motion: reduce)').matches

// 和数字时钟同一套过渡：旧数字加速上移变浅，新数字从偏低处减速升到原位。
const syncDigitSlots = (digits: string[]) => {
  for (let index = 0; index < 4; index += 1) {
    const char = digits[index]
    if (!char) continue
    const slot = digitSlots.value[index]
    const current = slot.layers.find((layer) => layer.phase !== 'leave')
    if (!current) {
      slot.layers.push({ id: digitLayerId++, char, phase: 'settled' })
      continue
    }
    if (current.char === char) continue
    if (prefersReducedMotion()) {
      slot.layers = [{ id: digitLayerId++, char, phase: 'settled' }]
      continue
    }
    current.phase = 'leave'
    slot.layers.push({ id: digitLayerId++, char, phase: 'enter' })
  }
}

watch(
  () => props.digits,
  (digits) => syncDigitSlots(digits),
  { immediate: true },
)

const onDigitAnimationEnd = (index: number, id: number, event: AnimationEvent) => {
  if (event.target !== event.currentTarget) return
  const slot = digitSlots.value[index]
  const layer = slot?.layers.find((item) => item.id === id)
  if (!layer) return
  if (event.animationName.includes('leave') && layer.phase === 'leave') {
    slot.layers = slot.layers.filter((item) => item.id !== id)
    return
  }
  if (event.animationName.includes('enter') && layer.phase === 'enter') {
    layer.phase = 'settled'
  }
}
</script>

<style scoped>
@property --digit-wash {
  syntax: '<number>';
  inherits: false;
  initial-value: 0;
}

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

.float-digit-slot {
  position: relative;
}

.float-slot-0,
.float-slot-2 {
  z-index: 0;
}

.float-slot-1,
.float-slot-3 {
  z-index: 1;
}

.float-digit {
  --digit-wash: 0;
  display: block;
  color: color-mix(in srgb, var(--digit-color) calc((1 - var(--digit-wash)) * 100%), white);
  /* 样式轨道的 transform 会让 iPad 的 WebKit 把字形从中间切开。单独建一层合成层避开。 */
  transform: translate3d(0, 0, 0);
  -webkit-backface-visibility: hidden;
  backface-visibility: hidden;
}

.float-digit-0 {
  --digit-color: var(--color-0);
  rotate: -5deg;
  scale: 1 1.2;
}

.float-digit-1 {
  --digit-color: var(--color-1);
  font-size: 41vw;
  rotate: 1deg;
  scale: 1 1.1;
  mix-blend-mode: screen;
}

.float-colon {
  color: var(--color-colon);
  z-index: 2;
  transform: translate3d(0, -3vw, 0) scale(1.1) rotate(2deg);
  -webkit-backface-visibility: hidden;
  backface-visibility: hidden;
}

.float-digit-2 {
  --digit-color: var(--color-2);
  rotate: 5deg;
  scale: 1 1.15;
}

.float-digit-3 {
  --digit-color: var(--color-3);
  rotate: -2deg;
  scale: 1 1.15;
  mix-blend-mode: screen;
}

.float-digit-leave {
  position: absolute;
  left: 0;
  top: 0;
  animation: float-digit-leave 0.2s ease-in forwards;
}

.float-digit-enter {
  animation: float-digit-enter 0.7s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes float-digit-leave {
  from {
    translate: 0 0;
    opacity: 1;
    --digit-wash: 0;
  }
  to {
    translate: 0 -1.6vw;
    opacity: 0;
    --digit-wash: 1;
  }
}

@keyframes float-digit-enter {
  from {
    translate: 0 1vw;
    opacity: 0.2;
  }
  to {
    translate: 0 0;
    opacity: 1;
  }
}

@media (prefers-reduced-motion: reduce) {
  .float-digit-enter,
  .float-digit-leave {
    animation: none;
  }
}
</style>
