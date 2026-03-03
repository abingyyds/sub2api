<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-8">
      <!-- Title -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('pricing.title') }}</h1>
        <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('pricing.subtitle') }}</p>
      </div>

      <!-- Plans Grid -->
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <a
          v-for="plan in plans"
          :key="plan.name"
          :href="plan.url"
          target="_blank"
          rel="noopener noreferrer"
          class="card relative flex flex-col overflow-hidden transition-all hover:shadow-lg hover:-translate-y-1"
          :class="{ 'ring-2 ring-primary-500': plan.popular }"
        >
          <div v-if="plan.popular" class="absolute right-4 top-4">
            <span class="badge badge-primary text-xs">{{ t('pricing.recommended') }}</span>
          </div>
          <div class="flex flex-1 flex-col p-6">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">{{ plan.name }}</h3>
            <div class="mt-3 text-3xl font-bold text-primary-600">{{ plan.price }}</div>
            <ul class="mt-5 flex-1 space-y-2.5">
              <li
                v-for="(feature, i) in plan.features"
                :key="i"
                class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300"
              >
                <svg class="h-4 w-4 flex-shrink-0 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                </svg>
                {{ feature }}
              </li>
            </ul>
            <div class="mt-4 flex items-center justify-center rounded-lg bg-primary-600 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700">
              {{ t('pricing.buyNow') }}
            </div>
          </div>
        </a>
      </div>

      <p class="text-center text-sm text-gray-500 dark:text-dark-400">
        {{ t('pricing.stackable') }}
      </p>

      <!-- Purchase Method -->
      <div class="card p-6">
        <h2 class="mb-4 text-xl font-bold text-gray-900 dark:text-white">{{ t('pricing.howToBuy') }}</h2>
        <div class="space-y-4">
          <div class="flex items-start gap-4">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">1</div>
            <div>
              <p class="font-medium text-gray-900 dark:text-white">{{ t('pricing.step1') }}</p>
            </div>
          </div>
          <div class="flex items-start gap-4">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">2</div>
            <p class="font-medium text-gray-900 dark:text-white">{{ t('pricing.step2') }}</p>
          </div>
          <div class="flex items-start gap-4">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">3</div>
            <div>
              <p class="font-medium text-gray-900 dark:text-white">{{ t('pricing.step3') }}</p>
              <router-link to="/redeem" class="text-primary-600 hover:text-primary-700 dark:text-primary-400">
                {{ t('pricing.goRedeem') }}
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Warning -->
      <div class="rounded-lg bg-yellow-50 p-4 dark:bg-yellow-900/20">
        <div class="flex items-start gap-3">
          <svg class="mt-0.5 h-5 w-5 flex-shrink-0 text-yellow-600 dark:text-yellow-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
          </svg>
          <p class="text-sm text-yellow-800 dark:text-yellow-300">{{ t('pricing.warning') }}</p>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'

const { t } = useI18n()

const plans = computed(() => [
  {
    name: t('pricing.plans.payg.name'),
    price: '1:1',
    popular: false,
    features: ['Opus 4.5', 'Sonnet 4.5', 'GPT-5.2-Codex', t('pricing.plans.payg.desc')],
    url: 'https://fk.ccoder.me'
  },
  {
    name: t('pricing.plans.m289.name'),
    price: '¥289',
    popular: true,
    features: [t('pricing.daily', { amount: 30 }), t('pricing.noWeeklyLimit'), t('pricing.validity'), t('pricing.stackableTag')],
    url: 'https://fk.ccoder.me'
  },
  {
    name: t('pricing.plans.m389.name'),
    price: '¥389',
    popular: false,
    features: [t('pricing.daily', { amount: 40 }), t('pricing.noWeeklyLimit'), t('pricing.validity'), t('pricing.stackableTag')],
    url: 'https://fk.ccoder.me'
  },
  {
    name: t('pricing.plans.m459.name'),
    price: '¥459',
    popular: false,
    features: [t('pricing.daily', { amount: 50 }), t('pricing.noWeeklyLimit'), t('pricing.validity'), t('pricing.stackableTag')],
    url: 'https://fk.ccoder.me'
  },
  {
    name: t('pricing.plans.m559.name'),
    price: '¥559',
    popular: false,
    features: [t('pricing.daily', { amount: 60 }), t('pricing.noWeeklyLimit'), t('pricing.validity'), t('pricing.stackableTag')],
    url: 'https://fk.ccoder.me'
  },
  {
    name: t('pricing.plans.m1180.name'),
    price: '¥1180',
    popular: false,
    features: [t('pricing.daily', { amount: 120 }), t('pricing.noWeeklyLimit'), t('pricing.validity'), t('pricing.stackableTag')],
    url: 'https://fk.ccoder.me'
  }
])
</script>
