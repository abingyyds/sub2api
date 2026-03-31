<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-2xl">
        <SlideIn direction="up" :delay="100">
          <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft">
            <div class="flex flex-col items-center justify-center py-16 px-6 text-center">
              <div
                class="mb-6 inline-flex h-20 w-20 items-center justify-center rounded-2xl bg-primary-100 dark:bg-primary-900/30"
              >
                <Icon name="document" size="xl" class="text-primary-500 dark:text-primary-400" />
              </div>
              <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
                {{ t('nav.invoice') }}
              </h1>
              <p class="mt-3 text-sm text-gray-500 dark:text-dark-400">
                {{ t('dashboard.comingSoon') }}
              </p>
            </div>
          </div>
        </SlideIn>
      </div>
    </FadeIn>
  </AppLayout>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { FadeIn, SlideIn } from '@/components/animations'

const { t } = useI18n()
</script>
