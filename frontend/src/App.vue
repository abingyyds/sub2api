<script setup lang="ts">
import { RouterView, useRouter, useRoute } from 'vue-router'
import { onMounted, ref, watch } from 'vue'
import Toast from '@/components/common/Toast.vue'
import NavigationProgress from '@/components/common/NavigationProgress.vue'
import DiscoverySourceModal from '@/components/auth/DiscoverySourceModal.vue'
import { useAppStore, useAuthStore, useSubscriptionStore } from '@/stores'
import { getSetupStatus } from '@/api/setup'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()
const subscriptionStore = useSubscriptionStore()

const showDiscoverySource = ref(false)

/**
 * Update favicon dynamically
 * @param logoUrl - URL of the logo to use as favicon
 */
function updateFavicon(logoUrl: string) {
  // Find existing favicon link or create new one
  let link = document.querySelector<HTMLLinkElement>('link[rel="icon"]')
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  }
  link.type = logoUrl.endsWith('.svg') ? 'image/svg+xml' : 'image/x-icon'
  link.href = logoUrl
}

// Watch for site settings changes and update favicon/title
watch(
  () => appStore.siteFavicon || appStore.siteLogo,
  (newFavicon) => {
    if (newFavicon) {
      updateFavicon(newFavicon)
    }
  },
  { immediate: true }
)

watch(
  () => appStore.siteName,
  (newName) => {
    if (newName) {
      document.title = `${newName} - AI API Gateway`
    }
  },
  { immediate: true }
)

// Watch for authentication state and manage subscription data
watch(
  () => authStore.isAuthenticated,
  (isAuthenticated) => {
    if (isAuthenticated) {
      // User logged in: start polling (first fetch will happen on-demand)
      subscriptionStore.startPolling()

      // Show discovery source survey if user hasn't filled it in
      if (authStore.user && authStore.user.discovery_source == null) {
        showDiscoverySource.value = true
      }
    } else {
      // User logged out: clear data and stop polling
      subscriptionStore.clear()
      showDiscoverySource.value = false
    }
  },
  { immediate: true }
)

onMounted(() => {
  // 异步检查 setup 状态，不阻塞渲染
  getSetupStatus()
    .then((status) => {
      if (status.needs_setup && route.path !== '/setup') {
        router.replace('/setup')
      }
    })
    .catch(() => {
      // If setup endpoint fails, assume normal mode and continue
    })

  // Public settings are already loaded from injected config or will be loaded on-demand
})
</script>

<template>
  <NavigationProgress />
  <RouterView v-slot="{ Component }">
    <Transition name="page" mode="out-in">
      <component :is="Component" />
    </Transition>
  </RouterView>
  <Toast />
  <DiscoverySourceModal :show="showDiscoverySource" @close="showDiscoverySource = false" />
</template>
