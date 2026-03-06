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
          'flex items-center gap-2 rounded-lg border px-3 py-2.5 text-sm transition-colors',
          selectedSource === option.value
            ? 'border-primary-500 bg-primary-50 text-primary-700 dark:border-primary-400 dark:bg-primary-900/20 dark:text-primary-300'
            : 'border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50 dark:border-dark-600 dark:text-gray-300 dark:hover:border-dark-500 dark:hover:bg-dark-700'
        ]"
        @click="selectSource(option.value)"
      >
        <span class="text-base">{{ option.icon }}</span>
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

const sourceOptions = computed(() => [
  { value: 'douyin', icon: '🎵', label: t('auth.discoverySource.douyin') },
  { value: 'xiaohongshu', icon: '📕', label: t('auth.discoverySource.xiaohongshu') },
  { value: 'bilibili', icon: '📺', label: t('auth.discoverySource.bilibili') },
  { value: 'wechat', icon: '💬', label: t('auth.discoverySource.wechat') },
  { value: 'twitter', icon: '🐦', label: t('auth.discoverySource.twitter') },
  { value: 'telegram', icon: '✈️', label: t('auth.discoverySource.telegram') },
  { value: 'github', icon: '🐙', label: t('auth.discoverySource.github') },
  { value: 'search', icon: '🔍', label: t('auth.discoverySource.search') },
  { value: 'friend', icon: '👥', label: t('auth.discoverySource.friend') },
  { value: 'other', icon: '💡', label: t('auth.discoverySource.other') }
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
