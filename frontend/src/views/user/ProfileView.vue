<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-4xl space-y-6">
        <StaggerContainer :stagger-delay="100">
          <div class="grid grid-cols-1 gap-6 sm:grid-cols-3">
            <GlowCard glow-color="rgb(34, 197, 94)">
              <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft p-5 flex items-start gap-4">
                <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-emerald-100 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400">
                  <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path d="M21 12a2.25 2.25 0 00-2.25-2.25H15a3 3 0 11-6 0H5.25A2.25 2.25 0 003 12" /></svg>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('profile.accountBalance') }}</p>
                  <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ formatCurrency(user?.balance || 0) }}</p>
                </div>
              </div>
            </GlowCard>
            <GlowCard glow-color="rgb(245, 158, 11)">
              <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft p-5 flex items-start gap-4">
                <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-amber-100 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400">
                  <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path d="m3.75 13.5 10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" /></svg>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('profile.concurrencyLimit') }}</p>
                  <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ user?.concurrency || 0 }}</p>
                </div>
              </div>
            </GlowCard>
            <GlowCard glow-color="rgb(59, 130, 246)">
              <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft p-5 flex items-start gap-4">
                <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-primary-100 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400">
                  <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path d="M6.75 3v2.25M17.25 3v2.25" /></svg>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('profile.memberSince') }}</p>
                  <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ formatDate(user?.created_at || '', { year: 'numeric', month: 'long' }) }}</p>
                </div>
              </div>
            </GlowCard>
          </div>
        </StaggerContainer>
        <SlideIn direction="up" :delay="300">
          <ProfileInfoCard :user="user" />
        </SlideIn>
        <SlideIn direction="up" :delay="400">
          <div v-if="contactDisplayText" class="rounded-2xl border-2 border-primary-200 bg-primary-50 dark:border-primary-800/30 dark:bg-primary-900/20 shadow-soft p-6 cursor-pointer" @click="appStore.showContactModal = true">
            <div class="flex items-center gap-4">
              <div class="p-3 bg-primary-100 rounded-xl text-primary-600"><Icon name="chat" size="lg" /></div>
              <div><h3 class="font-semibold text-primary-800 dark:text-primary-200">{{ t('common.contactSupport') }}</h3><p class="text-sm font-medium">{{ contactDisplayText }}</p></div>
            </div>
          </div>
        </SlideIn>
        <SlideIn direction="up" :delay="500">
          <ProfileEditForm :initial-username="user?.username || ''" />
        </SlideIn>
        <SlideIn direction="up" :delay="600">
          <ProfilePasswordForm />
        </SlideIn>
        <SlideIn direction="up" :delay="700">
          <ProfileTotpCard />
        </SlideIn>
      </div>
    </FadeIn>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores'
import { formatDate } from '@/utils/format'
import AppLayout from '@/components/layout/AppLayout.vue'
import ProfileInfoCard from '@/components/user/profile/ProfileInfoCard.vue'
import ProfileEditForm from '@/components/user/profile/ProfileEditForm.vue'
import ProfilePasswordForm from '@/components/user/profile/ProfilePasswordForm.vue'
import ProfileTotpCard from '@/components/user/profile/ProfileTotpCard.vue'
import { Icon } from '@/components/icons'
import { FadeIn, SlideIn, StaggerContainer, GlowCard } from '@/components/animations'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()
const user = computed(() => authStore.user)

const contactDisplayText = computed(() => {
  const raw = appStore.contactInfo
  if (!raw) return ''
  try {
    const parsed = JSON.parse(raw)
    return parsed.wechat_id || parsed.email || ''
  } catch {
    return raw
  }
})

const formatCurrency = (v: number) => `$${v.toFixed(2)}`
</script>
