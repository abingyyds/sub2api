<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-8">
      <!-- Title -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('modelPlaza.title') }}</h1>
        <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('modelPlaza.subtitle') }}</p>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent"></div>
      </div>

      <!-- Empty State -->
      <div v-else-if="!groupModels.length" class="py-12 text-center text-gray-500 dark:text-dark-400">
        {{ t('modelPlaza.empty') }}
      </div>

      <!-- Group Cards -->
      <div v-else class="grid gap-6 md:grid-cols-2">
        <div
          v-for="item in groupModels"
          :key="item.group.id"
          class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-800"
        >
          <!-- Group Header -->
          <div class="mb-4 flex items-center justify-between">
            <div>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ item.group.name }}</h2>
              <p v-if="item.group.description" class="mt-1 text-sm text-gray-500 dark:text-dark-400">
                {{ item.group.description }}
              </p>
            </div>
            <span class="rounded-full px-3 py-1 text-xs font-medium"
              :class="platformBadgeClass(item.group.platform)">
              {{ item.group.platform }}
            </span>
          </div>

          <!-- Group Meta -->
          <div class="mb-4 flex flex-wrap gap-3 text-xs text-gray-500 dark:text-dark-400">
            <span v-if="item.group.rate_multiplier !== 1">
              {{ t('modelPlaza.rate') }}: {{ item.group.rate_multiplier }}x
            </span>
            <span v-if="item.group.subscription_type === 'subscription'" class="text-primary-600 dark:text-primary-400">
              {{ t('modelPlaza.subscriptionType') }}
            </span>
            <span v-else class="text-green-600 dark:text-green-400">
              {{ t('modelPlaza.standardType') }}
            </span>
            <span v-if="item.group.claude_code_only" class="text-amber-600 dark:text-amber-400">
              Claude Code Only
            </span>
          </div>

          <!-- Models List -->
          <div v-if="item.models && item.models.length">
            <div class="mb-2 text-sm font-medium text-gray-700 dark:text-dark-300">
              {{ t('modelPlaza.availableModels') }} ({{ item.models.length }})
            </div>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="model in item.models"
                :key="model"
                class="inline-block rounded-md bg-gray-100 px-2.5 py-1 text-xs font-mono text-gray-700 dark:bg-dark-700 dark:text-dark-300"
              >
                {{ model }}
              </span>
            </div>
          </div>
          <div v-else class="text-sm text-gray-400 dark:text-dark-500 italic">
            {{ t('modelPlaza.allModels') }}
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import { getModelPlaza, type GroupModels } from '@/api/model-plaza'

const { t } = useI18n()

const loading = ref(true)
const groupModels = ref<GroupModels[]>([])

function platformBadgeClass(platform: string) {
  const map: Record<string, string> = {
    anthropic: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
    openai: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    gemini: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
    antigravity: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400',
    multi: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300',
  }
  return map[platform] || map.multi
}

onMounted(async () => {
  try {
    groupModels.value = await getModelPlaza()
  } catch (e) {
    console.error('Failed to load model plaza:', e)
  } finally {
    loading.value = false
  }
})
</script>
