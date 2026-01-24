<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
// App.vue 不再自动检查认证
// 认证检查将在路由守卫中按需进行
</script>

<style>
#app {
  min-height: 100vh;
  background-color: var(--light);
}
</style>
