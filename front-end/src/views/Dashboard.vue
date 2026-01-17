<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <MainLayout>
    <template v-slot:content>
      <PersonalCenter v-if="currentTab === 'personal-center'" />
      <DataUpload v-else-if="currentTab === 'data-upload'" />
      <DataQuery v-else-if="currentTab === 'data-query'" />
      <NodeManagement v-else-if="currentTab === 'node-management'" />
      <LogMonitor v-else-if="currentTab === 'log-monitor'" />
    </template>
  </MainLayout>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MainLayout from '../components/MainLayout.vue'
import PersonalCenter from '../components/PersonalCenter.vue'
import DataUpload from '../components/DataUpload.vue'
import DataQuery from '../components/DataQuery.vue'
import NodeManagement from '../components/NodeManagement.vue'
import LogMonitor from '../components/LogMonitor.vue'

const route = useRoute()
const router = useRouter()

const currentTab = computed(() => {
  const hash = route.hash.replace('#', '')
  return hash || 'personal-center'
})

// 监听hash变化，如果没有hash，默认跳转到个人信息
watch(() => route.hash, (newHash) => {
  if (!newHash) {
    router.replace('/dashboard#personal-center')
  }
}, { immediate: true })
</script>
