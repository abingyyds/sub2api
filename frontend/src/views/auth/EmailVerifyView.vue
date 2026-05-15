<template>
  <AuthLayout>
    <div class="space-y-6">
      <!-- Title -->
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('auth.verifyYourEmail') }}
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
          We'll send a verification code to
          <span class="font-medium text-gray-700 dark:text-gray-300">{{ email }}</span>
        </p>
      </div>

      <!-- No Data Warning -->
      <div
        v-if="!hasRegisterData"
        class="rounded-xl border border-amber-200 bg-amber-50 p-4 dark:border-amber-800/50 dark:bg-amber-900/20"
      >
        <div class="flex items-start gap-3">
          <div class="flex-shrink-0">
            <Icon name="exclamationCircle" size="md" class="text-amber-500" />
          </div>
          <div class="text-sm text-amber-700 dark:text-amber-400">
            <p class="font-medium">{{ t('auth.sessionExpired') }}</p>
            <p class="mt-1">{{ t('auth.sessionExpiredDesc') }}</p>
          </div>
        </div>
      </div>

      <!-- Verification Form -->
      <form v-else @submit.prevent="handleVerify" class="space-y-5">
        <!-- Verification Code Input -->
        <div>
          <label for="code" class="input-label text-center">
            {{ t('auth.verificationCode') }}
          </label>
          <input
            id="code"
            v-model="verifyCode"
            type="text"
            required
            autocomplete="one-time-code"
            inputmode="numeric"
            maxlength="6"
            :disabled="isLoading"
            class="input py-3 text-center font-mono text-xl tracking-[0.5em]"
            :class="{ 'input-error': errors.code }"
            placeholder="000000"
          />
          <p v-if="errors.code" class="input-error-text text-center">
            {{ errors.code }}
          </p>
          <p v-else class="input-hint text-center">{{ t('auth.verificationCodeHint') }}</p>
        </div>

        <!-- Code Status -->
        <div
          v-if="codeSent"
          class="rounded-xl border border-green-200 bg-green-50 p-4 dark:border-green-800/50 dark:bg-green-900/20"
        >
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0">
              <Icon name="checkCircle" size="md" class="text-green-500" />
            </div>
            <p class="text-sm text-green-700 dark:text-green-400">
              Verification code sent! Please check your inbox.
            </p>
          </div>
        </div>

        <!-- Error Message -->
        <transition name="fade">
          <div
            v-if="errorMessage"
            class="rounded-xl border border-red-200 bg-red-50 p-4 dark:border-red-800/50 dark:bg-red-900/20"
          >
            <div class="flex items-start gap-3">
              <div class="flex-shrink-0">
                <Icon name="exclamationCircle" size="md" class="text-red-500" />
              </div>
              <p class="text-sm text-red-700 dark:text-red-400">
                {{ errorMessage }}
              </p>
            </div>
          </div>
        </transition>

        <!-- Submit Button -->
        <button type="submit" :disabled="isLoading || !verifyCode" class="btn btn-primary w-full">
          <svg
            v-if="isLoading"
            class="-ml-1 mr-2 h-4 w-4 animate-spin text-white"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          <Icon v-else name="checkCircle" size="md" class="mr-2" />
          {{ isLoading ? 'Verifying...' : 'Verify & Create Account' }}
        </button>

        <!-- Resend Code -->
        <div class="text-center">
          <button
            v-if="countdown > 0"
            type="button"
            disabled
            class="cursor-not-allowed text-sm text-gray-400 dark:text-dark-500"
          >
            Resend code in {{ countdown }}s
          </button>
          <button
            v-else
            type="button"
            @click="handleResendCode"
            :disabled="isSendingCode"
            class="text-sm text-primary-600 transition-colors hover:text-primary-500 disabled:cursor-not-allowed disabled:opacity-50 dark:text-primary-400 dark:hover:text-primary-300"
          >
            <span v-if="isSendingCode">{{ t('auth.sendingCode') }}</span>
            <span v-else>{{ t('auth.resendCode') }}</span>
          </button>
        </div>
      </form>
    </div>

    <!-- Footer -->
    <template #footer>
      <button
        @click="handleBack"
        class="flex items-center gap-2 text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-gray-300"
      >
        <Icon name="arrowLeft" size="sm" />
        Back to registration
      </button>
    </template>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { AuthLayout } from '@/components/layout'
import Icon from '@/components/icons/Icon.vue'
import { useAuthStore, useAppStore } from '@/stores'
import { sendVerifyCode } from '@/api/auth'
import { redirectAfterAuth } from '@/utils/postAuthRedirect'

const { t } = useI18n()

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

const isLoading = ref<boolean>(false)
const isSendingCode = ref<boolean>(false)
const errorMessage = ref<string>('')
const codeSent = ref<boolean>(false)
const verifyCode = ref<string>('')
const countdown = ref<number>(0)
let countdownTimer: ReturnType<typeof setInterval> | null = null

const email = ref<string>('')
const password = ref<string>('')
const inviteCode = ref<string>('')
const termsAccepted = ref<boolean>(false)
const privacyAccepted = ref<boolean>(false)
const legalCommitmentAccepted = ref<boolean>(false)
const hasRegisterData = ref<boolean>(false)

const siteName = ref<string>('cCoder.me')

const errors = ref({
  code: ''
})

onMounted(async () => {
  const registerDataStr = sessionStorage.getItem('register_data')
  if (registerDataStr) {
    try {
      const registerData = JSON.parse(registerDataStr)
      email.value = registerData.email || ''
      password.value = registerData.password || ''
      inviteCode.value = registerData.invite_code || ''
      termsAccepted.value = registerData.terms_accepted === true
      privacyAccepted.value = registerData.privacy_accepted === true
      legalCommitmentAccepted.value = registerData.legal_commitment_accepted === true
      hasRegisterData.value = !!(
        email.value &&
        password.value &&
        termsAccepted.value &&
        privacyAccepted.value &&
        legalCommitmentAccepted.value
      )
    } catch {
      hasRegisterData.value = false
    }
  }

  try {
    const settings = await appStore.fetchPublicSettings()
    if (settings) {
      siteName.value = settings.site_name || 'cCoder.me'
    }
  } catch (error) {
    console.error('Failed to load public settings:', error)
  }

  if (hasRegisterData.value) {
    await sendCode()
  }
})

onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
})

function startCountdown(seconds: number): void {
  countdown.value = seconds

  if (countdownTimer) {
    clearInterval(countdownTimer)
  }

  countdownTimer = setInterval(() => {
    if (countdown.value > 0) {
      countdown.value--
    } else {
      if (countdownTimer) {
        clearInterval(countdownTimer)
        countdownTimer = null
      }
    }
  }, 1000)
}

async function sendCode(): Promise<void> {
  isSendingCode.value = true
  errorMessage.value = ''

  try {
    const response = await sendVerifyCode({
      email: email.value
    })

    codeSent.value = true
    startCountdown(response.countdown)
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { detail?: string } } }

    if (err.response?.data?.detail) {
      errorMessage.value = err.response.data.detail
    } else if (err.message) {
      errorMessage.value = err.message
    } else {
      errorMessage.value = 'Failed to send verification code. Please try again.'
    }

    appStore.showError(errorMessage.value)
  } finally {
    isSendingCode.value = false
  }
}

async function handleResendCode(): Promise<void> {
  await sendCode()
}

function validateForm(): boolean {
  errors.value.code = ''

  if (!verifyCode.value.trim()) {
    errors.value.code = 'Verification code is required'
    return false
  }

  if (!/^\d{6}$/.test(verifyCode.value.trim())) {
    errors.value.code = 'Please enter a valid 6-digit code'
    return false
  }

  return true
}

async function handleVerify(): Promise<void> {
  errorMessage.value = ''

  if (!validateForm()) {
    return
  }

  isLoading.value = true

  try {
    await authStore.register({
      email: email.value,
      password: password.value,
      verify_code: verifyCode.value.trim(),
      invite_code: inviteCode.value || undefined,
      terms_accepted: termsAccepted.value,
      privacy_accepted: privacyAccepted.value,
      legal_commitment_accepted: legalCommitmentAccepted.value
    })

    sessionStorage.removeItem('register_data')

    appStore.showSuccess('Account created successfully! Welcome to ' + siteName.value + '.')

    await redirectAfterAuth(router, '/pricing')
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { detail?: string } } }

    if (err.response?.data?.detail) {
      errorMessage.value = err.response.data.detail
    } else if (err.message) {
      errorMessage.value = err.message
    } else {
      errorMessage.value = 'Verification failed. Please try again.'
    }

    appStore.showError(errorMessage.value)
  } finally {
    isLoading.value = false
  }
}

function handleBack(): void {
  sessionStorage.removeItem('register_data')
  router.push('/register')
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
