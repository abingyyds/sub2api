<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-6xl space-y-8">
        <!-- Title -->
        <SlideIn direction="up" :delay="100">
          <div class="text-center">
            <h1 class="text-3xl md:text-4xl font-extrabold text-gray-900 dark:text-white tracking-tight">{{ t('pricing.title') }}</h1>
            <p class="mt-3 text-gray-500 dark:text-dark-400 text-lg">{{ t('pricing.subtitle') }}</p>
            <!-- Current plan info -->
            <p v-if="currentPlanInfo" class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ currentPlanInfo }}</p>
          </div>
        </SlideIn>

        <!-- Loading -->
        <div v-if="loading" class="flex justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
        </div>

        <!-- Tab Switcher -->
        <SlideIn v-if="!loading" direction="up" :delay="150">
          <div class="flex justify-center">
            <div class="inline-flex rounded-xl bg-gray-100 dark:bg-dark-800 p-1 gap-1">
              <!-- Newcomer Tab (hidden if not eligible) -->
              <button
                v-if="newcomerEligible"
                class="relative px-6 py-2.5 rounded-lg text-sm font-semibold transition-all"
                :class="activeTab === 'newcomer'
                  ? 'bg-white dark:bg-dark-700 text-gray-900 dark:text-white shadow-sm'
                  : 'text-gray-600 dark:text-dark-400 hover:text-gray-900 dark:hover:text-white'"
                @click="setActiveTab('newcomer')"
              >
                {{ t('pricing.tabs.newcomer') }}
                <span class="absolute -top-1 -right-1 flex h-5 w-5">
                  <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                  <span class="relative inline-flex rounded-full h-5 w-5 bg-green-500 items-center justify-center">
                    <svg class="h-3 w-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                  </span>
                </span>
              </button>
              <!-- Subscription Tab -->
              <button
                class="px-6 py-2.5 rounded-lg text-sm font-semibold transition-all"
                :class="activeTab === 'subscription'
                  ? 'bg-white dark:bg-dark-700 text-gray-900 dark:text-white shadow-sm'
                  : 'text-gray-600 dark:text-dark-400 hover:text-gray-900 dark:hover:text-white'"
                @click="setActiveTab('subscription')"
              >
                {{ t('pricing.tabs.subscription') }}
              </button>
              <!-- Recharge Tab -->
              <button
                class="px-6 py-2.5 rounded-lg text-sm font-semibold transition-all"
                :class="activeTab === 'recharge'
                  ? 'bg-white dark:bg-dark-700 text-gray-900 dark:text-white shadow-sm'
                  : 'text-gray-600 dark:text-dark-400 hover:text-gray-900 dark:hover:text-white'"
                @click="setActiveTab('recharge')"
              >
                {{ t('pricing.tabs.recharge') }}
              </button>
            </div>
          </div>
        </SlideIn>

        <!-- Promo Code Input -->
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

        <template v-if="!loading">
          <!-- ==================== Newcomer Section ==================== -->
          <template v-if="activeTab === 'newcomer'">
            <SlideIn direction="up" :delay="250">
              <div class="mb-6 text-center">
                <h2 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('pricing.newcomerSection') }}</h2>
                <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('pricing.newcomerSectionDesc') }}</p>
              </div>
            </SlideIn>

            <!-- Newcomer Plans -->
            <div v-if="showRechargeCatalogLoader" class="flex justify-center py-12">
              <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
            </div>
            <StaggerContainer v-else-if="newcomerPlans.length > 0" :stagger-delay="100" :delay="300">
              <div class="grid gap-8 md:grid-cols-2 max-w-4xl mx-auto">
                <GlowCard
                  v-for="rp in newcomerPlans"
                  :key="rp.key"
                  :glow-color="rp.pay_amount_fen >= 5000 ? 'rgb(99, 102, 241)' : 'rgb(34, 197, 94)'"
                >
                  <div
                    class="relative flex flex-col rounded-2xl border-2 bg-white dark:bg-dark-900 transition-all hover:shadow-xl hover:-translate-y-1 cursor-pointer h-full"
                    :class="rp.pay_amount_fen >= 5000
                      ? 'border-indigo-400 dark:border-indigo-500/60 shadow-lg shadow-indigo-500/10'
                      : 'border-green-400 dark:border-green-500/60 shadow-lg shadow-green-500/10'"
                    :style="orderFlowBusy ? 'pointer-events:none;' : ''"
                    @click="handleRechargePreset(rp)"
                  >
                    <!-- Badge -->
                    <div class="absolute -top-3 left-1/2 -translate-x-1/2">
                      <span class="text-white text-xs font-bold px-4 py-1 rounded-full shadow-md whitespace-nowrap"
                        :class="rp.pay_amount_fen >= 5000
                          ? 'bg-gradient-to-r from-indigo-500 to-purple-500'
                          : 'bg-gradient-to-r from-green-500 to-emerald-500'"
                      >
                        {{ rp.pay_amount_fen >= 5000 ? '🚀 技术指导 · 手把手带飞' : '🎉 上手即用 · 零门槛' }}
                      </span>
                    </div>

                    <div class="flex flex-1 flex-col p-8">
                      <!-- Title -->
                      <h3 class="text-sm font-medium text-gray-500 dark:text-dark-400 mt-2">{{ rp.name || t('pricing.purchaseQuota') }}</h3>
                      <!-- Balance amount (USD) -->
                      <div class="mt-2 flex items-baseline gap-1">
                        <span class="text-4xl font-extrabold text-gray-900 dark:text-white">${{ rp.balance_amount.toFixed(2) }}</span>
                        <span class="text-sm text-gray-500 dark:text-dark-400">USD</span>
                      </div>
                      <!-- Special price label -->
                      <p v-if="rp.pay_amount_fen !== rp.balance_amount * 100" class="mt-1 text-xs font-semibold"
                        :class="rp.pay_amount_fen >= 5000 ? 'text-indigo-600 dark:text-indigo-400' : 'text-green-600 dark:text-green-400'"
                      >
                        {{ t('pricing.newcomerSpecialPrice') }}
                      </p>

                      <!-- Equivalent value -->
                      <div class="mt-2 text-xs font-bold text-green-600 dark:text-green-400">
                        ≈ 官方 ${{ (rp.balance_amount * 2).toFixed(2) }} 等效算力
                      </div>

                      <!-- Feature highlights -->
                      <ul class="mt-5 space-y-3 text-sm text-gray-600 dark:text-dark-400 flex-grow">
                        <template v-if="rp.pay_amount_fen >= 5000">
                          <li class="flex items-start gap-2.5">
                            <Icon name="check" size="sm" class="text-indigo-500 mt-0.5 flex-shrink-0" />
                            <span><strong class="text-gray-900 dark:text-white">1 对 1 远程技术指导：</strong>专人手把手帮你配置 Claude Code 开发环境</span>
                          </li>
                          <li class="flex items-start gap-2.5">
                            <Icon name="check" size="sm" class="text-indigo-500 mt-0.5 flex-shrink-0" />
                            <span><strong class="text-gray-900 dark:text-white">环境搭建全包：</strong>API 密钥配置、IDE 插件安装、代理设置一步到位</span>
                          </li>
                          <li class="flex items-start gap-2.5">
                            <Icon name="check" size="sm" class="text-indigo-500 mt-0.5 flex-shrink-0" />
                            <span><strong class="text-gray-900 dark:text-white">${{ rp.balance_amount.toFixed(2) }} 充足额度：</strong>≈ 官方 ${{ (rp.balance_amount * 2).toFixed(0) }} 等效算力，足够深度体验</span>
                          </li>
                          <li class="flex items-start gap-2.5">
                            <Icon name="check" size="sm" class="text-indigo-500 mt-0.5 flex-shrink-0" />
                            <span><strong class="text-gray-900 dark:text-white">新手零压力：</strong>不会配置？不会用？全程带你搞定，确保跑通</span>
                          </li>
                        </template>
                        <template v-else>
                          <li class="flex items-start gap-2.5">
                            <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                            <span><strong class="text-gray-900 dark:text-white">零门槛体验：</strong>注册即买，立刻接入 Claude / Codex / Gemini</span>
                          </li>
                          <li class="flex items-start gap-2.5">
                            <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                            <span><strong class="text-gray-900 dark:text-white">超值性价比：</strong>¥{{ (rp.pay_amount_fen / 100).toFixed(rp.pay_amount_fen % 100 === 0 ? 0 : 2) }} 即享 ${{ rp.balance_amount.toFixed(2) }} 额度，先用再说</span>
                          </li>
                          <li class="flex items-start gap-2.5">
                            <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                            <span><strong class="text-gray-900 dark:text-white">即刻生效：</strong>支付后秒级到账，无需等待审核</span>
                          </li>
                        </template>
                      </ul>

                      <div class="mt-auto pt-5">
                        <!-- Price comparison -->
                        <div class="mb-3 text-center">
                          <div class="text-xs text-gray-500 dark:text-dark-400 line-through">
                            原价 ¥{{ (rp.balance_amount * 7.2).toFixed(0) }}
                          </div>
                          <div class="text-2xl font-extrabold text-gray-900 dark:text-white">
                            ¥{{ (rp.pay_amount_fen / 100).toFixed(rp.pay_amount_fen % 100 === 0 ? 0 : 2) }}
                          </div>
                          <div class="text-xs font-bold mt-1"
                            :class="rp.pay_amount_fen >= 5000 ? 'text-indigo-600 dark:text-indigo-400' : 'text-green-600 dark:text-green-400'"
                          >
                            立省 ¥{{ ((rp.balance_amount * 7.2) - (rp.pay_amount_fen / 100)).toFixed(0) }}，仅需 {{ ((rp.pay_amount_fen / 100) / (rp.balance_amount * 7.2) * 100).toFixed(0) }}% 价格！
                          </div>
                        </div>
                        <!-- CTA button -->
                        <MagneticButton>
                          <div
                            class="flex items-center justify-center rounded-xl py-3.5 text-base font-bold text-white transition-all shadow-lg"
                            :class="rp.pay_amount_fen >= 5000
                              ? 'bg-gradient-to-r from-indigo-500 to-purple-500 hover:from-indigo-600 hover:to-purple-600 shadow-indigo-500/30'
                              : 'bg-gradient-to-r from-green-500 to-emerald-500 hover:from-green-600 hover:to-emerald-600 shadow-green-500/30'"
                          >
                            <template v-if="isProcessingOrderAction('recharge', rp.key)">
                              <div class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
                              {{ t('common.processing') }}
                            </template>
                            <template v-else>
                              {{ rp.pay_amount_fen >= 5000 ? '立即抢购 + 赠技术指导' : '立即抢购' }}
                            </template>
                          </div>
                        </MagneticButton>
                        <!-- Purchase limit note -->
                        <p v-if="rp.max_purchases && rp.max_purchases > 0" class="mt-2 text-center text-xs text-gray-400 dark:text-dark-500">
                          {{ t('pricing.purchaseLimitNote', { count: rp.max_purchases }) }}
                        </p>
                      </div>
                    </div>
                  </div>
                </GlowCard>
              </div>
            </StaggerContainer>
            <div v-else class="text-center py-12 text-gray-500 dark:text-dark-400">
              {{ t('pricing.noNewcomerPlans') }}
            </div>
          </template>

          <!-- ==================== Subscription Section (套餐订阅) ==================== -->
          <template v-if="activeTab === 'subscription'">
          <SlideIn direction="up" :delay="250">
            <div class="mb-6 text-center">
              <h2 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('pricing.subscriptionSection') }}</h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('pricing.subscriptionSectionDesc') }}</p>
              <div class="mt-4 inline-block bg-gradient-to-r from-orange-500 to-red-500 text-white text-sm font-bold px-6 py-2.5 rounded-full shadow-lg">
                🔥 特价 Claude Code API：¥0.28 = 官方 $1，仅官方价格 3.9%！
              </div>
              <p class="text-xs text-gray-400 dark:text-dark-500 mt-2">套餐 $1 = 官方 $2.5 算力 × 80% 性能 = 官方 $2 等效算力</p>
            </div>
          </SlideIn>

          <StaggerContainer v-if="decoratedPlans.length > 0" :stagger-delay="120" :delay="300">
            <div class="grid gap-8 lg:grid-cols-3 items-stretch">
              <GlowCard
                v-for="dp in decoratedPlans"
                :key="dp.plan.key"
                :glow-color="dp.featured ? 'rgb(217, 119, 87)' : dp.plan.amount_fen >= 100000 ? 'rgb(168, 85, 247)' : 'rgb(234, 179, 8)'"
              >
                <div
                  class="relative flex flex-col rounded-3xl border-2 bg-white dark:bg-dark-900 transition-all hover:shadow-lg cursor-pointer h-full"
                  :class="dp.featured
                    ? 'border-[3px] border-primary-500/80 dark:border-primary-400/60 lg:scale-105 shadow-2xl shadow-primary-500/20 z-10'
                    : dp.plan.amount_fen >= 100000
                      ? 'border-purple-100 dark:border-purple-800/30 shadow-soft'
                      : 'border-orange-100 dark:border-orange-800/30 shadow-soft'"
                  :style="orderFlowBusy ? 'pointer-events:none;' : ''"
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
                    <div class="mb-2"
                      :class="dp.featured ? 'text-5xl' : 'text-4xl'"
                    >
                      <span class="font-extrabold text-gray-900 dark:text-white">¥{{ (dp.plan.amount_fen / 100).toFixed(0) }}</span>
                    </div>

                    <!-- Quota info with equivalent -->
                    <div class="text-sm text-gray-500 dark:text-dark-400 mb-1">
                      {{ dp.highlight }}
                    </div>
                    <div class="text-xs font-bold text-green-600 dark:text-green-400 mb-6">
                      ≈ 官方 ${{ getEquivalentQuota(dp.plan.amount_fen) }} 等效算力
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
                        :disabled="orderFlowBusy"
                      >
                        <template v-if="isProcessingOrderAction('subscription', dp.plan.key)">
                          <span class="inline-flex items-center justify-center">
                            <span class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"></span>
                            {{ t('common.processing') }}
                          </span>
                        </template>
                        <template v-else>
                          {{ dp.cta }}
                        </template>
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
            {{ t('pricing.repurchaseRule') }}
          </p>

          </template>

          <!-- ==================== Recharge Section (灵活额度) ==================== -->
          <template v-if="activeTab === 'recharge'">
          <SlideIn direction="up" :delay="400">
            <div class="mb-2 mt-4">
              <h2 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('pricing.balanceSection') }}</h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('pricing.balanceSectionDesc') }}</p>
            </div>
          </SlideIn>

          <!-- Recharge Plans -->
          <div v-if="showRechargeCatalogLoader" class="flex justify-center py-12">
            <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
          </div>
          <StaggerContainer v-else-if="regularRechargePlans.length > 0" :stagger-delay="100" :delay="450">
            <div class="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
              <GlowCard
                v-for="rp in regularRechargePlans"
                :key="rp.key"
                :glow-color="rp.popular ? 'rgb(217, 119, 87)' : 'rgb(59, 130, 246)'"
              >
                <div
                  class="relative flex flex-col rounded-2xl border-2 bg-white dark:bg-dark-900 transition-all hover:shadow-lg hover:-translate-y-1 cursor-pointer h-full"
                  :class="rp.popular
                    ? 'border-primary-500/80 dark:border-primary-400/60 shadow-lg shadow-primary-500/10'
                    : 'border-gray-200 dark:border-dark-600 shadow-soft'"
                  :style="orderFlowBusy ? 'pointer-events:none;' : ''"
                  @click="handleRechargePreset(rp)"
                >
                  <!-- Popular badge -->
                  <div v-if="rp.popular" class="absolute -top-3 left-1/2 -translate-x-1/2">
                    <span class="bg-gradient-to-r from-orange-500 to-red-500 text-white text-xs font-bold px-4 py-1 rounded-full shadow-md whitespace-nowrap">
                      {{ t('pricing.mostPopular') }}
                    </span>
                  </div>

                  <div class="flex flex-1 flex-col p-5">
                    <!-- Title -->
                    <h3 class="text-sm font-medium text-gray-500 dark:text-dark-400">{{ rp.name || t('pricing.purchaseQuota') }}</h3>
                    <!-- Balance amount (USD) -->
                    <div class="mt-2 flex items-baseline gap-1">
                      <span class="text-2xl font-extrabold text-gray-900 dark:text-white">${{ rp.balance_amount.toFixed(2) }}</span>
                      <span class="text-sm text-gray-500 dark:text-dark-400">USD</span>
                    </div>
                    <!-- Description -->
                    <p v-if="rp.description" class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ rp.description }}</p>

                    <div class="mt-auto pt-4">
                      <!-- CTA button -->
                      <div
                        class="flex items-center justify-center rounded-xl py-2.5 px-3 text-xs sm:text-sm font-bold text-white transition-all whitespace-nowrap overflow-hidden truncate"
                        :class="rp.popular
                          ? 'bg-primary-600 hover:bg-primary-700 shadow-md'
                          : 'bg-gray-900 hover:bg-gray-800 dark:bg-dark-500 dark:hover:bg-dark-400'"
                      >
                        <template v-if="isProcessingOrderAction('recharge', rp.key)">
                          <div class="mr-1.5 h-3.5 w-3.5 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
                          {{ t('common.processing') }}
                        </template>
                        <template v-else>
                          {{ t('pricing.rechargePayBtn') }} ¥{{ (rp.pay_amount_fen / 100).toFixed(rp.pay_amount_fen % 100 === 0 ? 0 : 2) }}
                        </template>
                      </div>
                      <!-- Purchase limit note -->
                      <p v-if="rp.max_purchases && rp.max_purchases > 0" class="mt-2 text-center text-xs text-gray-400 dark:text-dark-500">
                        {{ t('pricing.purchaseLimitNote', { count: rp.max_purchases }) }}
                      </p>
                    </div>
                  </div>
                </div>
              </GlowCard>

              <!-- Custom Amount Card - same grid -->
              <GlowCard glow-color="rgb(99, 102, 241)">
                <div
                  class="relative flex flex-col rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-600 dark:bg-dark-900 shadow-soft h-full"
                >
                  <div class="flex flex-1 flex-col p-5">
                    <h3 class="text-sm font-medium text-gray-500 dark:text-dark-400">{{ t('pricing.customRecharge') }}</h3>

                    <div class="mt-3">
                      <div class="relative">
                        <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 dark:text-dark-400 font-bold text-sm">¥</span>
                        <input
                          v-model="customInput"
                          type="number"
                          min="1"
                          step="1"
                          :placeholder="rechargeMinAmount > 0 ? t('pricing.minAmountHint', { amount: rechargeMinAmount }) : t('recharge.inputPlaceholder')"
                          class="w-full rounded-lg border border-gray-200 bg-gray-50 py-2.5 pl-8 pr-3 text-sm font-bold text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:bg-white focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-800 dark:text-white dark:placeholder-dark-500 dark:focus:bg-dark-700 transition-all"
                          @click.stop
                        />
                      </div>
                      <p v-if="rechargeMinAmount > 0" class="mt-1 text-xs text-gray-400 dark:text-dark-500">
                        {{ t('pricing.minAmountNote', { amount: rechargeMinAmount }) }}
                      </p>
                    </div>

                    <!-- Amount Summary -->
                    <div v-if="customFinalAmount > 0" class="mt-3 rounded-lg bg-primary-50 dark:bg-primary-900/10 border border-primary-100 dark:border-primary-800/30 p-2.5">
                      <div class="flex items-center justify-between">
                        <span class="text-xs text-gray-600 dark:text-dark-400">{{ t('recharge.rechargeAmount') }}</span>
                        <span class="text-sm font-extrabold text-primary-600">¥{{ customFinalAmount.toFixed(2) }}</span>
                      </div>
                    </div>

                    <div class="mt-auto pt-4">
                      <div
                        class="flex items-center justify-center rounded-xl py-2.5 px-3 text-xs sm:text-sm font-bold text-white transition-all whitespace-nowrap overflow-hidden truncate bg-primary-600 hover:bg-primary-700 shadow-md cursor-pointer"
                        :class="{ 'opacity-50 cursor-not-allowed': customFinalAmount <= 0 || (rechargeMinAmount > 0 && customFinalAmount < rechargeMinAmount) || orderFlowBusy }"
                        @click="handleCustomRecharge"
                      >
                        <template v-if="isProcessingOrderAction('custom')">
                          <div class="h-3.5 w-3.5 animate-spin rounded-full border-2 border-white border-t-transparent mr-1.5"></div>
                          {{ t('recharge.creating') }}
                        </template>
                        <template v-else>
                          {{ customFinalAmount > 0 ? t('recharge.payNow', { amount: customFinalAmount.toFixed(2) }) : t('recharge.enterAmount') }}
                        </template>
                      </div>
                    </div>
                  </div>
                </div>
              </GlowCard>
            </div>
          </StaggerContainer>

          <!-- ==================== FAQ Section ==================== -->
          </template>

          <!-- FAQ Section (shown on all tabs) -->
          <SlideIn direction="up" :delay="500">
            <div class="mt-8">
              <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">{{ t('pricing.faq.title') }}</h2>
              <div class="space-y-4">
                <div
                  v-for="(item, idx) in faqItems"
                  :key="idx"
                  class="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-800"
                >
                  <div class="border-b border-gray-100 bg-gray-50/90 px-5 py-4 dark:border-dark-700 dark:bg-dark-900/80">
                    <div class="flex items-start gap-3">
                      <span class="inline-flex h-7 min-w-7 items-center justify-center rounded-full bg-primary-600 px-2 text-xs font-bold text-white">
                        Q{{ idx + 1 }}
                      </span>
                      <h3 class="pt-0.5 text-base font-semibold leading-7 text-gray-900 dark:text-white">
                        {{ item.q }}
                      </h3>
                    </div>
                  </div>
                  <div class="px-5 py-5">
                    <div class="space-y-3 text-sm leading-7 text-gray-700 dark:text-dark-200" v-html="formatFaqAnswer(item.a)"></div>
                  </div>
                </div>
              </div>
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
          @click.self="cancelPayMethodSelection"
        >
          <div class="w-full max-w-md rounded-2xl bg-white p-8 shadow-2xl dark:bg-dark-800 mx-4 border border-gray-100 dark:border-dark-700">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-6 text-center">选择支付方式</h3>

            <div class="space-y-3 mb-8">
              <div
                v-for="method in availablePayMethods"
                :key="method"
                class="flex items-center gap-4 rounded-xl border-2 p-4 cursor-pointer transition-all"
                :class="selectedPayMethod === method
                  ? (payMethodMeta(method).color === 'green' ? 'border-green-500 bg-green-50 dark:bg-green-900/20 shadow-sm' : 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 shadow-sm')
                  : (payMethodMeta(method).color === 'green' ? 'border-gray-200 dark:border-dark-600 hover:border-green-300' : 'border-gray-200 dark:border-dark-600 hover:border-blue-300')"
                @click="selectedPayMethod = method"
              >
                <input type="radio" :checked="selectedPayMethod === method" class="h-5 w-5" :class="payMethodMeta(method).color === 'green' ? 'text-green-600' : 'text-blue-600'" />
                <svg class="h-8 w-8" viewBox="0 0 24 24" fill="none" v-html="payMethodMeta(method).icon"></svg>
                <span class="font-bold text-gray-900 dark:text-white">{{ payMethodMeta(method).label }}</span>
              </div>
            </div>

            <div class="flex gap-3">
              <button
                class="flex-1 rounded-xl border-2 border-gray-200 dark:border-dark-600 py-3 font-bold text-gray-700 dark:text-dark-300 hover:bg-gray-50 dark:hover:bg-dark-700 transition"
                @click="cancelPayMethodSelection"
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
              <div class="relative rounded-xl border-2 border-gray-100 dark:border-dark-700 p-2">
                <canvas ref="qrCanvas" class="rounded-lg"></canvas>
                <div
                  v-if="qrLoading"
                  class="absolute inset-0 flex items-center justify-center rounded-xl bg-white/90 dark:bg-dark-800/90"
                >
                  <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
                </div>
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
                  {{ t('pricing.payment.viewOrderHistory') }}
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
import type { PaymentPlan, RechargePlan, PayMethod } from '@/api/payment'
import { FadeIn, SlideIn, StaggerContainer, GlowCard, MagneticButton } from '@/components/animations'
import { useAppStore, useAuthStore, useSubscriptionStore } from '@/stores'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()
const subscriptionStore = useSubscriptionStore()

const loading = ref(true)
const plans = ref<PaymentPlan[]>([])
const rechargePlans = ref<RechargePlan[]>([])
const rechargeMinAmount = ref(0)
const rechargeCatalogLoading = ref(false)
const rechargeCatalogLoaded = ref(false)
const creatingOrder = ref(false)
const preparingOrder = ref(false)
const availablePayMethods = ref<PayMethod[]>([])
const payMethodsLoaded = ref(false)
const selectedPayMethod = ref<PayMethod>('wechat')
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
const paymentSuccessHandled = ref(false)
const customInput = ref('')
const activeOrderActionKey = ref('')

// Tab state
const activeTab = ref<'newcomer' | 'subscription' | 'recharge'>('subscription')
const newcomerEligible = ref(false)
const hasTabInteraction = ref(false)

// Computed: filter newcomer plans
const newcomerPlans = computed(() => rechargePlans.value.filter(p => p.is_newcomer))
const regularRechargePlans = computed(() => rechargePlans.value.filter(p => !p.is_newcomer))
const showRechargeCatalogLoader = computed(() => (
  rechargeCatalogLoading.value
  && activeTab.value !== 'subscription'
  && rechargePlans.value.length === 0
))

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
let rechargeCatalogPromise: Promise<void> | null = null
let payMethodsPromise: Promise<void> | null = null

let pollTimer: ReturnType<typeof setInterval> | null = null
let countdownTimer: ReturnType<typeof setInterval> | null = null
const orderFlowBusy = computed(() => preparingOrder.value || creatingOrder.value)

function makeOrderActionKey(type: 'subscription' | 'recharge' | 'custom', key?: string) {
  return `${type}:${key || 'default'}`
}

function isProcessingOrderAction(type: 'subscription' | 'recharge' | 'custom', key?: string) {
  return orderFlowBusy.value && activeOrderActionKey.value === makeOrderActionKey(type, key)
}

function getCreateOrderErrorMessage(err: any) {
  if (err?.reason === 'SUBSCRIPTION_REPURCHASE_BLOCKED') {
    return t('pricing.payment.subscriptionRepurchaseBlocked')
  }
  return err?.message || err?.response?.data?.message || t('pricing.payment.createFailed')
}

const customFinalAmount = computed(() => {
  if (customInput.value !== '') {
    const val = parseFloat(customInput.value)
    return isNaN(val) || val <= 0 ? 0 : val
  }
  return 0
})

// Current plan info line
const currentPlanInfo = ref('')

const faqItems = computed(() => [
  { q: t('pricing.faq.q1'), a: t('pricing.faq.a1') },
  { q: t('pricing.faq.q2'), a: t('pricing.faq.a2') },
  { q: t('pricing.faq.q3'), a: t('pricing.faq.a3') },
  { q: t('pricing.faq.q4'), a: t('pricing.faq.a4') },
  { q: t('pricing.faq.q5'), a: t('pricing.faq.a5') },
  { q: t('pricing.faq.q6'), a: t('pricing.faq.a6') },
])

const faqHighlightPatterns = [
  /每30天自动重置/g,
  /30天自动重置/g,
  /每30天/g,
  /订阅额度/g,
  /灵活额度/g,
  /官方 Token 价格/g,
  /计费标准/g,
  /永久有效/g,
  /优先消耗订阅额度/g,
  /套餐到期后失效/g,
  /自动重置/g,
]

function escapeHtml(value: string) {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function highlightFaqLine(line: string) {
  let highlighted = line

  for (const pattern of faqHighlightPatterns) {
    highlighted = highlighted.replace(
      pattern,
      '<span class="rounded-md bg-amber-100 px-1.5 py-0.5 font-semibold text-amber-900 dark:bg-amber-500/20 dark:text-amber-100">$&</span>'
    )
  }

  highlighted = highlighted.replace(
    /^(订阅额度：|灵活额度：|Subscription quota:|Flexible balance:)/,
    '<span class="font-semibold text-gray-900 dark:text-white">$1</span>'
  )

  return highlighted
}

function formatFaqAnswer(answer: string) {
  return answer
    .split('\n')
    .map(line => line.trim())
    .filter(Boolean)
    .map((line, index) => {
      const escapedLine = escapeHtml(line)
      const formattedLine = highlightFaqLine(escapedLine)
      const lineClass = index === 0
        ? 'rounded-xl bg-primary-50/80 px-4 py-3 text-gray-800 dark:bg-primary-900/20 dark:text-dark-100'
        : 'px-1 text-gray-700 dark:text-dark-200'
      return `<p class="${lineClass}">${formattedLine}</p>`
    })
    .join('')
}

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
      'Monthly quota $360',
      'WeChat Pay',
    ],
    cta: '',
    featured: false,
  },
  '49900': {
    badge: '',
    audience: '',
    highlight: '',
    features: [
      'Monthly quota $650',
      'WeChat Pay',
    ],
    cta: '',
    featured: true,
  },
  '99900': {
    badge: '',
    audience: '',
    highlight: '',
    features: [
      'Monthly quota $1400',
      'WeChat Pay',
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
    highlight: '每月额度 $360',
    features: [
      '<strong class="text-gray-900 dark:text-white">每周 $90 高速配额：</strong>专属通道，响应丝滑无卡顿',
      '<strong class="text-gray-900 dark:text-white">总额 $360/月：</strong>充足算力，轻松覆盖日常开发',
      '<strong class="text-gray-900 dark:text-white">永不掉线保障：</strong>额度用尽自动转按量，业务 7x24 在线',
      '<strong class="text-gray-900 dark:text-white">超高性价比：</strong>个人开发首选，一次付费全月无忧',
    ],
    cta: '购买开发者版',
    featured: false,
  },
  '49900': {
    badge: '👑 生产环境推荐 / 主推套餐',
    audience: '[适合：全职开发 / 生产级应用]',
    highlight: '每月额度 $650',
    features: [
      '<strong class="text-gray-900 dark:text-white text-base">总额 $650：</strong>充足算力储备，应对高频调用',
      '<strong class="text-gray-900 dark:text-white text-base">每周 $162.5 高速配额：</strong>生产级优先级，拒绝排队等待',
      '<strong class="text-gray-900 dark:text-white text-base">透明计费机制：</strong>按量清晰结算，每一分实打实可用',
      '<strong class="text-gray-900 dark:text-white text-base">熔断智能防护：</strong>异常流量自动隔离，生产环境稳定如磐石',
    ],
    cta: '立即购买，释放生产力',
    featured: true,
  },
  '99900': {
    badge: '💎 团队专属 / SLA 保障',
    audience: '[适合：开发团队 / AI 产品项目]',
    highlight: '每月额度 $1,400',
    features: [
      '<strong class="text-gray-900 dark:text-white">总额 $1,400/月：</strong>顶配算力，满足团队级高频需求',
      '<strong class="text-gray-900 dark:text-white">每周 $350 高速配额：</strong>团队级 VIP 优先级调度，秒级响应',
      '<strong class="text-gray-900 dark:text-white">高可用服务保障：</strong>极致架构，满足企业级稳定运营标准',
      '<strong class="text-gray-900 dark:text-white">无感溢出续航：</strong>额度海量储备，确保 AI 工作流永不断供',
    ],
    cta: '购买旗舰版',
    featured: false,
  },
}

function hideRateFeatures(plan: PaymentPlan): PaymentPlan {
  return {
    ...plan,
    features: (plan.features || []).filter(feature => !/(倍率|费率|\d+(?:\.\d+)?\s*x\b)/i.test(feature))
  }
}

const decoratedPlans = computed<DecoratedPlan[]>(() => {
  const targetAmounts = [29900, 49900, 99900]
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

// Calculate equivalent official quota
function getEquivalentQuota(amountFen: number): string {
  // Map amount to monthly quota
  const quotaMap: Record<number, number> = {
    29900: 360,
    49900: 650,
    99900: 1400,
  }
  const quota = quotaMap[amountFen] || 0
  // $1 套餐 = 官方 $2.5 × 80% = 官方 $2 等效
  const equivalent = quota * 2
  return equivalent.toLocaleString()
}

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

function applyRechargeCatalog(
  rechargeInfo: Awaited<ReturnType<typeof paymentAPI.getRechargeInfo>>,
  eligible: boolean
) {
  rechargePlans.value = rechargeInfo.plans || []
  rechargeMinAmount.value = rechargeInfo.min_amount || 0
  newcomerEligible.value = eligible && rechargePlans.value.some((plan) => plan.is_newcomer)
  rechargeCatalogLoaded.value = true

  if (!hasTabInteraction.value && activeTab.value === 'subscription' && plans.value.length === 0) {
    activeTab.value = newcomerEligible.value ? 'newcomer' : 'recharge'
  }
}

async function ensureRechargeCatalogLoaded(): Promise<void> {
  if (rechargeCatalogLoaded.value) {
    return
  }

  if (rechargeCatalogPromise) {
    return rechargeCatalogPromise
  }

  rechargeCatalogLoading.value = true

  const request = Promise.allSettled([
    paymentAPI.getRechargeInfo(),
    paymentAPI.getNewcomerStatus()
  ])
    .then((results) => {
      const [rechargeInfoResult, newcomerStatusResult] = results
      if (rechargeInfoResult.status !== 'fulfilled') {
        return
      }

      const eligible = newcomerStatusResult.status === 'fulfilled'
        ? newcomerStatusResult.value.eligible
        : false

      applyRechargeCatalog(rechargeInfoResult.value, eligible)
    })
    .finally(() => {
      rechargeCatalogLoading.value = false
      rechargeCatalogPromise = null
    })

  rechargeCatalogPromise = request

  return request
}

function applyPayMethods(methods: PayMethod[]) {
  availablePayMethods.value = methods || []
  payMethodsLoaded.value = true

  if (methods.length > 0) {
    selectedPayMethod.value = methods[0]
  }
}

async function ensurePayMethodsLoaded(): Promise<void> {
  if (payMethodsLoaded.value) {
    return
  }

  if (payMethodsPromise) {
    return payMethodsPromise
  }

  const request = paymentAPI.getPayMethods()
    .then((methods) => {
      applyPayMethods(methods)
    })
    .catch(() => {})
    .finally(() => {
      payMethodsPromise = null
    })

  payMethodsPromise = request

  return request
}

function scheduleSecondaryPricingDataLoad() {
  const load = () => {
    void ensureRechargeCatalogLoaded()
    void ensurePayMethodsLoaded()
  }

  if (typeof window.requestIdleCallback === 'function') {
    window.requestIdleCallback(load, { timeout: 1500 })
    return
  }

  window.setTimeout(load, 300)
}

function cancelPayMethodSelection() {
  showPayMethodModal.value = false
  pendingPaymentAction.value = null
  preparingOrder.value = false
  activeOrderActionKey.value = ''
}

function setActiveTab(tab: 'newcomer' | 'subscription' | 'recharge') {
  hasTabInteraction.value = true
  activeTab.value = tab

  if (tab !== 'subscription') {
    void ensureRechargeCatalogLoaded()
  }
}

onMounted(async () => {
  window.addEventListener('focus', handlePaymentReturn)
  document.addEventListener('visibilitychange', handlePaymentVisibilityChange)

  const requestedTab = route.query.tab
  if (requestedTab === 'newcomer' || requestedTab === 'subscription' || requestedTab === 'recharge') {
    activeTab.value = requestedTab
  }

  try {
    const allPlans = await paymentAPI.getPlans()
    plans.value = allPlans
      .filter(p => (p.type || 'subscription') === 'subscription')
      .map(hideRateFeatures)
  } catch {
    // silently fail
  } finally {
    loading.value = false
    scheduleSecondaryPricingDataLoad()
  }
})

onUnmounted(() => {
  window.removeEventListener('focus', handlePaymentReturn)
  document.removeEventListener('visibilitychange', handlePaymentVisibilityChange)
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
// === Payment method helpers ===
const WECHAT_ICON = '<rect width="24" height="24" rx="4" fill="#07C160"/><path d="M16.7 10.5c-.2 0-.4 0-.6.1.1-.4.1-.8.1-1.2 0-2.8-2.7-5-5.9-5C7 4.4 4.4 6.6 4.4 9.4c0 1.5.8 2.9 2.1 3.9l-.5 1.6 1.9-1c.6.2 1.2.3 1.8.3h.2c-.1-.3-.1-.7-.1-1 0-2.4 2.3-4.4 5.1-4.4l.2-.1c-.1-.1-.2-.2-.4-.2zm-3.3-2.2c.4 0 .7.3.7.7s-.3.7-.7.7-.7-.3-.7-.7.3-.7.7-.7zm-4.2 1.4c-.4 0-.7-.3-.7-.7s.3-.7.7-.7.7.3.7.7-.3.7-.7.7zm10.4 4.5c0-2.2-2.2-4-4.8-4s-4.8 1.8-4.8 4 2.2 4 4.8 4c.5 0 1-.1 1.5-.2l1.5.8-.4-1.3c1.2-.8 2.2-2 2.2-3.3zm-6.4-.6c-.3 0-.6-.3-.6-.6s.3-.6.6-.6.6.3.6.6-.3.6-.6.6zm3.2 0c-.3 0-.6-.3-.6-.6s.3-.6.6-.6.6.3.6.6-.3.6-.6.6z" fill="white"/>'
const ALIPAY_ICON = '<rect width="24" height="24" rx="4" fill="#1677FF"/><path d="M17.5 14.2c-1.2-.5-2.3-1-3.2-1.4.4-.8.7-1.7.9-2.6h-2.5v-1h3V8.4h-3V6.5h-1.4v1.9h-3v.8h3v1h-2.5v.8h4.6c-.2.7-.5 1.3-.8 1.9-1.3-.5-2.7-.8-3.8-.8-2 0-3.3.9-3.3 2.3 0 1.5 1.3 2.3 3.2 2.3 1.5 0 2.9-.6 4-1.5.9.5 1.9 1 3 1.5l.8-1.2zM9.2 16.3c-1.3 0-1.9-.5-1.9-1.2 0-.8.7-1.3 1.9-1.3.9 0 1.9.2 2.9.7-.9.8-1.9 1.8-2.9 1.8z" fill="white"/>'

const PAY_METHOD_META: Record<string, { label: string; icon: string; color: 'green' | 'blue' }> = {
  wechat: { label: '微信支付', icon: WECHAT_ICON, color: 'green' },
  alipay: { label: '支付宝', icon: ALIPAY_ICON, color: 'blue' },
  epay_alipay: { label: 'Epay-支付宝', icon: ALIPAY_ICON, color: 'blue' },
  epay_wxpay: { label: 'Epay-微信', icon: WECHAT_ICON, color: 'green' },
}

function payMethodMeta(method: string) {
  return PAY_METHOD_META[method] || { label: method, icon: '', color: 'blue' as const }
}

async function showPayMethodOrDirect(action: () => Promise<void>) {
  preparingOrder.value = true
  pendingPaymentAction.value = action
  await ensurePayMethodsLoaded()

  if (availablePayMethods.value.length === 0) {
    pendingPaymentAction.value = null
    preparingOrder.value = false
    alert(t('pricing.payment.noMethodsAvailable'))
    return
  }

  if (availablePayMethods.value.length <= 1) {
    // Only one method (or none) — skip the modal, execute directly
    await confirmPayMethod()
  } else {
    preparingOrder.value = false
    showPayMethodModal.value = true
  }
}

async function confirmPayMethod() {
  showPayMethodModal.value = false
  if (pendingPaymentAction.value) {
    try {
      preparingOrder.value = true
      await pendingPaymentAction.value()
    } finally {
      pendingPaymentAction.value = null
      preparingOrder.value = false
    }
  } else {
    preparingOrder.value = false
  }
}

// === Subscription purchase ===
async function handleBuy(plan: PaymentPlan) {
  if (orderFlowBusy.value) return
  activeOrderActionKey.value = makeOrderActionKey('subscription', plan.key)

  const code = promoCode.value.trim()
  if (code) {
    preparingOrder.value = true
    const valid = await validatePromoForAmount(code, plan.amount_fen)
    if (!valid) {
      preparingOrder.value = false
      activeOrderActionKey.value = ''
      return
    }
  }

  try {
    await showPayMethodOrDirect(async () => {
      paymentOrderType.value = 'subscription'
      creatingOrder.value = true
      try {
        const order = await paymentAPI.createOrder(plan.key, code || undefined, selectedPayMethod.value)
        await showQRModal(order)
      } catch (err: any) {
        alert(getCreateOrderErrorMessage(err))
      } finally {
        creatingOrder.value = false
      }
    })
  } finally {
    preparingOrder.value = false
    if (!showPaymentModal.value) {
      activeOrderActionKey.value = ''
    }
  }
}

// === Recharge preset ===
async function handleRechargePreset(rp: RechargePlan) {
  if (orderFlowBusy.value) return
  activeOrderActionKey.value = makeOrderActionKey('recharge', rp.key)

  const code = promoCode.value.trim()
  if (code) {
    preparingOrder.value = true
    const valid = await validatePromoForAmount(code, rp.pay_amount_fen)
    if (!valid) {
      preparingOrder.value = false
      activeOrderActionKey.value = ''
      return
    }
  }

  try {
    await showPayMethodOrDirect(async () => {
      paymentOrderType.value = 'recharge'
      creatingOrder.value = true
      try {
        const payYuan = rp.pay_amount_fen / 100
        const order = await paymentAPI.createRechargeOrder(payYuan, code || undefined, selectedPayMethod.value, rp.key)
        await showQRModal(order)
      } catch (err: any) {
        const msg = err?.message || err?.response?.data?.message || t('recharge.payment.createFailed')
        alert(msg)
      } finally {
        creatingOrder.value = false
      }
    })
  } finally {
    preparingOrder.value = false
    if (!showPaymentModal.value) {
      activeOrderActionKey.value = ''
    }
  }
}

// === Custom recharge ===
async function handleCustomRecharge() {
  if (customFinalAmount.value <= 0 || orderFlowBusy.value) return
  if (rechargeMinAmount.value > 0 && customFinalAmount.value < rechargeMinAmount.value) return
  activeOrderActionKey.value = makeOrderActionKey('custom')

  const code = promoCode.value.trim()
  const amountFen = Math.round(customFinalAmount.value * 100)
  if (code) {
    preparingOrder.value = true
    const valid = await validatePromoForAmount(code, amountFen)
    if (!valid) {
      preparingOrder.value = false
      activeOrderActionKey.value = ''
      return
    }
  }

  try {
    await showPayMethodOrDirect(async () => {
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
    })
  } finally {
    preparingOrder.value = false
    if (!showPaymentModal.value) {
      activeOrderActionKey.value = ''
    }
  }
}

// === Shared QR modal logic ===
async function showQRModal(order: { order_no: string; code_url: string | null; amount_fen: number; expired_at: string }) {
  currentOrderNo.value = order.order_no
  currentOrderAmount.value = (order.amount_fen / 100).toFixed(order.amount_fen % 100 === 0 ? 0 : 2)
  paymentStatus.value = 'pending'
  paymentSuccessHandled.value = false
  qrLoading.value = true
  showPaymentModal.value = true

  const expiresAt = new Date(order.expired_at).getTime()
  countdown.value = Math.max(0, Math.floor((expiresAt - Date.now()) / 1000))

  try {
    if (order.code_url) {
      await nextTick()
      if (qrCanvas.value) {
        await QRCode.toCanvas(qrCanvas.value, order.code_url, {
          width: 192,
          margin: 2,
          color: { dark: '#000000', light: '#ffffff' },
        })
      }
    }
  } finally {
    qrLoading.value = false
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
    await checkCurrentOrderStatus()
  }, 3000)
}

async function checkCurrentOrderStatus() {
  if (!showPaymentModal.value || paymentStatus.value !== 'pending' || !currentOrderNo.value) {
    return
  }

  try {
    const order = await paymentAPI.queryOrder(currentOrderNo.value)
    if (order.status === 'paid' && !paymentSuccessHandled.value) {
      paymentSuccessHandled.value = true
      paymentStatus.value = 'paid'
      clearTimers()
      await syncPaymentSuccessState()
    } else if (order.status === 'closed') {
      paymentStatus.value = 'closed'
      clearTimers()
    }
  } catch {
    // ignore poll errors
  }
}

function handlePaymentReturn() {
  void checkCurrentOrderStatus()
}

function handlePaymentVisibilityChange() {
  if (!document.hidden) {
    void checkCurrentOrderStatus()
  }
}

async function syncPaymentSuccessState() {
  const refreshTasks: Promise<unknown>[] = [
    authStore.refreshUser(),
  ]

  if (paymentOrderType.value === 'subscription') {
    refreshTasks.push(subscriptionStore.fetchActiveSubscriptions(true))
  }

  const results = await Promise.allSettled(refreshTasks)
  const failed = results.find(result => result.status === 'rejected')
  if (failed) {
    console.error('Failed to refresh payment success state:', failed.reason)
    appStore.showError('支付成功，但页面数据刷新失败，请手动刷新页面查看最新状态')
    return
  }

  appStore.showSuccess(
    paymentOrderType.value === 'subscription'
      ? t('pricing.payment.paymentSuccessToast')
      : t('pricing.payment.rechargeSuccessToast'),
    5000
  )
}

function cancelPayment() {
  showPaymentModal.value = false
  clearTimers()
  activeOrderActionKey.value = ''
}

function goAfterPayment() {
  showPaymentModal.value = false
  clearTimers()
  activeOrderActionKey.value = ''
  router.push({
    path: '/order-history',
    query: {
      highlight: currentOrderNo.value,
      success: '1'
    }
  })
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

</style>
