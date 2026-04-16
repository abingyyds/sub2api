<template>
  <div class="min-h-screen bg-gray-50 dark:bg-dark-950">
    <!-- Background Decoration -->
    <div class="pointer-events-none fixed inset-0 bg-mesh-gradient"></div>

    <!-- Sidebar -->
    <AppSidebar />

    <!-- Main Content Area -->
    <div
      class="relative min-h-screen transition-all duration-300"
      :class="[sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-64']"
    >
      <!-- Header -->
      <AppHeader />

      <!-- Main Content -->
      <main class="p-4 md:p-6 lg:p-8">
        <slot />
      </main>
    </div>

    <!-- Announcement Popup -->
    <AnnouncementPopup
      v-if="shouldRenderAnnouncementPopup"
      :announcements="announcements"
    />

    <!-- Contact Modal -->
    <ContactModal
      v-if="appStore.showContactModal"
      :show="appStore.showContactModal"
      @close="appStore.showContactModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import '@/styles/onboarding.css'
import { computed, defineAsyncComponent, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import { useOnboardingTour } from '@/composables/useOnboardingTour'
import { useOnboardingStore } from '@/stores/onboarding'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const AnnouncementPopup = defineAsyncComponent(
  () => import('@/components/common/AnnouncementPopup.vue')
)
const ContactModal = defineAsyncComponent(() => import('@/components/common/ContactModal.vue'))

const appStore = useAppStore()
const authStore = useAuthStore()
const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const isAdmin = computed(() => authStore.isAdmin)
const announcements = computed(() => Array.isArray(appStore.announcements) ? appStore.announcements : [])
const shouldRenderAnnouncementPopup = computed(() => announcements.value.length > 0)

const { replayTour } = useOnboardingTour({
  storageKey: isAdmin.value ? 'admin_guide' : 'user_guide',
  autoStart: true
})

const onboardingStore = useOnboardingStore()

onMounted(() => {
  onboardingStore.setReplayCallback(replayTour)
  const loadAnnouncements = () => {
    void appStore.fetchAnnouncements()
  }

  const scheduleAnnouncementsLoad = () => {
    if (typeof window.requestIdleCallback === 'function') {
      window.requestIdleCallback(loadAnnouncements, { timeout: 3000 })
    } else {
      window.setTimeout(loadAnnouncements, 1200)
    }
  }

  if (document.readyState === 'complete') {
    scheduleAnnouncementsLoad()
  } else {
    window.addEventListener('load', scheduleAnnouncementsLoad, { once: true })
  }
})

defineExpose({ replayTour })
</script>
