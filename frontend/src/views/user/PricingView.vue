<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-8">
      <!-- Title -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('pricing.title') }}</h1>
        <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('pricing.subtitle') }}</p>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
      </div>

      <!-- Plans Grid -->
      <div v-else-if="plans.length > 0" class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="(plan, index) in plans"
          :key="plan.key"
          class="card relative flex flex-col overflow-hidden transition-all hover:shadow-lg hover:-translate-y-1 cursor-pointer"
          :class="{ 'ring-2 ring-primary-500': index === 0 }"
          @click="handleBuy(plan)"
        >
          <div v-if="index === 0" class="absolute right-4 top-4">
            <span class="badge badge-primary text-xs">{{ t('pricing.recommended') }}</span>
          </div>
          <div class="flex flex-1 flex-col p-6">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">{{ plan.name }}</h3>
            <div class="mt-3 text-3xl font-bold text-primary-600">
              ¥{{ (plan.amount_fen / 100).toFixed(plan.amount_fen % 100 === 0 ? 0 : 2) }}
            </div>
            <ul class="mt-5 flex-1 space-y-2.5">
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <svg class="h-4 w-4 flex-shrink-0 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                </svg>
                {{ t('pricing.validityDays', { days: plan.validity_days }) }}
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <svg class="h-4 w-4 flex-shrink-0 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                </svg>
                {{ t('pricing.wechatPay') }}
              </li>
            </ul>
            <div
              class="mt-4 flex items-center justify-center rounded-lg bg-primary-600 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
              :class="{ 'opacity-50 pointer-events-none': creatingOrder }"
            >
              {{ t('pricing.buyNow') }}
            </div>
          </div>
        </div>
      </div>

      <!-- No Plans -->
      <div v-else class="text-center py-12 text-gray-500 dark:text-dark-400">
        {{ t('pricing.noPlans') }}
      </div>

      <p v-if="plans.length > 0" class="text-center text-sm text-gray-500 dark:text-dark-400">
        {{ t('pricing.stackable') }}
      </p>
    </div>

    <!-- Payment QR Code Modal -->
    <Teleport to="body">
      <div
        v-if="showPaymentModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="cancelPayment"
      >
        <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-2xl dark:bg-dark-700 mx-4">
          <!-- Header -->
          <div class="mb-6 text-center">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('pricing.payment.scanToPay') }}</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
              {{ t('pricing.payment.amount') }}: <span class="font-bold text-primary-600">¥{{ currentOrderAmount }}</span>
            </p>
          </div>

          <!-- QR Code -->
          <div class="flex justify-center mb-6">
            <div v-if="qrLoading" class="flex h-48 w-48 items-center justify-center">
              <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
            </div>
            <canvas v-else ref="qrCanvas" class="rounded-lg"></canvas>
          </div>

          <!-- Status -->
          <div class="mb-4 text-center">
            <div v-if="paymentStatus === 'pending'" class="flex items-center justify-center gap-2 text-sm text-gray-500 dark:text-dark-400">
              <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
              {{ t('pricing.payment.waitingPayment') }}
            </div>
            <div v-else-if="paymentStatus === 'paid'" class="flex items-center justify-center gap-2 text-sm text-green-600">
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ t('pricing.payment.paymentSuccess') }}
            </div>
            <div v-else-if="paymentStatus === 'closed'" class="text-sm text-red-500">
              {{ t('pricing.payment.orderExpired') }}
            </div>
          </div>

          <!-- Countdown -->
          <div v-if="paymentStatus === 'pending' && countdown > 0" class="mb-4 text-center text-xs text-gray-400 dark:text-dark-500">
            {{ t('pricing.payment.expiresIn', { minutes: Math.floor(countdown / 60), seconds: countdown % 60 }) }}
          </div>

          <!-- Actions -->
          <div class="flex gap-3">
            <button
              v-if="paymentStatus !== 'paid'"
              class="btn flex-1 border border-gray-300 dark:border-dark-500"
              @click="cancelPayment"
            >
              {{ t('pricing.payment.cancel') }}
            </button>
            <button
              v-if="paymentStatus === 'paid'"
              class="btn btn-primary flex-1"
              @click="goToSubscriptions"
            >
              {{ t('pricing.payment.viewSubscription') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import QRCode from 'qrcode'
import AppLayout from '@/components/layout/AppLayout.vue'
import { paymentAPI } from '@/api/payment'
import type { PaymentPlan } from '@/api/payment'

const { t } = useI18n()
const router = useRouter()

const loading = ref(true)
const plans = ref<PaymentPlan[]>([])
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

onMounted(async () => {
  try {
    const allPlans = await paymentAPI.getPlans()
    // Only show subscription plans on this page
    plans.value = allPlans.filter(p => (p.type || 'subscription') === 'subscription')
  } catch {
    // silently fail, show empty state
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  clearTimers()
})

function clearTimers() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
}

async function handleBuy(plan: PaymentPlan) {
  if (creatingOrder.value) return

  creatingOrder.value = true
  try {
    const order = await paymentAPI.createOrder(plan.key)
    currentOrderNo.value = order.order_no
    currentOrderAmount.value = (order.amount_fen / 100).toFixed(order.amount_fen % 100 === 0 ? 0 : 2)
    paymentStatus.value = 'pending'
    showPaymentModal.value = true

    // Calculate countdown
    const expiresAt = new Date(order.expired_at).getTime()
    countdown.value = Math.max(0, Math.floor((expiresAt - Date.now()) / 1000))

    // Render QR code
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

    // Start polling for payment status
    startPolling()

    // Start countdown
    countdownTimer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        paymentStatus.value = 'closed'
        clearTimers()
      }
    }, 1000)
  } catch (err: any) {
    const msg = err?.message || err?.response?.data?.message || t('pricing.payment.createFailed')
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

function goToSubscriptions() {
  showPaymentModal.value = false
  clearTimers()
  router.push('/subscriptions')
}
</script>
