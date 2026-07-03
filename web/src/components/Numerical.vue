<template>
  <div v-if="digits.length > 0" class="numerical-wrapper" :style="{ opacity: displayOpacity(brightness) }">
    <time class="numerical-time">{{ digits[0] }}{{ digits[1] }}<span class="numerical-colon">:</span>{{ digits[2] }}{{ digits[3] }}</time>
    <div class="numerical-side">
      <span class="numerical-day" :style="{ color: isNightMode ? 'var(--color-0)' : '#fff' }">{{ dayOfMonth }} </span><span class="numerical-weekday">周{{ weekday }}</span>
    </div>
  </div>
</template>

<script lang="ts">
export const numericalColors = [
  ['rgb(232,118,102)', 'rgb(242,182,125)', 'rgb(232,118,102)', 'rgb(242,182,125)', 'rgb(232,118,102)', 'rgb(242,182,125)'],
  ['rgb(65,135,225)', 'rgb(135,205,248)', 'rgb(65,135,225)', 'rgb(135,205,248)', 'rgb(65,135,225)', 'rgb(135,205,248)'],
  ['rgb(148,88,228)', 'rgb(198,158,242)', 'rgb(148,88,228)', 'rgb(198,158,242)', 'rgb(148,88,228)', 'rgb(198,158,242)'],
  ['rgb(52,178,102)', 'rgb(138,225,168)', 'rgb(52,178,102)', 'rgb(138,225,168)', 'rgb(52,178,102)', 'rgb(138,225,168)'],
  ['rgb(218,178,48)', 'rgb(245,218,108)', 'rgb(218,178,48)', 'rgb(245,218,108)', 'rgb(218,178,48)', 'rgb(245,218,108)'],
  ['rgb(228,72,132)', 'rgb(242,158,182)', 'rgb(228,72,132)', 'rgb(242,158,182)', 'rgb(228,72,132)', 'rgb(242,158,182)'],
  ['rgb(42,168,188)', 'rgb(102,215,225)', 'rgb(42,168,188)', 'rgb(102,215,225)', 'rgb(42,168,188)', 'rgb(102,215,225)'],
  ['rgb(222,72,62)', 'rgb(242,158,98)', 'rgb(222,72,62)', 'rgb(242,158,98)', 'rgb(222,72,62)', 'rgb(242,158,98)'],
]

export const numericalNightModeColors = ['rgb(148,4,7)', 'rgb(111,26,23)']
</script>

<script setup lang="ts">
import {computed, onBeforeUnmount, onMounted, ref} from 'vue'

defineProps<{
  digits: string[]
  displayOpacity: (opacity: number) => number
  colonOpacity: number
  brightness: number
  isNightMode: boolean
}>()

const now = ref(new Date())
let dateTimer: number | undefined

const updateDate = () => {
  now.value = new Date()
}

onMounted(() => {
  dateTimer = window.setInterval(updateDate, 60_000)
})

onBeforeUnmount(() => {
  if (dateTimer !== undefined) clearInterval(dateTimer)
})

const dayOfMonth = computed(() => now.value.getDate())
const weekday = computed(() => ['日', '一', '二', '三', '四', '五', '六'][now.value.getDay()])
</script>

<style scoped>
.numerical-wrapper {
  display: flex;
  align-items: flex-start;
  justify-content: center;
  max-width: 100vw;
  padding: 0 1vw;
}

.numerical-time {
  font-family: 'AFCamberwell-One-Regular', serif;
  font-size: 50vw;
  line-height: 1;
  letter-spacing: -0.8vw;
  background: linear-gradient(to bottom, var(--color-0), var(--color-1));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  color: transparent;
}

.numerical-colon {
  padding: 0 0.3vw;
  vertical-align: 4vw;
}

.numerical-side {
  display: flex;
  flex-direction: row;
  align-items: baseline;
  margin-left: 5vw;
  padding-top: 9vw;
  white-space: nowrap;
}

.numerical-day {
  font-size: 5vw;
  white-space: nowrap;
  padding-right: 1vw;
}

.numerical-weekday {
  font-size: 5vw;
  color: var(--color-0);
}
</style>
