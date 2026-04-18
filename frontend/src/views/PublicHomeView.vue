<template>
  <component :is="currentView" />
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted } from 'vue'
import { useAppStore } from '@/stores'

const MarketingHomeView = defineAsyncComponent(() => import('@/views/MarketingHomeView.vue'))
const HomeView = defineAsyncComponent(() => import('@/views/HomeView.vue'))

const appStore = useAppStore()

const currentView = computed(() => {
  if (appStore.cachedPublicSettings?.is_subsite) {
    return HomeView
  }
  return MarketingHomeView
})

onMounted(() => {
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>
