<template>
  <component :is="currentView" />
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import MarketingHomeView from '@/views/MarketingHomeView.vue'
import HomeView from '@/views/HomeView.vue'

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
