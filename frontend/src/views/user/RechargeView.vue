<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-2xl space-y-8">
        <!-- Title -->
        <SlideIn direction="up" :delay="100">
          <div class="text-center">
            <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('recharge.title') }}</h1>
            <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('recharge.subtitle') }}</p>
          </div>
        </SlideIn>

        <SlideIn direction="up" :delay="200">
          <GlowCard glow-color="rgb(59, 130, 246)">
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
                    class="w-full rounded-lg border border-gray-300 bg-white py-3 pl-8 pr-4 text-lg font-medium text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-800 dark:text-white dark:placeholder-dark-500 transition-all duration-300"
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

              <!-- Submit Button -->
              <MagneticButton>
                <button
                  class="w-full rounded-lg bg-primary-600 py-3 text-base font-medium text-white transition-all duration-300 hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
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
              </MagneticButton>
            </div>
          </GlowCard>
        </SlideIn>
      </div>
    </FadeIn>

    <!-- Payment Method Selection Modal -->
    <Teleport to="body">
      <div
        v-if="showPayMethodModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="showPayMethodModal = false"
      >
        <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-2xl dark:bg-dark-700 mx-4">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-4">选择支付方式</h3>

          <div class="space-y-3 mb-6">
            <!-- WeChat Pay -->
            <div
              class="flex items-center gap-3 rounded-lg border-2 p-4 cursor-pointer transition-all"
              :class="selectedPayMethod === 'wechat' ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20' : 'border-gray-200 dark:border-dark-600 hover:border-primary-300'"
              @click="selectedPayMethod = 'wechat'"
            >
              <input type="radio" :checked="selectedPayMethod === 'wechat'" class="h-5 w-5 text-primary-600" />
              <img src="data:image/png;base64,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" class="h-8 w-8" alt="微信支付" />
              <span class="font-medium text-gray-900 dark:text-white">微信支付</span>
            </div>

            <!-- Alipay -->
            <div
              class="flex items-center gap-3 rounded-lg border-2 p-4 cursor-pointer transition-all"
              :class="selectedPayMethod === 'alipay' ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20' : 'border-gray-200 dark:border-dark-600 hover:border-primary-300'"
              @click="selectedPayMethod = 'alipay'"
            >
              <input type="radio" :checked="selectedPayMethod === 'alipay'" class="h-5 w-5 text-primary-600" />
              <svg class="h-8 w-8 text-blue-500" viewBox="0 0 1024 1024" fill="currentColor">
                <path d="M1024 701.9v162.8c0 88.4-71.7 160.1-160.1 160.1H160.1C71.7 1024.8 0 953.1 0 864.7V159.9C0 71.5 71.7-0.2 160.1-0.2h703.8c88.4 0 160.1 71.7 160.1 160.1v390.7c-65.9-50.6-181.4-130.4-320-130.4-210.3 0-348.8 122.6-348.8 235.4 0 112.8 97.7 194.5 236.7 194.5 152.5 0 278.1-77.3 342.1-139.1zM736 511.5H288v-64h448v64z m0-128H288v-64h448v64z m0-128H288v-64h448v64z"/>
              </svg>
              <span class="font-medium text-gray-900 dark:text-white">支付宝</span>
            </div>
          </div>

          <div class="flex gap-3">
            <button
              class="btn flex-1 border border-gray-300 dark:border-dark-500"
              @click="showPayMethodModal = false"
            >
              取消
            </button>
            <button
              class="btn btn-primary flex-1"
              @click="confirmPayMethod"
            >
              确认支付
            </button>
          </div>
        </div>
      </div>
    </Teleport>

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
import { FadeIn, SlideIn, GlowCard, MagneticButton } from '@/components/animations'
import QRCode from 'qrcode'
import AppLayout from '@/components/layout/AppLayout.vue'
import { paymentAPI } from '@/api/payment'

const { t } = useI18n()
const router = useRouter()

const presets = [10, 50, 100, 200, 500, 1000]
const selectedAmount = ref<number | null>(null)
const customInput = ref<string>('')
const creatingOrder = ref(false)
const selectedPayMethod = ref<'wechat' | 'alipay'>('wechat')
const showPayMethodModal = ref(false)
const pendingPaymentAction = ref<(() => Promise<void>) | null>(null)

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

async function confirmPayMethod() {
  showPayMethodModal.value = false
  if (pendingPaymentAction.value) {
    await pendingPaymentAction.value()
    pendingPaymentAction.value = null
  }
}

async function handleRecharge() {
  if (finalAmount.value <= 0 || creatingOrder.value) return

  // Show payment method selection modal
  pendingPaymentAction.value = async () => {
    creatingOrder.value = true
    try {
      const order = await paymentAPI.createRechargeOrder(finalAmount.value, undefined, selectedPayMethod.value)
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
  showPayMethodModal.value = true
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
