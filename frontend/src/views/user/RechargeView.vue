<template>
  <AppLayout>
    <div class="mx-auto max-w-2xl space-y-8">
      <!-- Title -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('recharge.title') }}</h1>
        <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('recharge.subtitle') }}</p>
      </div>

      <div class="card p-6 space-y-6">
        <!-- Quick Select Amounts -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-3">
            {{ t('recharge.quickSelect') }}
          </label>
          <div class="grid grid-cols-3 gap-3">
            <button
              v-for="preset in presets"
              :key="preset"
              class="rounded-lg border-2 px-4 py-3 text-center font-medium transition-all hover:border-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20"
              :class="selectedAmount === preset
                ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-400'
                : 'border-gray-200 text-gray-700 dark:border-dark-600 dark:text-dark-300'"
              @click="selectPreset(preset)"
            >
              ¥{{ preset }}
            </button>
          </div>
        </div>

        <!-- Custom Input -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-2">
            {{ t('recharge.customAmount') }}
          </label>
          <div class="relative">
            <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 dark:text-dark-400 font-medium">¥</span>
            <input
              v-model="customInput"
              type="number"
              min="1"
              step="1"
              :placeholder="t('recharge.inputPlaceholder')"
              class="w-full rounded-lg border border-gray-300 bg-white py-3 pl-8 pr-4 text-lg font-medium text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-800 dark:text-white dark:placeholder-dark-500"
              @input="onCustomInput"
            />
          </div>
        </div>

        <!-- Amount Summary -->
        <div v-if="finalAmount > 0" class="rounded-lg bg-gray-50 dark:bg-dark-800 p-4">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600 dark:text-dark-400">{{ t('recharge.rechargeAmount') }}</span>
            <span class="text-lg font-bold text-primary-600">¥{{ finalAmount.toFixed(2) }}</span>
          </div>
          <div class="mt-1 flex items-center justify-between">
            <span class="text-sm text-gray-600 dark:text-dark-400">{{ t('recharge.payAmount') }}</span>
            <span class="text-2xl font-bold text-gray-900 dark:text-white">¥{{ finalAmount.toFixed(2) }}</span>
          </div>
        </div>

        <!-- Payment Method -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-3">
            {{ t('recharge.paymentMethod') }}
          </label>
          <div class="flex items-center gap-3 rounded-lg border-2 border-primary-500 bg-primary-50 dark:bg-primary-900/20 p-3">
            <svg class="h-6 w-6 text-green-600" viewBox="0 0 24 24" fill="currentColor">
              <path d="M9.5 4C5.36 4 2 6.69 2 10c0 1.89 1.08 3.56 2.78 4.66L4 17l2.5-1.5C7.55 15.82 8.5 16 9.5 16c.34 0 .68-.02 1-.06A5.95 5.95 0 0110 14c0-3.31 2.69-6 6-6 .34 0 .68.03 1 .08C16.32 5.68 13.17 4 9.5 4zM7 9a1 1 0 110 2 1 1 0 010-2zm5 0a1 1 0 110 2 1 1 0 010-2zm4 3c-2.76 0-5 1.79-5 4s2.24 4 5 4c.71 0 1.39-.11 2-.31L20 21l-.5-1.8C20.45 18.22 21 17.16 21 16c0-2.21-2.24-4-5-4zm-1.5 2.5a.75.75 0 110 1.5.75.75 0 010-1.5zm3 0a.75.75 0 110 1.5.75.75 0 010-1.5z"/>
            </svg>
            <span class="font-medium text-gray-900 dark:text-white">{{ t('recharge.wechatPay') }}</span>
          </div>
        </div>

        <!-- Submit Button -->
        <button
          class="w-full rounded-lg bg-primary-600 py-3 text-base font-medium text-white transition-colors hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="finalAmount <= 0 || creatingOrder"
          @click="handleRecharge"
        >
          <span v-if="creatingOrder" class="flex items-center justify-center gap-2">
            <div class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
            {{ t('recharge.creating') }}
          </span>
          <span v-else>
            {{ finalAmount > 0 ? t('recharge.payNow', { amount: finalAmount.toFixed(2) }) : t('recharge.enterAmount') }}
          </span>
        </button>
      </div>
    </div>

    <!-- Payment QR Code Modal -->
    <Teleport to="body">
      <div
        v-if="showPaymentModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="cancelPayment"
      >
        <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-2xl dark:bg-dark-700 mx-4">
          <div class="mb-6 text-center">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('recharge.payment.scanToPay') }}</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
              {{ t('recharge.payment.amount') }}: <span class="font-bold text-primary-600">¥{{ currentOrderAmount }}</span>
            </p>
          </div>

          <div class="flex justify-center mb-6">
            <div v-if="qrLoading" class="flex h-48 w-48 items-center justify-center">
              <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
            </div>
            <canvas v-else ref="qrCanvas" class="rounded-lg"></canvas>
          </div>

          <div class="mb-4 text-center">
            <div v-if="paymentStatus === 'pending'" class="flex items-center justify-center gap-2 text-sm text-gray-500 dark:text-dark-400">
              <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
              {{ t('recharge.payment.waiting') }}
            </div>
            <div v-else-if="paymentStatus === 'paid'" class="flex items-center justify-center gap-2 text-sm text-green-600">
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ t('recharge.payment.success') }}
            </div>
            <div v-else-if="paymentStatus === 'closed'" class="text-sm text-red-500">
              {{ t('recharge.payment.expired') }}
            </div>
          </div>

          <div v-if="paymentStatus === 'pending' && countdown > 0" class="mb-4 text-center text-xs text-gray-400 dark:text-dark-500">
            {{ t('recharge.payment.expiresIn', { minutes: Math.floor(countdown / 60), seconds: countdown % 60 }) }}
          </div>

          <div class="flex gap-3">
            <button
              v-if="paymentStatus !== 'paid'"
              class="btn flex-1 border border-gray-300 dark:border-dark-500"
              @click="cancelPayment"
            >
              {{ t('recharge.payment.cancel') }}
            </button>
            <button
              v-if="paymentStatus === 'paid'"
              class="btn btn-primary flex-1"
              @click="goToDashboard"
            >
              {{ t('recharge.payment.viewDashboard') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import QRCode from 'qrcode'
import AppLayout from '@/components/layout/AppLayout.vue'
import { paymentAPI } from '@/api/payment'

const { t } = useI18n()
const router = useRouter()

const presets = [10, 50, 100, 200, 500, 1000]
const selectedAmount = ref<number | null>(null)
const customInput = ref<string>('')
const creatingOrder = ref(false)

const showPaymentModal = ref(false)
const qrLoading = ref(false)
const qrCanvas = ref<HTMLCanvasElement | null>(null)
const paymentStatus = ref<'pending' | 'paid' | 'closed'>('pending')
const currentOrderNo = ref('')
const currentOrderAmount = ref('')
const countdown = ref(0)

let pollTimer: ReturnType<typeof setInterval> | null = null
let countdownTimer: ReturnType<typeof setInterval> | null = null

const finalAmount = computed(() => {
  if (customInput.value !== '') {
    const val = parseFloat(customInput.value)
    return isNaN(val) || val <= 0 ? 0 : val
  }
  return selectedAmount.value ?? 0
})

function selectPreset(amount: number) {
  selectedAmount.value = amount
  customInput.value = ''
}

function onCustomInput() {
  selectedAmount.value = null
}

function clearTimers() {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
  if (countdownTimer) { clearInterval(countdownTimer); countdownTimer = null }
}

onUnmounted(() => clearTimers())

async function handleRecharge() {
  if (finalAmount.value <= 0 || creatingOrder.value) return

  creatingOrder.value = true
  try {
    const order = await paymentAPI.createRechargeOrder(finalAmount.value)
    currentOrderNo.value = order.order_no
    currentOrderAmount.value = (order.amount_fen / 100).toFixed(order.amount_fen % 100 === 0 ? 0 : 2)
    paymentStatus.value = 'pending'
    showPaymentModal.value = true

    const expiresAt = new Date(order.expired_at).getTime()
    countdown.value = Math.max(0, Math.floor((expiresAt - Date.now()) / 1000))

    if (order.code_url) {
      qrLoading.value = false
      await nextTick()
      if (qrCanvas.value) {
        await QRCode.toCanvas(qrCanvas.value, order.code_url, {
          width: 192,
          margin: 2,
          color: { dark: '#000000', light: '#ffffff' },
        })
      }
    }

    startPolling()

    countdownTimer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        paymentStatus.value = 'closed'
        clearTimers()
      }
    }, 1000)
  } catch (err: any) {
    const msg = err?.message || err?.response?.data?.message || t('recharge.payment.createFailed')
    alert(msg)
  } finally {
    creatingOrder.value = false
  }
}

function startPolling() {
  pollTimer = setInterval(async () => {
    try {
      const order = await paymentAPI.queryOrder(currentOrderNo.value)
      if (order.status === 'paid') {
        paymentStatus.value = 'paid'
        clearTimers()
      } else if (order.status === 'closed') {
        paymentStatus.value = 'closed'
        clearTimers()
      }
    } catch {
      // ignore poll errors
    }
  }, 3000)
}

function cancelPayment() {
  showPaymentModal.value = false
  clearTimers()
}

function goToDashboard() {
  showPaymentModal.value = false
  clearTimers()
  router.push('/dashboard')
}
</script>
