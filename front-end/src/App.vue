<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()

onMounted(() => {
  // 检查用户是否已登录
  const token = localStorage.getItem('token')
  if (token) {
    authStore.checkAuth()
  }
})
</script>

<style>
#app {
  min-height: 100vh;
  background-color: var(--light);
}
</style>
