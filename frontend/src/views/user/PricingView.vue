<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-6xl space-y-8">
        <!-- Title -->
        <SlideIn direction="up" :delay="100">
          <div class="text-center">
            <h1 class="text-3xl md:text-4xl font-extrabold text-gray-900 dark:text-white tracking-tight">{{ t('pricing.title') }}</h1>
            <p class="mt-3 text-gray-500 dark:text-dark-400 text-lg">{{ t('pricing.subtitle') }}</p>
          </div>
        </SlideIn>

        <!-- Tabs -->
        <SlideIn direction="up" :delay="150">
          <div class="flex justify-center">
            <div class="inline-flex rounded-full border border-gray-200 bg-white p-1 shadow-sm dark:border-dark-600 dark:bg-dark-800">
              <button
                class="rounded-full px-6 py-2.5 text-sm font-bold transition-all"
                :class="activeTab === 'subscription'
                  ? 'bg-primary-600 text-white shadow-md'
                  : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200'"
                @click="activeTab = 'subscription'"
              >
                {{ t('pricing.tabSubscription') }}
              </button>
              <button
                class="rounded-full px-6 py-2.5 text-sm font-bold transition-all"
                :class="activeTab === 'recharge'
                  ? 'bg-primary-600 text-white shadow-md'
                  : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200'"
                @click="activeTab = 'recharge'"
              >
                {{ t('pricing.tabRecharge') }}
              </button>
            </div>
          </div>
        </SlideIn>

        <!-- Loading -->
        <div v-if="loading" class="flex justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
        </div>

        <!-- Promo Code Input (shared across tabs) -->
        <SlideIn v-if="!loading" direction="up" :delay="200">
          <div class="mx-auto max-w-lg">
            <GlowCard glow-color="rgb(34, 197, 94)">
              <div class="rounded-2xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800 shadow-sm">
                <div class="flex items-center gap-3">
                  <div class="relative flex-1">
                    <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                      <Icon name="gift" size="md" :class="promoValidation.valid ? 'text-green-500' : 'text-gray-400 dark:text-dark-500'" />
                    </div>
                    <input
                      v-model="promoCode"
                      type="text"
                      class="w-full rounded-xl border border-gray-200 bg-gray-50 py-3 pl-10 pr-10 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:bg-white focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-dark-500 dark:focus:bg-dark-800 transition-all"
                      :placeholder="t('pricing.promoCodePlaceholder')"
                      @input="handlePromoCodeInput"
                    />
                    <div v-if="promoValidating" class="absolute inset-y-0 right-0 flex items-center pr-3">
                      <svg class="h-4 w-4 animate-spin text-gray-400" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                    </div>
                    <div v-else-if="promoValidation.valid" class="absolute inset-y-0 right-0 flex items-center pr-3">
                      <Icon name="checkCircle" size="md" class="text-green-500" />
                    </div>
                    <div v-else-if="promoValidation.invalid" class="absolute inset-y-0 right-0 flex items-center pr-3">
                      <Icon name="exclamationCircle" size="md" class="text-red-500" />
                    </div>
                  </div>
                </div>
                <!-- Promo code validation result -->
                <transition name="fade">
                  <div v-if="promoValidation.valid" class="mt-3 flex items-center gap-2 rounded-xl bg-green-50 px-4 py-2.5 dark:bg-green-900/20 border border-green-100 dark:border-green-800/30">
                    <Icon name="gift" size="sm" class="text-green-600 dark:text-green-400" />
                    <span class="text-sm font-medium text-green-700 dark:text-green-400">
                      {{ formatDiscountText() }}
                    </span>
                  </div>
                  <p v-else-if="promoValidation.invalid" class="mt-2 text-sm text-red-500 pl-1">
                    {{ promoValidation.message }}
                  </p>
                </transition>
              </div>
            </GlowCard>
          </div>
        </SlideIn>

        <!-- ==================== Subscription Tab ==================== -->
        <template v-if="!loading && activeTab === 'subscription'">
          <StaggerContainer v-if="decoratedPlans.length > 0" :stagger-delay="120" :delay="250">
            <div class="grid gap-8 lg:grid-cols-3 items-stretch">
              <GlowCard
                v-for="dp in decoratedPlans"
                :key="dp.plan.key"
                :glow-color="dp.featured ? 'rgb(217, 119, 87)' : dp.plan.amount_fen >= 100000 ? 'rgb(168, 85, 247)' : 'rgb(234, 179, 8)'"
              >
                <div
                  class="relative flex flex-col overflow-hidden rounded-3xl border-2 bg-white dark:bg-dark-900 transition-all hover:shadow-lg cursor-pointer h-full"
                  :class="dp.featured
                    ? 'border-[3px] border-primary-500/80 dark:border-primary-400/60 lg:scale-105 shadow-2xl shadow-primary-500/20 z-10'
                    : dp.plan.amount_fen >= 100000
                      ? 'border-purple-100 dark:border-purple-800/30 shadow-soft'
                      : 'border-orange-100 dark:border-orange-800/30 shadow-soft'"
                  @click="handleBuy(dp.plan)"
                >
                  <!-- Top badge -->
                  <div
                    v-if="dp.badge"
                    class="absolute -top-3.5 text-xs font-bold px-4 py-1 rounded-full shadow-sm whitespace-nowrap"
                    :class="dp.featured
                      ? 'left-1/2 -translate-x-1/2 -top-4 bg-gradient-to-r from-orange-500 to-red-500 text-white px-5 py-1.5 shadow-lg tracking-wide'
                      : dp.plan.amount_fen >= 100000
                        ? 'right-8 bg-purple-50 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 border border-purple-200 dark:border-purple-700'
                        : 'left-8 bg-orange-50 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400 border border-orange-200 dark:border-orange-700'"
                  >
                    {{ dp.badge }}
                  </div>

                  <div class="flex flex-1 flex-col p-8">
                    <!-- Target audience -->
                    <p class="text-sm font-semibold mb-2 mt-2"
                      :class="dp.featured ? 'text-primary-600 dark:text-primary-400' : 'text-gray-500 dark:text-dark-400'"
                    >{{ dp.audience }}</p>

                    <!-- Plan name -->
                    <h3 class="font-bold text-gray-900 dark:text-white mb-2"
                      :class="dp.featured ? 'text-3xl' : 'text-2xl'"
                    >{{ dp.plan.name }}</h3>

                    <!-- Price -->
                    <div class="mb-4"
                      :class="dp.featured ? 'text-5xl' : 'text-4xl'"
                    >
                      <span class="font-extrabold text-gray-900 dark:text-white">¥{{ (dp.plan.amount_fen / 100).toFixed(0) }}</span>
                    </div>

                    <!-- Highlight text -->
                    <div v-if="dp.highlight"
                      class="rounded-xl p-2.5 mb-6 text-sm text-center font-bold"
                      :class="dp.featured
                        ? 'bg-primary-500/10 dark:bg-primary-400/10 border border-primary-500/20 dark:border-primary-400/20 text-primary-600 dark:text-primary-400'
                        : 'bg-gray-50 dark:bg-dark-800 border border-gray-100 dark:border-dark-700 text-gray-600 dark:text-dark-400'"
                    >
                      {{ dp.highlight }}
                    </div>

                    <!-- Feature list -->
                    <ul class="flex-1 space-y-4 mb-8 text-sm text-gray-600 dark:text-dark-400">
                      <li v-for="(feat, fi) in dp.features" :key="fi" class="flex items-start gap-3">
                        <Icon
                          :name="(dp.featured ? getFeatureIcon(fi) : 'check') as any"
                          size="sm"
                          :class="dp.featured
                            ? getFeatureIconColor(fi)
                            : dp.plan.amount_fen >= 100000
                              ? 'text-gray-700 dark:text-dark-300 mt-0.5 flex-shrink-0'
                              : 'text-green-500 mt-0.5 flex-shrink-0'"
                          class="mt-0.5 flex-shrink-0"
                        />
                        <div class="leading-relaxed" v-html="feat"></div>
                      </li>
                    </ul>

                    <!-- CTA button -->
                    <MagneticButton>
                      <button
                        class="w-full rounded-xl font-bold transition-all duration-200"
                        :class="dp.featured
                          ? 'py-4 bg-primary-600 hover:bg-primary-700 text-white text-lg shadow-[0_4px_14px_0_rgba(217,119,87,0.39)] transform hover:-translate-y-0.5'
                          : dp.plan.amount_fen >= 100000
                            ? 'py-3.5 bg-purple-50 dark:bg-purple-900/20 hover:bg-purple-100 dark:hover:bg-purple-900/30 text-purple-700 dark:text-purple-400 border border-purple-200 dark:border-purple-700'
                            : 'py-3.5 bg-orange-50 dark:bg-orange-900/20 hover:bg-orange-100 dark:hover:bg-orange-900/30 text-orange-700 dark:text-orange-400 border border-orange-200 dark:border-orange-700'"
                        :disabled="creatingOrder"
                      >
                        {{ dp.cta }}
                      </button>
                    </MagneticButton>
                  </div>
                </div>
              </GlowCard>
            </div>
          </StaggerContainer>

          <div v-else class="text-center py-12 text-gray-500 dark:text-dark-400">
            {{ t('pricing.noPlans') }}
          </div>
          <p v-if="decoratedPlans.length > 0" class="text-center text-sm text-gray-500 dark:text-dark-400">
            {{ t('pricing.stackable') }}
          </p>
        </template>

        <!-- ==================== Recharge Tab ==================== -->
        <template v-if="!loading && activeTab === 'recharge'">
          <!-- Recharge Plans (if configured) -->
          <StaggerContainer v-if="rechargePlans.length > 0" :stagger-delay="100" :delay="250">
            <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
              <GlowCard
                v-for="rp in rechargePlans"
                :key="rp.key"
                :glow-color="rp.popular ? 'rgb(217, 119, 87)' : 'rgb(59, 130, 246)'"
              >
                <div
                  class="relative flex flex-col overflow-hidden rounded-2xl border-2 bg-white dark:bg-dark-900 transition-all hover:shadow-lg hover:-translate-y-1 cursor-pointer h-full"
                  :class="rp.popular
                    ? 'border-primary-500/80 dark:border-primary-400/60 shadow-lg shadow-primary-500/10'
                    : 'border-gray-200 dark:border-dark-600 shadow-soft'"
                  @click="handleRechargePreset(rp)"
                >
                  <div v-if="rp.popular" class="absolute -top-3 left-1/2 -translate-x-1/2">
                    <span class="bg-gradient-to-r from-orange-500 to-red-500 text-white text-xs font-bold px-4 py-1 rounded-full shadow-md whitespace-nowrap">
                      {{ t('pricing.mostPopular') }}
                    </span>
                  </div>
                  <div class="flex flex-1 flex-col p-6">
                    <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ rp.name }}</h3>
                    <div class="mt-3 flex items-baseline gap-1">
                      <span class="text-xs text-gray-500 dark:text-dark-400">{{ t('pricing.rechargeGet') }}</span>
                      <span class="text-3xl font-extrabold text-primary-600">${{ rp.balance_amount }}</span>
                    </div>
                    <p v-if="rp.description" class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ rp.description }}</p>
                    <div class="mt-auto pt-5">
                      <div v-if="rp.pay_amount_fen !== rp.balance_amount * 100" class="mb-1.5 text-center">
                        <span class="text-xs text-gray-400 line-through dark:text-dark-500">¥{{ rp.balance_amount * 100 / 100 }}</span>
                      </div>
                      <div
                        class="flex items-center justify-center rounded-xl py-2.5 text-sm font-bold text-white transition-all"
                        :class="rp.popular
                          ? 'bg-primary-600 hover:bg-primary-700 shadow-md'
                          : 'bg-gray-900 hover:bg-gray-800 dark:bg-dark-500 dark:hover:bg-dark-400'"
                      >
                        {{ t('pricing.rechargeNow') }} ¥{{ (rp.pay_amount_fen / 100).toFixed(rp.pay_amount_fen % 100 === 0 ? 0 : 2) }}
                      </div>
                    </div>
                  </div>
                </div>
              </GlowCard>
            </div>
          </StaggerContainer>

          <!-- Custom Amount Input -->
          <SlideIn direction="up" :delay="350">
            <div class="mx-auto max-w-lg">
              <GlowCard glow-color="rgb(99, 102, 241)">
                <div class="rounded-2xl border-2 border-gray-200 bg-white p-6 space-y-5 dark:border-dark-700 dark:bg-dark-900 shadow-soft">
                  <h3 class="text-center text-xl font-bold text-gray-900 dark:text-white">{{ t('pricing.customRecharge') }}</h3>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-2">
                      {{ t('recharge.customAmount') }}
                    </label>
                    <div class="relative">
                      <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 dark:text-dark-400 font-bold text-lg">¥</span>
                      <input
                        v-model="customInput"
                        type="number"
                        min="1"
                        step="1"
                        :placeholder="rechargeMinAmount > 0 ? t('pricing.minAmountHint', { amount: rechargeMinAmount }) : t('recharge.inputPlaceholder')"
                        class="w-full rounded-xl border-2 border-gray-200 bg-gray-50 py-3.5 pl-10 pr-4 text-lg font-bold text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:bg-white focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-800 dark:text-white dark:placeholder-dark-500 dark:focus:bg-dark-700 transition-all"
                      />
                    </div>
                    <p v-if="rechargeMinAmount > 0" class="mt-1.5 text-xs text-gray-400 dark:text-dark-500">
                      {{ t('pricing.minAmountNote', { amount: rechargeMinAmount }) }}
                    </p>
                  </div>

                  <!-- Amount Summary -->
                  <div v-if="customFinalAmount > 0" class="rounded-xl bg-primary-50 dark:bg-primary-900/10 border border-primary-100 dark:border-primary-800/30 p-4">
                    <div class="flex items-center justify-between">
                      <span class="text-sm font-medium text-gray-600 dark:text-dark-400">{{ t('recharge.rechargeAmount') }}</span>
                      <span class="text-xl font-extrabold text-primary-600">¥{{ customFinalAmount.toFixed(2) }}</span>
                    </div>
                  </div>

                  <!-- Submit -->
                  <MagneticButton>
                    <button
                      class="w-full rounded-xl bg-primary-600 py-3.5 text-base font-bold text-white transition-all hover:bg-primary-700 hover:-translate-y-0.5 shadow-md disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                      :disabled="customFinalAmount <= 0 || (rechargeMinAmount > 0 && customFinalAmount < rechargeMinAmount) || creatingOrder"
                      @click="handleCustomRecharge"
                    >
                      <span v-if="creatingOrder" class="flex items-center justify-center gap-2">
                        <div class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
                        {{ t('recharge.creating') }}
                      </span>
                      <span v-else>
                        {{ customFinalAmount > 0 ? t('recharge.payNow', { amount: customFinalAmount.toFixed(2) }) : t('recharge.enterAmount') }}
                      </span>
                    </button>
                  </MagneticButton>
                </div>
              </GlowCard>
            </div>
          </SlideIn>
        </template>
      </div>
    </FadeIn>

    <!-- Payment Method Selection Modal -->
    <Teleport to="body">
      <transition name="fade">
        <div
          v-if="showPayMethodModal"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
          @click.self="showPayMethodModal = false"
        >
          <div class="w-full max-w-md rounded-2xl bg-white p-8 shadow-2xl dark:bg-dark-800 mx-4 border border-gray-100 dark:border-dark-700">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-6 text-center">选择支付方式</h3>

            <div class="space-y-3 mb-8">
              <!-- WeChat Pay -->
              <div
                class="flex items-center gap-4 rounded-xl border-2 p-4 cursor-pointer transition-all"
                :class="selectedPayMethod === 'wechat' ? 'border-green-500 bg-green-50 dark:bg-green-900/20 shadow-sm' : 'border-gray-200 dark:border-dark-600 hover:border-green-300'"
                @click="selectedPayMethod = 'wechat'"
              >
                <input type="radio" :checked="selectedPayMethod === 'wechat'" class="h-5 w-5 text-green-600" />
                <svg class="h-8 w-8" viewBox="0 0 24 24" fill="none">
                  <rect width="24" height="24" rx="4" fill="#07C160"/>
                  <path d="M16.7 10.5c-.2 0-.4 0-.6.1.1-.4.1-.8.1-1.2 0-2.8-2.7-5-5.9-5C7 4.4 4.4 6.6 4.4 9.4c0 1.5.8 2.9 2.1 3.9l-.5 1.6 1.9-1c.6.2 1.2.3 1.8.3h.2c-.1-.3-.1-.7-.1-1 0-2.4 2.3-4.4 5.1-4.4l.2-.1c-.1-.1-.2-.2-.4-.2zm-3.3-2.2c.4 0 .7.3.7.7s-.3.7-.7.7-.7-.3-.7-.7.3-.7.7-.7zm-4.2 1.4c-.4 0-.7-.3-.7-.7s.3-.7.7-.7.7.3.7.7-.3.7-.7.7zm10.4 4.5c0-2.2-2.2-4-4.8-4s-4.8 1.8-4.8 4 2.2 4 4.8 4c.5 0 1-.1 1.5-.2l1.5.8-.4-1.3c1.2-.8 2.2-2 2.2-3.3zm-6.4-.6c-.3 0-.6-.3-.6-.6s.3-.6.6-.6.6.3.6.6-.3.6-.6.6zm3.2 0c-.3 0-.6-.3-.6-.6s.3-.6.6-.6.6.3.6.6-.3.6-.6.6z" fill="white"/>
                </svg>
                <span class="font-bold text-gray-900 dark:text-white">微信支付</span>
              </div>

              <!-- Alipay -->
              <div
                class="flex items-center gap-4 rounded-xl border-2 p-4 cursor-pointer transition-all"
                :class="selectedPayMethod === 'alipay' ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 shadow-sm' : 'border-gray-200 dark:border-dark-600 hover:border-blue-300'"
                @click="selectedPayMethod = 'alipay'"
              >
                <input type="radio" :checked="selectedPayMethod === 'alipay'" class="h-5 w-5 text-blue-600" />
                <svg class="h-8 w-8" viewBox="0 0 24 24" fill="none">
                  <rect width="24" height="24" rx="4" fill="#1677FF"/>
                  <path d="M17.5 14.2c-1.2-.5-2.3-1-3.2-1.4.4-.8.7-1.7.9-2.6h-2.5v-1h3V8.4h-3V6.5h-1.4v1.9h-3v.8h3v1h-2.5v.8h4.6c-.2.7-.5 1.3-.8 1.9-1.3-.5-2.7-.8-3.8-.8-2 0-3.3.9-3.3 2.3 0 1.5 1.3 2.3 3.2 2.3 1.5 0 2.9-.6 4-1.5.9.5 1.9 1 3 1.5l.8-1.2zM9.2 16.3c-1.3 0-1.9-.5-1.9-1.2 0-.8.7-1.3 1.9-1.3.9 0 1.9.2 2.9.7-.9.8-1.9 1.8-2.9 1.8z" fill="white"/>
                </svg>
                <span class="font-bold text-gray-900 dark:text-white">支付宝</span>
              </div>

              <!-- Epay-支付宝 -->
              <div
                class="flex items-center gap-4 rounded-xl border-2 p-4 cursor-pointer transition-all"
                :class="selectedPayMethod === 'epay_alipay' ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 shadow-sm' : 'border-gray-200 dark:border-dark-600 hover:border-blue-300'"
                @click="selectedPayMethod = 'epay_alipay'"
              >
                <input type="radio" :checked="selectedPayMethod === 'epay_alipay'" class="h-5 w-5 text-blue-600" />
                <svg class="h-8 w-8" viewBox="0 0 24 24" fill="none">
                  <rect width="24" height="24" rx="4" fill="#1677FF"/>
                  <path d="M17.5 14.2c-1.2-.5-2.3-1-3.2-1.4.4-.8.7-1.7.9-2.6h-2.5v-1h3V8.4h-3V6.5h-1.4v1.9h-3v.8h3v1h-2.5v.8h4.6c-.2.7-.5 1.3-.8 1.9-1.3-.5-2.7-.8-3.8-.8-2 0-3.3.9-3.3 2.3 0 1.5 1.3 2.3 3.2 2.3 1.5 0 2.9-.6 4-1.5.9.5 1.9 1 3 1.5l.8-1.2zM9.2 16.3c-1.3 0-1.9-.5-1.9-1.2 0-.8.7-1.3 1.9-1.3.9 0 1.9.2 2.9.7-.9.8-1.9 1.8-2.9 1.8z" fill="white"/>
                </svg>
                <span class="font-bold text-gray-900 dark:text-white">Epay-支付宝</span>
              </div>

              <!-- Epay-微信 -->
              <div
                class="flex items-center gap-4 rounded-xl border-2 p-4 cursor-pointer transition-all"
                :class="selectedPayMethod === 'epay_wxpay' ? 'border-green-500 bg-green-50 dark:bg-green-900/20 shadow-sm' : 'border-gray-200 dark:border-dark-600 hover:border-green-300'"
                @click="selectedPayMethod = 'epay_wxpay'"
              >
                <input type="radio" :checked="selectedPayMethod === 'epay_wxpay'" class="h-5 w-5 text-green-600" />
                <svg class="h-8 w-8" viewBox="0 0 24 24" fill="none">
                  <rect width="24" height="24" rx="4" fill="#07C160"/>
                  <path d="M16.7 10.5c-.2 0-.4 0-.6.1.1-.4.1-.8.1-1.2 0-2.8-2.7-5-5.9-5C7 4.4 4.4 6.6 4.4 9.4c0 1.5.8 2.9 2.1 3.9l-.5 1.6 1.9-1c.6.2 1.2.3 1.8.3h.2c-.1-.3-.1-.7-.1-1 0-2.4 2.3-4.4 5.1-4.4l.2-.1c-.1-.1-.2-.2-.4-.2zm-3.3-2.2c.4 0 .7.3.7.7s-.3.7-.7.7-.7-.3-.7-.7.3-.7.7-.7zm-4.2 1.4c-.4 0-.7-.3-.7-.7s.3-.7.7-.7.7.3.7.7-.3.7-.7.7zm10.4 4.5c0-2.2-2.2-4-4.8-4s-4.8 1.8-4.8 4 2.2 4 4.8 4c.5 0 1-.1 1.5-.2l1.5.8-.4-1.3c1.2-.8 2.2-2 2.2-3.3zm-6.4-.6c-.3 0-.6-.3-.6-.6s.3-.6.6-.6.6.3.6.6-.3.6-.6.6zm3.2 0c-.3 0-.6-.3-.6-.6s.3-.6.6-.6.6.3.6.6-.3.6-.6.6z" fill="white"/>
                </svg>
                <span class="font-bold text-gray-900 dark:text-white">Epay-微信</span>
              </div>
            </div>

            <div class="flex gap-3">
              <button
                class="flex-1 rounded-xl border-2 border-gray-200 dark:border-dark-600 py-3 font-bold text-gray-700 dark:text-dark-300 hover:bg-gray-50 dark:hover:bg-dark-700 transition"
                @click="showPayMethodModal = false"
              >
                取消
              </button>
              <MagneticButton class="flex-1">
                <button
                  class="w-full rounded-xl bg-primary-600 py-3 font-bold text-white hover:bg-primary-700 transition shadow-md"
                  @click="confirmPayMethod"
                >
                  确认支付
                </button>
              </MagneticButton>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Payment QR Code Modal -->
    <Teleport to="body">
      <transition name="fade">
        <div
          v-if="showPaymentModal"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
          @click.self="cancelPayment"
        >
          <div class="w-full max-w-md rounded-2xl bg-white p-8 shadow-2xl dark:bg-dark-800 mx-4 border border-gray-100 dark:border-dark-700">
            <!-- Header -->
            <div class="mb-6 text-center">
              <h3 class="text-xl font-bold text-gray-900 dark:text-white">{{ t('pricing.payment.scanToPay') }}</h3>
              <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
                {{ t('pricing.payment.amount') }}: <span class="text-xl font-extrabold text-primary-600">¥{{ currentOrderAmount }}</span>
              </p>
            </div>

            <!-- QR Code -->
            <div class="flex justify-center mb-6">
              <div v-if="qrLoading" class="flex h-48 w-48 items-center justify-center">
                <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
              </div>
              <div v-else class="rounded-xl border-2 border-gray-100 dark:border-dark-700 p-2">
                <canvas ref="qrCanvas" class="rounded-lg"></canvas>
              </div>
            </div>

            <!-- Status -->
            <div class="mb-4 text-center">
              <div v-if="paymentStatus === 'pending'" class="flex items-center justify-center gap-2 text-sm text-gray-500 dark:text-dark-400">
                <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
                {{ t('pricing.payment.waitingPayment') }}
              </div>
              <div v-else-if="paymentStatus === 'paid'" class="flex items-center justify-center gap-2 text-sm font-bold text-green-600">
                <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ paymentOrderType === 'subscription' ? t('pricing.payment.paymentSuccess') : t('recharge.payment.success') }}
              </div>
              <div v-else-if="paymentStatus === 'closed'" class="text-sm font-bold text-red-500">
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
                class="flex-1 rounded-xl border-2 border-gray-200 dark:border-dark-600 py-3 font-bold text-gray-700 dark:text-dark-300 hover:bg-gray-50 dark:hover:bg-dark-700 transition"
                @click="cancelPayment"
              >
                {{ t('pricing.payment.cancel') }}
              </button>
              <MagneticButton v-if="paymentStatus === 'paid'" class="flex-1">
                <button
                  class="w-full rounded-xl bg-primary-600 py-3 font-bold text-white hover:bg-primary-700 transition shadow-md"
                  @click="goAfterPayment"
                >
                  {{ paymentOrderType === 'subscription' ? t('pricing.payment.viewSubscription') : t('recharge.payment.viewDashboard') }}
                </button>
              </MagneticButton>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import QRCode from 'qrcode'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { paymentAPI } from '@/api/payment'
import { validatePromoCode } from '@/api/auth'
import type { PaymentPlan, RechargePlan } from '@/api/payment'
import { FadeIn, SlideIn, StaggerContainer, GlowCard, MagneticButton } from '@/components/animations'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

// Tab state
const activeTab = ref<'subscription' | 'recharge'>(
  route.query.tab === 'recharge' ? 'recharge' : 'subscription'
)

const loading = ref(true)
const plans = ref<PaymentPlan[]>([])
const rechargePlans = ref<RechargePlan[]>([])
const rechargeMinAmount = ref(0)
const creatingOrder = ref(false)
const selectedPayMethod = ref<'wechat' | 'alipay' | 'epay_alipay' | 'epay_wxpay'>('wechat')
const showPayMethodModal = ref(false)
const showPaymentModal = ref(false)
const pendingPaymentAction = ref<(() => Promise<void>) | null>(null)
const qrLoading = ref(false)
const qrCanvas = ref<HTMLCanvasElement | null>(null)
const paymentStatus = ref<'pending' | 'paid' | 'closed'>('pending')
const paymentOrderType = ref<'subscription' | 'recharge'>('subscription')
const currentOrderNo = ref('')
const currentOrderAmount = ref('')
const countdown = ref(0)
const customInput = ref('')

// Promo code state
const promoCode = ref('')
const promoValidating = ref(false)
const promoValidation = ref<{
  valid: boolean
  invalid: boolean
  discountType: string
  discountAmount: number
  discountFen: number
  finalAmountFen: number
  minOrderAmount: number
  message: string
}>({
  valid: false,
  invalid: false,
  discountType: '',
  discountAmount: 0,
  discountFen: 0,
  finalAmountFen: 0,
  minOrderAmount: 0,
  message: ''
})
let promoValidateTimeout: ReturnType<typeof setTimeout> | null = null

let pollTimer: ReturnType<typeof setInterval> | null = null
let countdownTimer: ReturnType<typeof setInterval> | null = null

const customFinalAmount = computed(() => {
  if (customInput.value !== '') {
    const val = parseFloat(customInput.value)
    return isNaN(val) || val <= 0 ? 0 : val
  }
  return 0
})

// === Feature icons for featured plan ===
const featureIcons = ['database', 'bolt', 'dollar', 'shield'] as const
const featureIconColors = [
  'text-primary-600 dark:text-primary-400',
  'text-yellow-500',
  'text-green-500',
  'text-blue-500',
]

function getFeatureIcon(index: number): string {
  return featureIcons[index] || 'check'
}

function getFeatureIconColor(index: number): string {
  return featureIconColors[index] || 'text-green-500'
}

// === Hardcoded plan decoration data ===
interface DecoratedPlan {
  plan: PaymentPlan
  badge: string
  audience: string
  highlight: string
  features: string[]
  cta: string
  featured: boolean
}

const planDecorations: Record<string, Omit<DecoratedPlan, 'plan'>> = {
  '29900': {
    badge: '',
    audience: '',
    highlight: '',
    features: [
      '$30/day',
      '$600/month',
      '',
      '',
      '',
    ],
    cta: '',
    featured: false,
  },
  '58900': {
    badge: '',
    audience: '',
    highlight: '',
    features: [
      '$60/day',
      '$1,300/month',
      '',
      '',
      '',
      '',
    ],
    cta: '',
    featured: true,
  },
  '118900': {
    badge: '',
    audience: '',
    highlight: '',
    features: [
      '$120/day',
      '$2,500/month',
      '',
      '',
      '',
    ],
    cta: '',
    featured: false,
  },
}

// Chinese decoration text (hardcoded)
const planDecoZh: Record<string, Omit<DecoratedPlan, 'plan'>> = {
  '29900': {
    badge: '🔥 个人首选 / 火爆热销',
    audience: '[适合：独立开发者 / 个人项目]',
    highlight: '总额$600：官方价格体系，无需复杂换算',
    features: [
      '<strong class="text-gray-900 dark:text-white">每日 $30 高速配额：</strong>专属通道，响应丝滑无卡顿',
      '<strong class="text-gray-900 dark:text-white">总额 $600/月：</strong>充足算力，轻松覆盖日常开发',
      '<strong class="text-gray-900 dark:text-white">永不掉线保障：</strong>额度用尽自动转按量，业务 7x24 在线',
      '<strong class="text-gray-900 dark:text-white">超高性价比：</strong>个人开发首选，一次付费全月无忧',
    ],
    cta: '购买开发者版',
    featured: false,
  },
  '58900': {
    badge: '👑 生产环境推荐 / 主推套餐',
    audience: '[适合：全职开发 / 生产级应用]',
    highlight: '折合 $1=¥0.45，立省 55% 成本！',
    features: [
      '<strong class="text-gray-900 dark:text-white text-base">总额 $1,300：</strong>充足算力储备，应对高频调用',
      '<strong class="text-gray-900 dark:text-white text-base">每日 $60 高速配额：</strong>生产级优先级，拒绝排队等待',
      '<strong class="text-gray-900 dark:text-white text-base">无损计费机制：</strong>拒绝倍率陷阱，每一分实打实可用',
      '<strong class="text-gray-900 dark:text-white text-base">熔断智能防护：</strong>异常流量自动隔离，生产环境稳定如磐石',
    ],
    cta: '立即购买，释放生产力',
    featured: true,
  },
  '118900': {
    badge: '💎 团队专属 / SLA 保障',
    audience: '[适合：开发团队 / AI 产品项目]',
    highlight: '总额$2,500：顶配储备，支撑海量高并发接入',
    features: [
      '<strong class="text-gray-900 dark:text-white">总额 $2,500/月：</strong>顶配算力，满足团队级高频需求',
      '<strong class="text-gray-900 dark:text-white">每日 $120 高速配额：</strong>团队级 VIP 优先级调度，秒级响应',
      '<strong class="text-gray-900 dark:text-white">高可用服务保障：</strong>极致架构，满足企业级稳定运营标准',
      '<strong class="text-gray-900 dark:text-white">无感溢出续航：</strong>额度海量储备，确保 AI 工作流永不断供',
    ],
    cta: '购买旗舰版',
    featured: false,
  },
}

const decoratedPlans = computed<DecoratedPlan[]>(() => {
  const targetAmounts = [29900, 58900, 118900]
  const matched = plans.value.filter(p => targetAmounts.includes(p.amount_fen))
  const source = matched.length > 0 ? matched : plans.value

  return source.map(plan => {
    const key = String(plan.amount_fen)
    const deco = planDecoZh[key] || planDecorations[key]
    if (deco) {
      return { ...deco, plan }
    }
    return {
      plan,
      badge: '',
      audience: '',
      highlight: plan.description || '',
      features: [
        ...(plan.features || []),
        '微信扫码支付',
      ],
      cta: '立即购买',
      featured: false,
    }
  })
})

// === Promo Code Validation ===
function resetPromoValidation() {
  promoValidation.value = {
    valid: false,
    invalid: false,
    discountType: '',
    discountAmount: 0,
    discountFen: 0,
    finalAmountFen: 0,
    minOrderAmount: 0,
    message: ''
  }
}

function handlePromoCodeInput() {
  const code = promoCode.value.trim()
  resetPromoValidation()

  if (!code) {
    promoValidating.value = false
    return
  }
}

async function validatePromoForAmount(code: string, amountFen: number): Promise<boolean> {
  if (!code.trim()) return true

  promoValidating.value = true
  try {
    const result = await validatePromoCode(code, amountFen)
    if (result.valid) {
      promoValidation.value = {
        valid: true,
        invalid: false,
        discountType: result.discount_type || '',
        discountAmount: result.discount_amount || 0,
        discountFen: result.discount_fen || 0,
        finalAmountFen: result.final_amount_fen || amountFen,
        minOrderAmount: result.min_order_amount || 0,
        message: ''
      }
      return true
    } else {
      promoValidation.value = {
        valid: false,
        invalid: true,
        discountType: '',
        discountAmount: 0,
        discountFen: 0,
        finalAmountFen: 0,
        minOrderAmount: result.min_order_amount || 0,
        message: getPromoErrorMessage(result.error_code)
      }
      return false
    }
  } catch {
    promoValidation.value = {
      valid: false,
      invalid: true,
      discountType: '',
      discountAmount: 0,
      discountFen: 0,
      finalAmountFen: 0,
      minOrderAmount: 0,
      message: t('pricing.promoCodeInvalid')
    }
    return false
  } finally {
    promoValidating.value = false
  }
}

function getPromoErrorMessage(errorCode?: string): string {
  switch (errorCode) {
    case 'PROMO_CODE_NOT_FOUND':
      return t('pricing.promoCodeNotFound')
    case 'PROMO_CODE_EXPIRED':
      return t('pricing.promoCodeExpired')
    case 'PROMO_CODE_DISABLED':
      return t('pricing.promoCodeDisabled')
    case 'PROMO_CODE_MAX_USED':
      return t('pricing.promoCodeMaxUsed')
    case 'PROMO_CODE_ALREADY_USED':
      return t('pricing.promoCodeAlreadyUsed')
    case 'PROMO_CODE_MIN_ORDER':
      return t('pricing.promoCodeMinOrder', { amount: (promoValidation.value.minOrderAmount / 100).toFixed(0) })
    default:
      return t('pricing.promoCodeInvalid')
  }
}

function formatDiscountText(): string {
  const v = promoValidation.value
  if (!v.valid) return ''
  const amountYuan = (v.discountFen / 100).toFixed(2)
  if (v.discountType === 'percentage') {
    return t('pricing.promoDiscountPercent', { percent: v.discountAmount, amount: amountYuan })
  }
  return t('pricing.promoDiscountFixed', { amount: amountYuan })
}

onMounted(async () => {
  try {
    const [allPlans, rechargeInfo] = await Promise.all([
      paymentAPI.getPlans(),
      paymentAPI.getRechargeInfo()
    ])
    plans.value = allPlans.filter(p => (p.type || 'subscription') === 'subscription')
    rechargePlans.value = rechargeInfo.plans || []
    rechargeMinAmount.value = rechargeInfo.min_amount || 0
  } catch {
    // silently fail
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  clearTimers()
  if (promoValidateTimeout) {
    clearTimeout(promoValidateTimeout)
  }
})

function clearTimers() {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
  if (countdownTimer) { clearInterval(countdownTimer); countdownTimer = null }
}

// === Payment method selection ===
async function confirmPayMethod() {
  showPayMethodModal.value = false
  if (pendingPaymentAction.value) {
    await pendingPaymentAction.value()
    pendingPaymentAction.value = null
  }
}

// === Subscription purchase ===
async function handleBuy(plan: PaymentPlan) {
  if (creatingOrder.value) return

  const code = promoCode.value.trim()
  if (code) {
    const valid = await validatePromoForAmount(code, plan.amount_fen)
    if (!valid) return
  }

  pendingPaymentAction.value = async () => {
    paymentOrderType.value = 'subscription'
    creatingOrder.value = true
    try {
      const order = await paymentAPI.createOrder(plan.key, code || undefined, selectedPayMethod.value)
      await showQRModal(order)
    } catch (err: any) {
      const msg = err?.message || err?.response?.data?.message || t('pricing.payment.createFailed')
      alert(msg)
    } finally {
      creatingOrder.value = false
    }
  }
  showPayMethodModal.value = true
}

// === Recharge preset ===
async function handleRechargePreset(rp: RechargePlan) {
  if (creatingOrder.value) return

  const code = promoCode.value.trim()
  if (code) {
    const valid = await validatePromoForAmount(code, rp.pay_amount_fen)
    if (!valid) return
  }

  pendingPaymentAction.value = async () => {
    paymentOrderType.value = 'recharge'
    creatingOrder.value = true
    try {
      const payYuan = rp.pay_amount_fen / 100
      const order = await paymentAPI.createRechargeOrder(payYuan, code || undefined, selectedPayMethod.value)
      await showQRModal(order)
    } catch (err: any) {
      const msg = err?.message || err?.response?.data?.message || t('recharge.payment.createFailed')
      alert(msg)
    } finally {
      creatingOrder.value = false
    }
  }
  showPayMethodModal.value = true
}

// === Custom recharge ===
async function handleCustomRecharge() {
  if (customFinalAmount.value <= 0 || creatingOrder.value) return
  if (rechargeMinAmount.value > 0 && customFinalAmount.value < rechargeMinAmount.value) return

  const code = promoCode.value.trim()
  const amountFen = Math.round(customFinalAmount.value * 100)
  if (code) {
    const valid = await validatePromoForAmount(code, amountFen)
    if (!valid) return
  }

  pendingPaymentAction.value = async () => {
    paymentOrderType.value = 'recharge'
    creatingOrder.value = true
    try {
      const order = await paymentAPI.createRechargeOrder(customFinalAmount.value, code || undefined, selectedPayMethod.value)
      await showQRModal(order)
    } catch (err: any) {
      const msg = err?.message || err?.response?.data?.message || t('recharge.payment.createFailed')
      alert(msg)
    } finally {
      creatingOrder.value = false
    }
  }
  showPayMethodModal.value = true
}

// === Shared QR modal logic ===
async function showQRModal(order: { order_no: string; code_url: string | null; amount_fen: number; expired_at: string }) {
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

function goAfterPayment() {
  showPaymentModal.value = false
  clearTimers()
  if (paymentOrderType.value === 'subscription') {
    router.push('/subscriptions')
  } else {
    router.push('/dashboard')
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
  transform: translateY(-4px);
}

.shadow-soft {
  box-shadow: 0 4px 24px -4px rgba(0, 0, 0, 0.08);
}
</style>
