<script setup lang="ts">
import Float from './components/Float.vue'
</script>

<template>
  <main id="app">
   <Float></Float>
  </main>
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
}

#app {
  width: 100vw;
  height: 100vh;
  background: #000;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0;
}
</style>