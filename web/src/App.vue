<script setup lang="ts">
import {onMounted, onUnmounted, useTemplateRef} from 'vue'
import Float from './components/Float.vue'

const clockRef = useTemplateRef('clock')

const positions = [
  'translate(0, 0)',
  'translate(1px, 0)',
  'translate(1px, 1px)',
  'translate(0, 1px)',
] as const

let step = 0
let intervalId: number | undefined

const applyAntiBurn = () => {
  const el = clockRef.value
  if (!el) return
  el.style.transform = positions[step]
  step = (step + 1) % positions.length
}

onMounted(() => {
  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    return
  }
  applyAntiBurn()
  intervalId = window.setInterval(applyAntiBurn, 30_000)
})

onUnmounted(() => {
  if (intervalId !== undefined) {
    clearInterval(intervalId)
  }
})
</script>

<template>
  <meta name="apple-mobile-web-app-capable" content="yes">
  <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
  <main id="app">
    <div ref="clock" class="clock">
      <Float/>
    </div>
  </main>
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  background: #000;
}

html::-webkit-scrollbar, body::-webkit-scrollbar {
  display: none;
}

#app {
  width: 100vw;
  height: 100vh;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0;
}

/* transform 由脚本每 30s 更新，无 CSS animation 时间线 */
.clock {
  padding-top: 150vh;
}

@media (prefers-reduced-motion: reduce) {
  .clock {
    transform: none;
  }
}
</style>