<script setup lang="ts">
import Float from './components/Float.vue'
</script>

<template>
  <meta name="apple-mobile-web-app-capable" content="yes">
  <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
  <main id="app">
   <Float/>
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
</style>