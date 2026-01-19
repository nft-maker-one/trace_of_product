<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div class="login-page">
    <component :is="currentComponent" @switch-to-login="showLogin" @switch-to-register="showRegister" />
  </div>
</template>

<script setup>
import { ref, markRaw } from 'vue'
import LoginForm from '../components/LoginForm.vue'
import RegisterForm from '../components/RegisterForm.vue'

const currentComponent = ref(markRaw(LoginForm))

function showLogin() {
  currentComponent.value = markRaw(LoginForm)
}

function showRegister() {
  currentComponent.value = markRaw(RegisterForm)
}
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

/* 添加切换动画 */
.component-enter-active,
.component-leave-active {
  transition: all 0.3s ease;
}

.component-enter-from,
.component-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
