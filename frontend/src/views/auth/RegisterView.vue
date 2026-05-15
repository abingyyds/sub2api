<template>
  <AuthLayout>
    <div class="space-y-6">
      <!-- Title -->
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('auth.createAccount') }}
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
          {{ t('auth.signUpToStart', { siteName }) }}
        </p>
      </div>

      <!-- LinuxDo Connect OAuth login -->
      <LinuxDoOAuthSection
        v-if="linuxdoOAuthEnabled"
        :disabled="isLoading"
        :legal-accepted="legalAccepted"
        include-legal-acceptance
      />

      <!-- Registration Disabled Message -->
      <div
        v-if="!registrationEnabled && settingsLoaded"
        class="rounded-xl border border-amber-200 bg-amber-50 p-4 dark:border-amber-800/50 dark:bg-amber-900/20"
      >
        <div class="flex items-start gap-3">
          <div class="flex-shrink-0">
            <Icon name="exclamationCircle" size="md" class="text-amber-500" />
          </div>
          <p class="text-sm text-amber-700 dark:text-amber-400">
            {{ t('auth.registrationDisabled') }}
          </p>
        </div>
      </div>

      <!-- Registration Form -->
      <form v-else @submit.prevent="handleRegister" class="space-y-5">
        <!-- Email Input -->
        <div>
          <label for="email" class="input-label">
            {{ t('auth.emailLabel') }}
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
              <Icon name="mail" size="md" class="text-gray-400 dark:text-dark-500" />
            </div>
            <input
              id="email"
              v-model="formData.email"
              type="email"
              required
              autofocus
              autocomplete="email"
              :disabled="isLoading"
              class="input pl-11 transition-all duration-300 focus:ring-2 focus:ring-primary-500"
              :class="{ 'input-error animate-shake': errors.email }"
              :placeholder="t('auth.emailPlaceholder')"
            />
          </div>
          <p v-if="errors.email" class="input-error-text">
            {{ errors.email }}
          </p>
        </div>

        <!-- Password Input -->
        <div>
          <label for="password" class="input-label">
            {{ t('auth.passwordLabel') }}
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
              <Icon name="lock" size="md" class="text-gray-400 dark:text-dark-500" />
            </div>
            <input
              id="password"
              v-model="formData.password"
              :type="showPassword ? 'text' : 'password'"
              required
              autocomplete="new-password"
              :disabled="isLoading"
              class="input pl-11 pr-11 transition-all duration-300 focus:ring-2 focus:ring-primary-500"
              :class="{ 'input-error animate-shake': errors.password }"
              :placeholder="t('auth.createPasswordPlaceholder')"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute inset-y-0 right-0 flex items-center pr-3.5 text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-dark-300"
            >
              <Icon v-if="showPassword" name="eyeOff" size="md" />
              <Icon v-else name="eye" size="md" />
            </button>
          </div>
          <p v-if="errors.password" class="input-error-text">
            {{ errors.password }}
          </p>
          <p v-else class="input-hint">
            {{ t('auth.passwordHint') }}
          </p>
        </div>

        <!-- Agreement review -->
        <div class="space-y-3">
          <p class="text-sm font-medium text-gray-700 dark:text-dark-200">
            {{ t('auth.legal.reviewRequired') }}
          </p>

          <button
            type="button"
            class="flex w-full items-center justify-between rounded-xl border border-gray-200 bg-white px-4 py-3 text-left transition-colors hover:border-primary-300 hover:bg-primary-50/40 dark:border-dark-700 dark:bg-dark-800 dark:hover:border-primary-700 dark:hover:bg-primary-900/20"
            @click="openLegal('terms')"
          >
            <span class="flex items-center gap-3">
              <Icon name="document" size="md" class="text-primary-500" />
              <span>
                <span class="block text-sm font-medium text-gray-900 dark:text-white">
                  {{ t('auth.legal.userAgreement') }}
                </span>
                <span class="block text-xs text-gray-500 dark:text-dark-400">
                  {{ termsAccepted ? t('auth.legal.accepted') : t('auth.legal.readBeforeAccept') }}
                </span>
              </span>
            </span>
            <Icon :name="termsAccepted ? 'checkCircle' : 'chevronRight'" size="md" :class="termsAccepted ? 'text-emerald-500' : 'text-gray-400'" />
          </button>

          <button
            type="button"
            class="flex w-full items-center justify-between rounded-xl border border-gray-200 bg-white px-4 py-3 text-left transition-colors hover:border-primary-300 hover:bg-primary-50/40 dark:border-dark-700 dark:bg-dark-800 dark:hover:border-primary-700 dark:hover:bg-primary-900/20"
            @click="openLegal('privacy')"
          >
            <span class="flex items-center gap-3">
              <Icon name="shield" size="md" class="text-primary-500" />
              <span>
                <span class="block text-sm font-medium text-gray-900 dark:text-white">
                  {{ t('auth.legal.privacyPolicy') }}
                </span>
                <span class="block text-xs text-gray-500 dark:text-dark-400">
                  {{ privacyAccepted ? t('auth.legal.accepted') : t('auth.legal.readBeforeAccept') }}
                </span>
              </span>
            </span>
            <Icon :name="privacyAccepted ? 'checkCircle' : 'chevronRight'" size="md" :class="privacyAccepted ? 'text-emerald-500' : 'text-gray-400'" />
          </button>

          <label
            class="flex cursor-pointer items-start gap-3 rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-sm text-gray-700 dark:border-dark-700 dark:bg-dark-900/60 dark:text-dark-200"
          >
            <input
              v-model="legalCommitmentAccepted"
              type="checkbox"
              class="mt-0.5 h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-600"
              :disabled="!termsAccepted || !privacyAccepted"
            />
            <span>{{ t('auth.legal.commitmentLabel') }}</span>
          </label>

          <p v-if="errors.legal" class="input-error-text">
            {{ errors.legal }}
          </p>
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
        <button
          type="submit"
          :disabled="isLoading || !legalAccepted"
          class="btn btn-primary w-full transition-all duration-300"
        >
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
          <Icon v-else name="userPlus" size="md" class="mr-2" />
          {{
            isLoading
              ? t('auth.processing')
              : emailVerifyEnabled
                ? t('auth.continue')
                : t('auth.createAccount')
          }}
        </button>
      </form>
    </div>

    <!-- Footer -->
    <template #footer>
      <p class="text-gray-500 dark:text-dark-400">
        {{ t('auth.alreadyHaveAccount') }}
        <router-link
          to="/login"
          class="font-medium text-primary-600 transition-colors hover:text-primary-500 dark:text-primary-400 dark:hover:text-primary-300"
        >
          {{ t('auth.signIn') }}
        </router-link>
      </p>
    </template>
  </AuthLayout>

  <LegalAgreementModal
    v-if="activeLegalDocument"
    :show="showLegalModal"
    :document="activeLegalDocument"
    :accept-label="t('auth.legal.acceptDocument')"
    :reject-label="t('common.cancel')"
    :read-prompt="t('auth.legal.scrollToAccept')"
    :accepted-label="t('auth.legal.finishedReading')"
    @accept="handleLegalAccept"
    @reject="showLegalModal = false"
  />
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { AuthLayout } from '@/components/layout'
import Icon from '@/components/icons/Icon.vue'
import LegalAgreementModal from '@/components/legal/LegalAgreementModal.vue'
import { useAuthStore, useAppStore } from '@/stores'
import { redirectAfterAuth } from '@/utils/postAuthRedirect'
import { getLegalDocument, type LegalDocument, type LegalDocumentKind } from '@/legal/documents'

const { t } = useI18n()

const LinuxDoOAuthSection = defineAsyncComponent(
  () => import('@/components/auth/LinuxDoOAuthSection.vue')
)

// ==================== Router & Stores ====================

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()

// ==================== State ====================

const isLoading = ref<boolean>(false)
const settingsLoaded = ref<boolean>(false)
const errorMessage = ref<string>('')
const showPassword = ref<boolean>(false)

// Public settings
const registrationEnabled = ref<boolean>(true)
const emailVerifyEnabled = ref<boolean>(false)
const siteName = ref<string>('cCoder.me')
const linuxdoOAuthEnabled = ref<boolean>(false)

const inviteCode = ref<string>('')
const termsAccepted = ref<boolean>(false)
const privacyAccepted = ref<boolean>(false)
const legalCommitmentAccepted = ref<boolean>(false)
const showLegalModal = ref<boolean>(false)
const activeLegalKind = ref<LegalDocumentKind>('terms')

const formData = reactive({
  email: '',
  password: ''
})

const errors = reactive({
  email: '',
  password: '',
  legal: ''
})

const legalAccepted = computed(
  () => termsAccepted.value && privacyAccepted.value && legalCommitmentAccepted.value
)
const activeLegalDocument = computed<LegalDocument>(() => getLegalDocument(activeLegalKind.value))

// ==================== Lifecycle ====================

onMounted(async () => {
  const inviteParam = route.query.invite as string
  if (inviteParam) {
    inviteCode.value = inviteParam
  }

  try {
    const settings = await appStore.fetchPublicSettings()
    if (settings) {
      registrationEnabled.value = settings.registration_enabled
      emailVerifyEnabled.value = settings.email_verify_enabled
      siteName.value = settings.site_name || 'cCoder.me'
      linuxdoOAuthEnabled.value = settings.linuxdo_oauth_enabled
    }
  } catch (error) {
    console.error('Failed to load public settings:', error)
  } finally {
    settingsLoaded.value = true
  }
})

// ==================== Validation ====================

function validateEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

function validateForm(): boolean {
  // Reset errors
  errors.email = ''
  errors.password = ''
  errors.legal = ''

  let isValid = true

  // Email validation
  if (!formData.email.trim()) {
    errors.email = t('auth.emailRequired')
    isValid = false
  } else if (!validateEmail(formData.email)) {
    errors.email = t('auth.invalidEmail')
    isValid = false
  }

  // Password validation
  if (!formData.password) {
    errors.password = t('auth.passwordRequired')
    isValid = false
  } else if (formData.password.length < 6) {
    errors.password = t('auth.passwordMinLength')
    isValid = false
  }

  if (!legalAccepted.value) {
    errors.legal = t('auth.legal.required')
    isValid = false
  }

  return isValid
}

// ==================== Form Handlers ====================

function openLegal(kind: LegalDocumentKind): void {
  activeLegalKind.value = kind
  showLegalModal.value = true
}

function handleLegalAccept(): void {
  if (activeLegalKind.value === 'terms') {
    termsAccepted.value = true
  } else {
    privacyAccepted.value = true
  }
  errors.legal = ''
  showLegalModal.value = false
}

async function handleRegister(): Promise<void> {
  // Clear previous error
  errorMessage.value = ''

  // Validate form
  if (!validateForm()) {
    return
  }

  isLoading.value = true

  try {
    // If email verification is enabled, redirect to verification page
    if (emailVerifyEnabled.value) {
      // Store registration data in sessionStorage
      sessionStorage.setItem(
        'register_data',
        JSON.stringify({
          email: formData.email,
          password: formData.password,
          invite_code: inviteCode.value || undefined,
          terms_accepted: true,
          privacy_accepted: true,
          legal_commitment_accepted: true
        })
      )

      // Navigate to email verification page
      await router.push('/email-verify')
      return
    }

    // Otherwise, directly register
    await authStore.register({
      email: formData.email,
      password: formData.password,
      invite_code: inviteCode.value || undefined,
      terms_accepted: true,
      privacy_accepted: true,
      legal_commitment_accepted: true
    })

    // Show success toast
    appStore.showSuccess(t('auth.accountCreatedSuccess', { siteName: siteName.value }))

    // Redirect to pricing page so new users can purchase/recharge
    await redirectAfterAuth(router, '/pricing')
  } catch (error: unknown) {
    // Handle registration error
    const err = error as { message?: string; response?: { data?: { detail?: string } } }

    if (err.response?.data?.detail) {
      errorMessage.value = err.response.data.detail
    } else if (err.message) {
      errorMessage.value = err.message
    } else {
      errorMessage.value = t('auth.registrationFailed')
    }

    // Also show error toast
    appStore.showError(errorMessage.value)
  } finally {
    isLoading.value = false
  }
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
