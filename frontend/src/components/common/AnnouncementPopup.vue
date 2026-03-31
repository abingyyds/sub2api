<template>
  <BaseDialog :show="visible" :title="currentAnnouncement?.title || ''" width="normal" @close="dismissToday">
    <div v-if="currentAnnouncement" class="space-y-4">
      <!-- Badges -->
      <div class="flex flex-wrap items-center gap-2">
        <span v-if="currentAnnouncement.version" class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs font-medium text-gray-600 dark:bg-dark-700 dark:text-dark-300">
          {{ currentAnnouncement.version }}
        </span>
        <span :class="['inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium', categoryStyle(currentAnnouncement.category)]">
          {{ currentAnnouncement.category }}
        </span>
        <span class="text-xs text-gray-400 dark:text-dark-500">
          {{ formatDisplayDate(currentAnnouncement) }}
        </span>
      </div>

      <!-- Content -->
      <div v-if="currentAnnouncement.content" class="whitespace-pre-wrap text-sm leading-relaxed text-gray-600 dark:text-dark-300">
        {{ currentAnnouncement.content }}
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-3">
        <button @click="dismissToday" class="btn btn-secondary">
          {{ t('changelog.closeForToday') }}
        </button>
        <button @click="dismissThis" class="btn btn-primary">
          {{ t('changelog.dismissThis') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import type { Announcement } from '@/types'

const { t } = useI18n()

const props = defineProps<{
  announcements: Announcement[]
}>()

const unreadAnnouncements = ref<Announcement[]>([])
const currentIndex = ref(0)

const visible = computed(() => unreadAnnouncements.value.length > 0 && currentIndex.value < unreadAnnouncements.value.length)
const currentAnnouncement = computed(() => unreadAnnouncements.value[currentIndex.value] || null)

function getTodayKey(): string {
  const d = new Date()
  return `ann_dismissed_today_${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function isDismissedToday(): boolean {
  return localStorage.getItem(getTodayKey()) === 'true'
}

function isDismissedPermanently(id: number): boolean {
  return localStorage.getItem(`ann_dismissed_${id}`) === 'true'
}

function dismissToday() {
  localStorage.setItem(getTodayKey(), 'true')
  unreadAnnouncements.value = []
}

function dismissThis() {
  if (currentAnnouncement.value) {
    localStorage.setItem(`ann_dismissed_${currentAnnouncement.value.id}`, 'true')
  }
  currentIndex.value++
}

function categoryStyle(category: string): string {
  switch (category) {
    case '新功能': return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
    case '修复': return 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400'
    case '通知': return 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400'
    default: return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
  }
}

function formatDisplayDate(ann: Announcement): string {
  const d = new Date(ann.published_at || ann.created_at)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

watch(() => props.announcements, (newVal) => {
  if (!newVal.length || isDismissedToday()) return
  unreadAnnouncements.value = newVal.filter(ann => !isDismissedPermanently(ann.id))
  currentIndex.value = 0
}, { immediate: true })
</script>
