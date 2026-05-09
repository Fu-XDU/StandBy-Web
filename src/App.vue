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

/* 阶跃式防烧屏：仅 1px 位移，每段约 30s 保持不动；无 opacity/缩放 连续插值，降低合成与绘制频率 */
.clock {
  padding-top: 150vh;
  animation: antiBurnStep 120s linear infinite;
}

@keyframes antiBurnStep {
  0%, 24.999% {
    transform: translate(0, 0);
  }
  25%, 49.999% {
    transform: translate(1px, 0);
  }
  50%, 74.999% {
    transform: translate(1px, 1px);
  }
  75%, 99.999% {
    transform: translate(0, 1px);
  }
  100% {
    transform: translate(0, 0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .clock {
    animation: none;
  }
}
</style>