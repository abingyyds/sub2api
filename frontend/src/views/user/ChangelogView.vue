<template>
  <AppLayout>
    <div class="mx-auto max-w-4xl">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('changelog.title') }}</h1>
        <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('changelog.subtitle') }}</p>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
      </div>

      <!-- Empty State -->
      <div v-else-if="!announcements.length" class="py-16 text-center">
        <svg class="mx-auto h-12 w-12 text-gray-300 dark:text-dark-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M10.34 15.84c-.688-.06-1.386-.09-2.09-.09H7.5a4.5 4.5 0 110-9h.75c.704 0 1.402-.03 2.09-.09m0 9.18c.253.962.584 1.892.985 2.783.247.55.06 1.21-.463 1.511l-.657.38c-.551.318-1.26.117-1.527-.461a20.845 20.845 0 01-1.44-4.282m3.102.069a18.03 18.03 0 01-.59-4.59c0-1.586.205-3.124.59-4.59m0 9.18a23.848 23.848 0 018.835 2.535M10.34 6.66a23.847 23.847 0 008.835-2.535m0 0A23.74 23.74 0 0018.795 3m.38 1.125a23.91 23.91 0 011.014 5.395m-1.014 8.855c-.118.38-.245.754-.38 1.125m.38-1.125a23.91 23.91 0 001.014-5.395m0-3.46c.495.413.811 1.035.811 1.73 0 .695-.316 1.317-.811 1.73m0-3.46a24.347 24.347 0 010 3.46" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900 dark:text-white">{{ t('changelog.noData') }}</h3>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ t('changelog.noDataDesc') }}</p>
      </div>

      <!-- Timeline -->
      <div v-else class="relative">
        <!-- Vertical line -->
        <div class="absolute left-[139px] top-0 hidden h-full w-0.5 bg-gray-200 dark:bg-dark-700 md:block"></div>

        <!-- Timeline entries -->
        <div v-for="(ann, index) in announcements" :key="ann.id" class="relative mb-8 last:mb-0">
          <div class="flex flex-col gap-4 md:flex-row">
            <!-- Date column (left side) -->
            <div class="flex w-full flex-shrink-0 items-start md:w-[120px] md:text-right">
              <div class="md:w-full">
                <div class="text-lg font-bold text-gray-900 dark:text-white">
                  {{ formatMonth(ann) }}{{ formatDay(ann) }}日
                </div>
                <div class="text-sm text-gray-500 dark:text-dark-400">
                  {{ formatYear(ann) }}
                </div>
              </div>
            </div>

            <!-- Timeline dot -->
            <div class="relative hidden flex-shrink-0 md:flex md:w-[40px] md:justify-center">
              <div class="relative z-10 mt-2 h-3 w-3 rounded-full border-2 border-primary-500 bg-white dark:bg-dark-900"
                :class="index === 0 ? 'ring-4 ring-primary-100 dark:ring-primary-900/30' : ''"></div>
            </div>

            <!-- Content card (right side) -->
            <div class="min-w-0 flex-1">
              <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-shadow hover:shadow-md dark:border-dark-700 dark:bg-dark-800">
                <!-- Title row with badges -->
                <div class="mb-3 flex flex-wrap items-center gap-2">
                  <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ ann.title }}</h3>
                  <span v-if="ann.version" class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs font-medium text-gray-600 dark:bg-dark-700 dark:text-dark-300">
                    {{ ann.version }}
                  </span>
                  <span :class="['inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium', categoryStyle(ann.category)]">
                    {{ ann.category }}
                  </span>
                </div>

                <!-- Content -->
                <div v-if="ann.content" class="whitespace-pre-wrap text-sm leading-relaxed text-gray-600 dark:text-dark-300">{{ ann.content }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { authAPI } from '@/api/auth'
import AppLayout from '@/components/layout/AppLayout.vue'
import type { Announcement } from '@/types'

const { t } = useI18n()

const announcements = ref<Announcement[]>([])
const loading = ref(true)

function getDate(ann: Announcement): Date {
  return new Date(ann.published_at || ann.created_at)
}

function formatMonth(ann: Announcement): string {
  return String(getDate(ann).getMonth() + 1)
}

function formatDay(ann: Announcement): string {
  return String(getDate(ann).getDate())
}

function formatYear(ann: Announcement): string {
  return String(getDate(ann).getFullYear())
}

function categoryStyle(category: string): string {
  switch (category) {
    case '新功能': return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
    case '修复': return 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400'
    case '通知': return 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400'
    default: return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
  }
}

onMounted(async () => {
  try {
    announcements.value = await authAPI.getActiveAnnouncements()
  } catch {
    // silently fail
  } finally {
    loading.value = false
  }
})
</script>
