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
              <img src="data:image/jpeg;base64,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" class="h-8 w-8" alt="支付宝" />
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
