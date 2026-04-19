<template>
  <div class="relative flex min-h-screen items-center justify-center overflow-hidden bg-gray-50 p-4 dark:bg-dark-950">
    <!-- Background -->
    <div
      class="absolute inset-0 bg-gradient-to-br from-white via-primary-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950"
    ></div>

    <!-- Content Container -->
    <div class="relative z-10 w-full max-w-md">
      <!-- Logo/Brand -->
      <div class="mb-6 text-center">
        <!-- Custom Logo or Default Logo -->
        <div
          class="mb-4 inline-flex h-14 w-14 items-center justify-center overflow-hidden rounded-2xl border border-white/70 bg-white/80 shadow-sm shadow-primary-500/10 dark:border-dark-700/70 dark:bg-dark-900/80"
        >
          <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
        </div>
        <h1 class="text-gradient mb-2 text-3xl font-bold">
          {{ siteName }}
        </h1>
        <p class="text-sm text-gray-500 dark:text-dark-400">
          {{ siteSubtitle }}
        </p>
      </div>

      <!-- Card Container -->
      <div class="rounded-2xl border border-white/70 bg-white/92 p-8 shadow-xl shadow-gray-200/40 backdrop-blur-sm dark:border-dark-800 dark:bg-dark-900/92 dark:shadow-black/20">
        <slot />
      </div>

      <!-- Footer Links -->
      <div class="mt-6 text-center text-sm">
        <slot name="footer" />
      </div>

      <!-- Copyright -->
      <div class="mt-8 text-center text-xs text-gray-400 dark:text-dark-500">
        &copy; {{ currentYear }} {{ siteName }}. All rights reserved.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'

const appStore = useAppStore()

const siteName = computed(() => appStore.siteName || 'cCoder.me')
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true }))
const siteSubtitle = computed(() => {
  const settings = appStore.cachedPublicSettings
  return settings?.site_subtitle || '新一代代码大师平台'
})

const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  // Ensure settings are loaded (uses cache if already fetched)
  if (!appStore.publicSettingsLoaded) {
    void appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.text-gradient {
  @apply bg-gradient-to-r from-primary-600 to-primary-500 bg-clip-text text-transparent;
}
</style>
