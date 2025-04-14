<script setup lang="ts">
import Float from './components/Float.vue'
</script>

<template>
  <meta name="apple-mobile-web-app-capable" content="yes">
  <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
  <main id="app">
    <div class="clock">
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

#app {
  width: 100vw;
  height: 100vh;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0;
}

.clock {
  margin-top: 50vh;
  animation: antiBurnMove 60s infinite linear,
  antiBurnOpacity 80s infinite ease-in-out,
  antiBurnTransform 120s infinite ease-in-out;
  will-change: transform, opacity;
}

@keyframes antiBurnMove {
  0% {
    transform: translate(0px, 0px);
  }
  25% {
    transform: translate(1px, 0px);
  }
  50% {
    transform: translate(1px, 1px);
  }
  75% {
    transform: translate(0px, 1px);
  }
  100% {
    transform: translate(0px, 0px);
  }
}

@keyframes antiBurnOpacity {
  0%, 100% {
    opacity: 0.96;
  }
  50% {
    opacity: 0.88;
  }
}

@keyframes antiBurnTransform {
  0%, 100% {
    transform: scale(1) rotate(0deg);
  }
  50% {
    transform: scale(1.01) rotate(0.3deg);
  }
}
</style>