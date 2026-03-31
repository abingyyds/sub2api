<template>
  <Teleport to="body">
    <transition name="fade">
      <div
        v-if="show"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
        @click.self="emit('close')"
      >
        <div class="mx-4 w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl dark:bg-dark-900 space-y-5">
          <!-- Header -->
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('contactUs.title') }}</h3>
            <button
              class="rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-dark-700 dark:hover:text-dark-200"
              @click="emit('close')"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- QR Code -->
          <div v-if="contactData.qrcode" class="flex justify-center">
            <div class="rounded-xl border border-gray-200 bg-white p-3 dark:border-dark-600 dark:bg-dark-800">
              <img :src="contactData.qrcode" alt="QR Code" class="h-48 w-48 object-contain" />
            </div>
          </div>
          <p v-if="contactData.qrcode" class="text-center text-xs text-gray-500 dark:text-dark-400">
            {{ t('contactUs.scanQrcode') }}
          </p>

          <!-- WeChat ID -->
          <div v-if="contactData.wechat_id" class="flex items-center justify-between rounded-xl bg-gray-50 px-4 py-3 dark:bg-dark-800">
            <div>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('contactUs.wechatId') }}</p>
              <p class="mt-0.5 font-medium text-gray-900 dark:text-white">{{ contactData.wechat_id }}</p>
            </div>
            <button
              class="rounded-lg px-3 py-1.5 text-xs font-medium transition"
              :class="copiedField === 'wechat' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' : 'bg-gray-200 text-gray-600 hover:bg-gray-300 dark:bg-dark-600 dark:text-dark-300 dark:hover:bg-dark-500'"
              @click="copyText(contactData.wechat_id, 'wechat')"
            >
              {{ copiedField === 'wechat' ? t('contactUs.copySuccess') : t('common.copy') }}
            </button>
          </div>

          <!-- Email -->
          <div v-if="contactData.email" class="flex items-center justify-between rounded-xl bg-gray-50 px-4 py-3 dark:bg-dark-800">
            <div>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('contactUs.email') }}</p>
              <p class="mt-0.5 font-medium text-gray-900 dark:text-white">{{ contactData.email }}</p>
            </div>
            <button
              class="rounded-lg px-3 py-1.5 text-xs font-medium transition"
              :class="copiedField === 'email' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' : 'bg-gray-200 text-gray-600 hover:bg-gray-300 dark:bg-dark-600 dark:text-dark-300 dark:hover:bg-dark-500'"
              @click="copyText(contactData.email, 'email')"
            >
              {{ copiedField === 'email' ? t('contactUs.copySuccess') : t('common.copy') }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'

const { t } = useI18n()
const appStore = useAppStore()

defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const copiedField = ref<string | null>(null)

interface ContactData {
  wechat_id: string
  email: string
  qrcode: string
}

const contactData = computed<ContactData>(() => {
  const raw = appStore.contactInfo
  if (!raw) return { wechat_id: '', email: '', qrcode: '' }

  // Try to parse as JSON first (new structured format)
  try {
    const parsed = JSON.parse(raw)
    return {
      wechat_id: parsed.wechat_id || '',
      email: parsed.email || '',
      qrcode: parsed.qrcode || '',
    }
  } catch {
    // Fallback: treat as plain text (old format, just show as wechat_id)
    return {
      wechat_id: raw,
      email: '',
      qrcode: '',
    }
  }
})

async function copyText(text: string, field: string) {
  try {
    await navigator.clipboard.writeText(text)
    copiedField.value = field
    setTimeout(() => {
      copiedField.value = null
    }, 2000)
  } catch {
    // fallback
  }
}
</script>
