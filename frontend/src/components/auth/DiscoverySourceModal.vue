<template>
  <BaseDialog
    :show="show"
    :title="t('auth.discoverySource.title')"
    width="narrow"
    :close-on-escape="false"
    @close="handleSkip"
  >
    <p class="mb-5 text-sm text-gray-500 dark:text-gray-400">
      {{ t('auth.discoverySource.subtitle') }}
    </p>

    <div class="grid grid-cols-2 gap-3">
      <button
        v-for="option in sourceOptions"
        :key="option.value"
        type="button"
        :class="[
          'flex items-center gap-2.5 rounded-lg border px-3 py-2.5 text-sm transition-colors',
          selectedSource === option.value
            ? 'border-primary-500 bg-primary-50 text-primary-700 dark:border-primary-400 dark:bg-primary-900/20 dark:text-primary-300'
            : 'border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50 dark:border-dark-600 dark:text-gray-300 dark:hover:border-dark-500 dark:hover:bg-dark-700'
        ]"
        @click="selectSource(option.value)"
      >
        <svg class="h-4 w-4 shrink-0" :style="{ color: option.color }" viewBox="0 0 24 24" fill="currentColor">
          <path v-for="(d, i) in option.paths" :key="i" :d="d" />
        </svg>
        <span>{{ option.label }}</span>
      </button>
    </div>

    <!-- Other input -->
    <div v-if="selectedSource === 'other'" class="mt-3">
      <input
        v-model="otherText"
        type="text"
        maxlength="50"
        class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white"
        :placeholder="t('auth.discoverySource.otherPlaceholder')"
      />
    </div>

    <template #footer>
      <div class="flex justify-between">
        <button
          type="button"
          class="text-sm text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
          :disabled="submitting"
          @click="handleSkip"
        >
          {{ t('auth.discoverySource.skip') }}
        </button>
        <button
          type="button"
          class="btn btn-primary"
          :disabled="!canSubmit || submitting"
          @click="handleSubmit"
        >
          {{ submitting ? t('auth.discoverySource.submitting') : t('auth.discoverySource.submit') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { authAPI } from '@/api'
import { useAuthStore } from '@/stores/auth'

defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const { t } = useI18n()
const authStore = useAuthStore()

const selectedSource = ref<string | null>(null)
const otherText = ref('')
const submitting = ref(false)

// Brand SVG icon paths
const sourceOptions = computed(() => [
  {
    value: 'douyin',
    color: '#000000',
    paths: ['M12.53.02C13.84 0 15.14.01 16.44 0c.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z'],
    label: t('auth.discoverySource.douyin')
  },
  {
    value: 'xiaohongshu',
    color: '#FF2442',
    paths: ['M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm4.64 7.27h-1.62c-.13 0-.24.07-.3.18l-1.57 2.9h-.06l-1.58-2.9c-.06-.11-.17-.18-.3-.18H9.6c-.22 0-.35.24-.24.43l2.58 4.33v2.74c0 .17.14.31.31.31h1.5c.17 0 .31-.14.31-.31v-2.74l2.58-4.33c.12-.19-.02-.43-.24-.43h.04z'],
    label: t('auth.discoverySource.xiaohongshu')
  },
  {
    value: 'bilibili',
    color: '#00A1D6',
    paths: ['M17.813 4.653h.854c1.51.054 2.769.578 3.773 1.574 1.004.995 1.524 2.249 1.56 3.76v7.36c-.036 1.51-.556 2.769-1.56 3.773s-2.262 1.524-3.773 1.56H5.333c-1.51-.036-2.769-.556-3.773-1.56S.036 18.858 0 17.347v-7.36c.036-1.511.556-2.765 1.56-3.76 1.004-.996 2.262-1.52 3.773-1.574h.774l-1.174-1.12a1.234 1.234 0 0 1-.373-.906c0-.356.124-.658.373-.907l.027-.027c.267-.249.573-.373.92-.373.347 0 .653.124.92.373L9.653 4.44c.071.071.134.142.187.213h4.267a.836.836 0 0 1 .16-.213l2.853-2.747c.267-.249.573-.373.92-.373.347 0 .662.151.929.4.267.249.391.551.391.907 0 .355-.124.657-.373.906zM5.333 7.24c-.746.018-1.373.276-1.88.773-.506.498-.769 1.13-.786 1.894v7.52c.017.764.28 1.395.786 1.893.507.498 1.134.756 1.88.773h13.334c.746-.017 1.373-.275 1.88-.773.506-.498.769-1.129.786-1.893v-7.52c-.017-.765-.28-1.396-.786-1.894-.507-.497-1.134-.755-1.88-.773zM8 11.107c.373 0 .684.124.933.373.25.249.383.569.4.96v1.173c-.017.391-.15.711-.4.96-.249.25-.56.374-.933.374s-.684-.125-.933-.374c-.25-.249-.383-.569-.4-.96V12.44c.017-.391.15-.711.4-.96.249-.249.56-.373.933-.373zm8 0c.373 0 .684.124.933.373.25.249.383.569.4.96v1.173c-.017.391-.15.711-.4.96-.249.25-.56.374-.933.374s-.684-.125-.933-.374c-.25-.249-.383-.569-.4-.96V12.44c.017-.391.15-.711.4-.96.249-.249.56-.373.933-.373z'],
    label: t('auth.discoverySource.bilibili')
  },
  {
    value: 'wechat',
    color: '#07C160',
    paths: ['M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm3.655 4.467c-2.322.037-4.592.826-6.27 2.22-1.93 1.6-2.88 3.98-2.036 6.315.48 1.27 1.39 2.39 2.59 3.17a.504.504 0 0 1 .18.56l-.28 1.06c-.009.04-.018.08-.018.12 0 .1.076.17.167.17a.18.18 0 0 0 .096-.03l1.37-.8a.628.628 0 0 1 .516-.07c.79.18 1.59.27 2.38.27 3.76 0 6.89-2.57 7.29-5.88.42-3.42-2.36-6.53-5.98-7.1zm-2.272 3.39c.496 0 .898.407.898.91a.904.904 0 0 1-.898.909.904.904 0 0 1-.898-.91c0-.502.402-.909.898-.909zm4.49 0c.497 0 .899.407.899.91a.904.904 0 0 1-.899.909.904.904 0 0 1-.898-.91c0-.502.402-.909.898-.909z'],
    label: t('auth.discoverySource.wechat')
  },
  {
    value: 'twitter',
    color: '#000000',
    paths: ['M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z'],
    label: t('auth.discoverySource.twitter')
  },
  {
    value: 'telegram',
    color: '#26A5E4',
    paths: ['M11.944 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.96 6.504-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z'],
    label: t('auth.discoverySource.telegram')
  },
  {
    value: 'github',
    color: '#181717',
    paths: ['M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12'],
    label: t('auth.discoverySource.github')
  },
  {
    value: 'search',
    color: '#4285F4',
    paths: ['M15.5 14h-.79l-.28-.27A6.471 6.471 0 0 0 16 9.5 6.5 6.5 0 1 0 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z'],
    label: t('auth.discoverySource.search')
  },
  {
    value: 'friend',
    color: '#6366F1',
    paths: ['M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z'],
    label: t('auth.discoverySource.friend')
  },
  {
    value: 'other',
    color: '#9CA3AF',
    paths: ['M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 17h-2v-2h2v2zm2.07-7.75l-.9.92C13.45 12.9 13 13.5 13 15h-2v-.5c0-1.1.45-2.1 1.17-2.83l1.24-1.26c.37-.36.59-.86.59-1.41 0-1.1-.9-2-2-2s-2 .9-2 2H8c0-2.21 1.79-4 4-4s4 1.79 4 4c0 .88-.36 1.68-.93 2.25z'],
    label: t('auth.discoverySource.other')
  }
])

const canSubmit = computed(() => {
  if (!selectedSource.value) return false
  if (selectedSource.value === 'other') return otherText.value.trim().length > 0
  return true
})

function selectSource(value: string) {
  selectedSource.value = value
  if (value !== 'other') {
    otherText.value = ''
  }
}

async function handleSubmit() {
  if (!canSubmit.value) return
  submitting.value = true
  try {
    const source = selectedSource.value === 'other'
      ? otherText.value.trim()
      : selectedSource.value!
    await authAPI.updateDiscoverySource(source)
    await authStore.refreshUser()
    emit('close')
  } catch (error) {
    console.error('Failed to update discovery source:', error)
    emit('close')
  } finally {
    submitting.value = false
  }
}

async function handleSkip() {
  submitting.value = true
  try {
    await authAPI.updateDiscoverySource('skip')
    await authStore.refreshUser()
  } catch {
    // Ignore errors on skip
  } finally {
    submitting.value = false
    emit('close')
  }
}
</script>
