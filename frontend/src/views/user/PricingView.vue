<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-6">
      <!-- Title -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('pricing.title') }}</h1>
        <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('pricing.subtitle') }}</p>
      </div>

      <!-- Tabs -->
      <div class="flex justify-center">
        <div class="inline-flex rounded-lg border border-gray-200 bg-gray-100 p-1 dark:border-dark-600 dark:bg-dark-800">
          <button
            class="rounded-md px-5 py-2 text-sm font-medium transition-all"
            :class="activeTab === 'subscription'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200'"
            @click="activeTab = 'subscription'"
          >
            {{ t('pricing.tabSubscription') }}
          </button>
          <button
            class="rounded-md px-5 py-2 text-sm font-medium transition-all"
            :class="activeTab === 'recharge'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200'"
            @click="activeTab = 'recharge'"
          >
            {{ t('pricing.tabRecharge') }}
          </button>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
      </div>

      <!-- Promo Code Input (shared across tabs) -->
      <div v-if="!loading" class="mx-auto max-w-lg">
        <div class="flex items-center gap-3">
          <div class="relative flex-1">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
              <Icon name="gift" size="md" :class="promoValidation.valid ? 'text-green-500' : 'text-gray-400 dark:text-dark-500'" />
            </div>
            <input
              v-model="promoCode"
              type="text"
              class="input pl-10 pr-10"
              :class="{
                'border-green-500 focus:border-green-500 focus:ring-green-500': promoValidation.valid,
                'border-red-500 focus:border-red-500 focus:ring-red-500': promoValidation.invalid
              }"
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
          <div v-if="promoValidation.valid" class="mt-2 flex items-center gap-2 rounded-lg bg-green-50 px-3 py-2 dark:bg-green-900/20">
            <Icon name="gift" size="sm" class="text-green-600 dark:text-green-400" />
            <span class="text-sm text-green-700 dark:text-green-400">
              {{ formatDiscountText() }}
            </span>
          </div>
          <p v-else-if="promoValidation.invalid" class="mt-1 text-sm text-red-500">
            {{ promoValidation.message }}
          </p>
        </transition>
      </div>

      <!-- ==================== Subscription Tab ==================== -->
      <template v-if="!loading && activeTab === 'subscription'">
        <div v-if="decoratedPlans.length > 0" class="grid gap-6 md:grid-cols-1 lg:grid-cols-3">
          <div
            v-for="dp in decoratedPlans"
            :key="dp.plan.key"
            class="relative flex flex-col overflow-hidden rounded-2xl border transition-all hover:shadow-xl hover:-translate-y-1 cursor-pointer"
            :class="dp.featured
              ? 'border-primary-500 ring-2 ring-primary-500 shadow-lg scale-[1.02]'
              : 'border-gray-200 dark:border-dark-600 shadow-sm'"
            @click="handleBuy(dp.plan)"
          >
            <!-- Top badge ribbon -->
            <div
              v-if="dp.badge"
              class="w-full py-1.5 text-center text-xs font-bold tracking-wide"
              :class="dp.featured
                ? 'bg-primary-600 text-white'
                : 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-300'"
            >
              {{ dp.badge }}
            </div>

            <div class="flex flex-1 flex-col p-6">
              <!-- Target audience -->
              <p class="text-xs font-medium text-primary-500 dark:text-primary-400">{{ dp.audience }}</p>

              <!-- Plan name -->
              <h3 class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ dp.plan.name }}</h3>

              <!-- Price -->
              <div class="mt-3 flex items-baseline gap-1">
                <span class="text-4xl font-extrabold text-gray-900 dark:text-white">¥{{ (dp.plan.amount_fen / 100).toFixed(0) }}</span>
                <span class="text-sm text-gray-500 dark:text-dark-400">/月</span>
              </div>

              <!-- Highlight text -->
              <p class="mt-3 text-sm font-medium text-primary-600 dark:text-primary-400">{{ dp.highlight }}</p>

              <!-- Divider -->
              <div class="my-4 border-t border-gray-100 dark:border-dark-600"></div>

              <!-- Feature list -->
              <ul class="flex-1 space-y-3">
                <li v-for="(feat, fi) in dp.features" :key="fi" class="flex items-start gap-2.5 text-sm text-gray-600 dark:text-dark-300">
                  <svg class="mt-0.5 h-4 w-4 flex-shrink-0 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                  <span>{{ feat }}</span>
                </li>
              </ul>

              <!-- CTA button -->
              <button
                class="mt-6 w-full rounded-xl py-3 text-sm font-bold transition-colors"
                :class="dp.featured
                  ? 'bg-primary-600 text-white hover:bg-primary-700 shadow-md'
                  : 'bg-gray-900 text-white hover:bg-gray-800 dark:bg-dark-500 dark:hover:bg-dark-400'"
                :disabled="creatingOrder"
              >
                {{ dp.cta }}
              </button>
            </div>
          </div>
        </div>
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
        <div v-if="rechargePlans.length > 0" class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          <div
            v-for="rp in rechargePlans"
            :key="rp.key"
            class="card relative flex flex-col overflow-hidden transition-all hover:shadow-lg hover:-translate-y-1 cursor-pointer"
            :class="{ 'ring-2 ring-primary-500': rp.popular }"
            @click="handleRechargePreset(rp)"
          >
            <div v-if="rp.popular" class="absolute right-4 top-4">
              <span class="badge badge-primary text-xs">{{ t('pricing.mostPopular') }}</span>
            </div>
            <div class="flex flex-1 flex-col p-5">
              <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ rp.name }}</h3>
              <div class="mt-2 flex items-baseline gap-1">
                <span class="text-xs text-gray-500 dark:text-dark-400">{{ t('pricing.rechargeGet') }}</span>
                <span class="text-2xl font-bold text-primary-600">${{ rp.balance_amount }}</span>
              </div>
              <p v-if="rp.description" class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ rp.description }}</p>
              <div class="mt-auto pt-4">
                <div v-if="rp.pay_amount_fen !== rp.balance_amount * 100" class="mb-1 text-center">
                  <span class="text-xs text-gray-400 line-through dark:text-dark-500">¥{{ rp.balance_amount * 100 / 100 }}</span>
                </div>
                <div
                  class="flex items-center justify-center rounded-lg bg-primary-600 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
                  :class="{ 'opacity-50 pointer-events-none': creatingOrder }"
                >
                  {{ t('pricing.rechargeNow') }} ¥{{ (rp.pay_amount_fen / 100).toFixed(rp.pay_amount_fen % 100 === 0 ? 0 : 2) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Custom Amount Input -->
        <div class="card mx-auto max-w-lg p-6 space-y-5">
          <h3 class="text-center text-lg font-semibold text-gray-900 dark:text-white">{{ t('pricing.customRecharge') }}</h3>

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
                :placeholder="rechargeMinAmount > 0 ? t('pricing.minAmountHint', { amount: rechargeMinAmount }) : t('recharge.inputPlaceholder')"
                class="w-full rounded-lg border border-gray-300 bg-white py-3 pl-8 pr-4 text-lg font-medium text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-800 dark:text-white dark:placeholder-dark-500"
              />
            </div>
            <p v-if="rechargeMinAmount > 0" class="mt-1 text-xs text-gray-400 dark:text-dark-500">
              {{ t('pricing.minAmountNote', { amount: rechargeMinAmount }) }}
            </p>
          </div>

          <!-- Amount Summary -->
          <div v-if="customFinalAmount > 0" class="rounded-lg bg-gray-50 dark:bg-dark-800 p-4">
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600 dark:text-dark-400">{{ t('recharge.rechargeAmount') }}</span>
              <span class="text-lg font-bold text-primary-600">¥{{ customFinalAmount.toFixed(2) }}</span>
            </div>
          </div>

          <!-- Submit -->
          <button
            class="w-full rounded-lg bg-primary-600 py-3 text-base font-medium text-white transition-colors hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
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
        </div>
      </template>
    </div>

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
              <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAJQAAACUCAMAAABC4vDmAAAAeFBMVEUtwQD///8AvAAdvwD8/vzj9eDx+u+05LJkzVXQ7szW8NSM2IVHxjFQyEmF1YAuwQ/3/Pbs+Oqe3Zfn9uTI68Pb8tml36O+6Ll20m5703Rozl4ywh2t4qYvwSV00WiW245RyUFWyVGm4J41wjGF1npAxDx30XhmzWUzHtIeAAAHK0lEQVR4nM1c6barIAxU4r7UtbVq7d7b93/DK2p3QBT0c/71tEenIYQwgSjqJwzHdC1PmQ2e5ZqO8UVC+fikuz4ggPk4KQrUL/RdnUrKttCsfN6YIcsmktJW85roixastF9SdvCHlBpagf1NyvxjSg0t85NUjP6aEQaK30mZi+BUszJfpOwFjF0LsB+ktOCvubwQaB2p1WIMVZtq1ZJazuBh4AGsSVnLImVhUvpCZt4DSK9JuYsyVG0qV1UM/69ZfMM3FGdhhqpN5SjDgzk8MAUjBYd1ZYhL1QkZQsrR2uZ5vi2OQf1pAmrgKrwBAVDg3w9p8pEiRmm89UvJqSFYCl8+jlDxm0q3CJ0sR1J5cVFCZWGGRELPrDXde7OGO4TiNZNRi2QTzEYLncx+Qh121iy0EGTscftEaE6f6EO1G8CoRVxOSwtt9X4SP1j7E44hVNkIShjnyWyFjtFITqqaVtPQemx2xsE4TcFKjFM9DQv5rNBYd3rhKpsV8MdLKrRcLis0PDqRsJXJCg5SOKmaL48V3LX+F3LBuEkjFZCzpjGQlv5DKo1THdvlrDiwlchJ1e5SbAXjFxcSpOy+YSOVk6oeJLC6yPPyFpF4nQDJNpSMPAZkG6r2ddHxo8fyZK335OqRrpP/kbDubBMfq+q5B+WdFcDC+AaBT1wz14KcfPIb02b7C+hMN1OrKKCC8J0mtjAj8ug9JxA4NDs9UjpifiHm6ighPFJVNw+nwGokEfbTbUiJvSHkVB7Zl1/ZGqKQip+/ANIGX8RS4BLf+OYTNFIum5RIZgyUhPPwfGZJIfUSAyuSB4gUfgJKQHhq/9Q0Oaoe/2tFfIAAqYqWILjtQ+FEjZ+dLQAR4+e6HE/qQnuluqmw1Hmlfq+q2bH+BbqTp6+IKO7R35nYu5Qtdhj2znS+k3stXGfn674YbylqFBoHPd3sW/VYKCLIS4S1KC4uIEObhVwSpcT00dt8E6oESCKVXKuOQE0kqC5esV25q+veul2q4dRgL4GSs+oGDRCyDqb9HmW0JN25H0bkICXu6Mah1WEBKiuLiFEtjHZ+NcBg9DjFibRCrZH8HXOjpm987kkgSCrcdlYq7F4xQrMLTmsRF1NurI+NmdCJc9uf8hUChFSErPnng4R30+MwlsjuuE15aYsfBUbebyzIRwtT7eaclpDRkfWXJ8qxW9GWUzVi+Ne9Q4goWV4fmu0mXIZUll7o0yBHxvTmOBgcR5o57CsXU3NPFppkF47jKGH0sBqjn4eNV4hISGGPX1F2fiw0m+PRU6RlVTFJocETKGscSrBqkrFNNXQKRfiMGkP54ARbLhoa1VcIO7lwNSC8MW1FETkocLChApoYMwDsw27Dtg/4MB+w9oPc2LNZDSg/Jk2I+nJDJ/e8AyXgZfvSIi+QEXttJgonZGDp41snTRWoc0+PyOqA8HfkDUrPCUruuB7iKlX1uXPuFDIgKZVp+16yLm6zz5pyzyb8EviSOXfdMJQE7++MARbx+UcmqfpFfNEKj953uH0IZKR4+ljkSuLj+1QsKHhspTUC7dcLHmGQFCc6wnAiPj3qzaz2HLbCE+YnW43G+pSqsaMCfuixn1WGSBmwiXcQEBCFo2sz+2jD0C9tg9cbGbDfBr8reLr1ji5lXTALr9jQXKNnWW5Y9S79hUIOH1pIt7IW0r014dHW0JXp7tqp/s2th/ggXDhI1dOEtTpHOHRSqjkTkmILMfrlN3SKgfPEKVmCbuGUsmvznOUuYIxfQ0qWJNmA88gOMB5hl5LUv4GkmC7TWEqqzN0b0ht8rrZGkhhvEcjBGyOBPegvjnyj9xbV0/x0O57u/1Zxtm6o6Xg9nj9OKafn5Ft3KiUWxHEVJijcM7aUmPr3iZ496cNQ/7qfa+7P6dtOEpd5hsjhuufU7SG0jC6agvA+9IUdn583epXBOh0JW1nH0lTOqzKofqEWsyVvVswfBoPvqE5QO3jRkzlLOBHaweGTsN3w0F8ZkJYn8K18kPOciZShJGBEXAGhZsX1I2KBfThivtWYE8N0GhoMgYo8AXIyBRln9z5YSYjqIkcXyLiIB1D5lzHFZTNJJ2c/WQmeX5/mJi3iF9oISCa6PTKqgNIh4ilGjgGMrzlEE91nUQSqWMlUdmpYjTufrU/bYAPKETXMZOLrZKNukUw5dh0GJ8ez9I6AclAiqhUzcKqB8gGZTHSch5QC6Mod3qP5rtVD+Y83Rb7PRgpbqzw7JDU2dMyth5C3StugtpvD0V9Ayik3P8X0MLuegock4ceYNPWuwWRtnbASEli5G5/PsZtbyudVc4TOCS2R8rgv3I9kRu9jAJedRsykwBrUmkAyEBxIuTm4I5o4SATRIMhcZruLRTYGWWQLlWU2m1lkW55lNjBaZKunZTbFWtAAvtqHLbPR2jJb0i2zed8y2xwusyGkusjWmRgLbDKKsYh2rP8BI6Zi4Ga6nYMAAAAASUVORK5CYII=" class="h-8 w-8" alt="微信支付" />
              <span class="font-medium text-gray-900 dark:text-white">微信支付</span>
            </div>

            <!-- Alipay -->
            <div
              class="flex items-center gap-3 rounded-lg border-2 p-4 cursor-pointer transition-all"
              :class="selectedPayMethod === 'alipay' ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20' : 'border-gray-200 dark:border-dark-600 hover:border-primary-300'"
              @click="selectedPayMethod = 'alipay'"
            >
              <input type="radio" :checked="selectedPayMethod === 'alipay'" class="h-5 w-5 text-primary-600" />
              <img src="data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAkGBwgHBgkIBwgKCgkLDRYPDQwMDRsUFRAWIB0iIiAdHx8kKDQsJCYxJx8fLT0tMTU3Ojo6Iys/RD84QzQ5OjcBCgoKDQwNGg8PGjclHyU3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3N//AABEIAJQAlAMBEQACEQEDEQH/xAAbAAEAAgMBAQAAAAAAAAAAAAAABQYBBAcDAv/EAEAQAAEDAwAECgcHAwQDAAAAAAEAAgMEBQYSEiExFiJBUVJhcZGT0RMUMlWBoeEHF0JicrHBQ1NzMzRj8CM0Nf/EABsBAQADAQEBAQAAAAAAAAAAAAAEBQYDAgEH/8QANREAAQMCAwUGBgEFAQEAAAAAAAECAwQRBRIhFTEyYZFBUVJToeETFiNxgbHRFDM0QsHwI//aAAwDAQACEQMRAD8A7igIi+3OltsYlqX4yMNYPad2LtBTvmdZiEapq4qZmaRSiXHS6vqHYpCKWPl1drj8VeQ4ZEzV/wBSmZqcankf0p6kLPVz1BLp55ZHc73kqeyNjOFEQq3zPfxqqngvZyCAIAgCAIAgCAIAgPpj3MOWOc087ThfFRF3nprlauhJ0mkV1pMalXJIzoS8YfPaostFBJvb0JsGJVMK6Ov99S2WXSimuDhDUgU852DJ4ruw8ip6nDnxJmbqhoqLF4qhcr/pd6KXBh4rexVxbn0gIi7Xllton1MzQcbGsB2uPMu1PA6d6MQjVdSymiWRxy65V9Rcat9TVv1nu3DkaOYdS0MLImZW7jEVFQ+eRXvXU1F1I4QBAEAQBAEAQBAEAQBAEAQGUBd9D9JpHFlurXazsYhkdvP5SqPEKJE/+rPyn/TT4TiSvX4Eq69i/8Lh607otVOaE5zprXOqboaYbI6bi4zvccZP8LRYZDkizrvcZDGqlZJ/hpub+yuqyKYIAgCAIAgCAtNo0MnuVvhq/XI4hKNYN1C7Z8lVz4m2KRWZb2LumwZ80SSZ0S5ufd9P7yj8I+a5bYb4PU77Af5noV3SC0Usta2mfMJSWB+sG43/9KwpahKhmdEsVVdSLSSJGq3IxSSEEAQBAEAQGWOLXBzSQ4HII5CvioipZT61VRUVDqljrHXS1w1Wrx3DD/wBQ2FZSqh+DKrDeUVT/AFEDZO3tOY1kxqKad2+SRzu8rUxsyMRvchh5X53q7vW54L2cggCAIAgCAIDrmin+3qD/AB D/3WTrf8h/3N3h3+LH9iXUYmnNftE/+5H/gb+5Wiwr+x+TJY7/kp9irKzKQIAgCAIAgAQIT1kv0lspHQNJwZC75DyUGppEmfmLSjrlgjycyBU4qwgCAkLRZ6u7vkZRhhdGAXa7sbCo9RUx06Ir+0l0tFLVKqR9hKcCbz0afxfoo21Kfn0Juw6vl1HAm89Gn8X6JtSn59BsOr5dRwJvPRp/F+ibUp+fQbDq+XUcCbz0afxU2pT8+g2JV8upf7DTSUVppaWfHpIo8O1TkKhqZGySue3cpp6SN0MDWO3ohILiSSl6XaO3C7XRtRSNi9GIg3L34OcnzVvQVsMMWV++5QYph09TMj47WsQnAm89Gn8X6KbtSn59Cu2HV8uo4E3no0/i/RNqU/PoNh1fLqOBN56NP4v0TalPz6DYdXy6njWaI3WjpZamZsPo4mlzsSZOAvceIwPejEvdTnLhFTExXutZOZAqcVYQBAEAQBAEBMaO3t9kkmkZA2b0rQMOdjGCodXSJUIiKtrFhQVy0iuVEvcneH83u6PxPooWyGeMstvu8scP5vd0fifRNkM8Y2+7yxw/m93R+J9E2Qzxjb7vLHD+X3dH4h8k2Q3xjb7vL9S4WasNfbaercwMMrNYtB3KoqI/hSuZe9i/ppVmhbIvaby5Hcq2kmlElmuDaZtKyUGMPyXY3k+SsqSgbUR5ldYpq/E1pZciNvoRXD+b3dH4n0UrZDPGQtvu8scP5vd0fifRNkM8Y2+7yxw/m93R+IfJNkM8Y2+7yzWuOmstdQz0poo2CZhYXB5OMrpFhjY3o9HbjlPjTpYnR5LXKmVaoUQQBAEAQBAEAQEnaLHW3gSGiawiPGtrOxvUaerjgtn7SbS0E1VdY+wkeBN56EHiqPtSn59CXsSq5dRwJvPQg8VNqU/PofNiVfLqOBN56EHiptSn59BsSq5dS/WCllorRS004Akjj1XYORlUNS9JJnPbuU1FJE6KBrHb0QkTuXEklL0u0duF2ujaijZGYxEG5c/G3JVvQVsUEWV++5QYnh09TMj491iE4E3noQeKpu1Kfn0K7YlVy6jgTeehB4qbUp+fQ+bEq+XU8qnRC60tNLUTNhEcTC92JOQDJXpmIwPcjUvdeR4kwepjYr3Wsmu8gFPKsIAgCAIAgCAIAgJ3R3SF9jbOGU4mEpB2v1cYUKro0qVS62sWdDiK0aORG3uTH3gS+72eKfJQ9jt8fp7k/5gd5fr7D7wJfdzPFPkmx2+P09x8wO8v19h94Evu5ninyTY7fH6e4+YHeX6+wH2gTe7meKfJNjt8fp7j5gd5fr7GfvBm93s8U+SbHb4/T3HzA7y/X2MfeBL7uZ4p8k2O3x+nuPmB3l+vsSujmk899przT+ptijazWc8PJx8lFq6FtOzNmuTaDE31cmTJZPv7Fowq0uSvacVYpbBMzPGnIiHx2n5Aqdh0eeoTlqVeLy/DpXJ36HLlpzFhAEAQBAEAQBAZQ+2PeGiqpxmGmmeOdrCV4dKxu9yHRkEj+FqqextFyaMmhqAP8AGVz/AKmHxIdP6Oo8C9DWlp54f9aGSP8AW0hdWva7cpydE9nEioeS9HiwQ+RAZCA6RoDb/VbUap4xJUu1h+kbv5WcxSbPLkTchr8Fp/hwZ13u/RaSqwuTnn2iV/pq+KiYeJC3Wd+o/T91f4TFlYsi9plccqM0qRJuT9lQVuUIQBAEAQBAEBYtHtFai6gTzOMFLyOI2v7PNV1XiDIPpbq79FtQ4U+pTO/6W/svVt0ftduAMFKwyDfJJxnHv3fBUktZNLxONNT0FPBwN17ySMsLDgyMad2C4BR8qrqSszU0ufbXsf7LmnsKKipvPqORdyhzWuGHAEda+IosikfWWO11oPrFFC4nlA1T3jau8dVNHwuI0tFTS8bEIOu0EoJQTRzS07uQHjt81NjxWVvGlytmwKB39tVT1K7X6G3Wly6GNtUznjOD3FWEeJwP36KVU+DVMfCmZORF0NrqKm6Q0D4nxyPeAdYEFo5T3KVJUMZEsl7ohCgpXyTtiVLKp2CCJkMLIoxhjGhrR1BZJzlcqqpvGtRrURNyHlcKqOjo5amY4ZE0uPkvUcbpHoxu9TzNK2JivduQ47XVMlZVzVM3tyvLj1dS10UbY2I1vYYCaV0sivd2muuhyCAIAgCAICc0Ss4u9x1ZhmnhGtJ18w+Kg19SsEem9SzwujSpm+rhTVS86QX6nsVOxjYw+dw/8cLdgAHKeYKkpaR9S7l3mkra+OjYidvYhz+46QXO4OPpqp7WZ2Rx8VoV/DRwxcLdTKz4hUTcTtO5NCLJztO0nlKlEO59RyviOY3uYedpwvitRd6H1r3N4VsStHpNd6Qj0dY97R+GUBwPeoslDBJvb0J0WKVUe51/vqWO26eMdhlypiz/AJIdo7lXS4S5NY3X+5bwY61dJW2+xa7fdKO4x69HUMlHKAdo7RvCq5YZIls9LF1DUxTpeN1zc3hcjufJjaXBxaNYbjjal13HzKl79p9bhhfD6VTTmK51VNFTUNM+WDOvK5hGSRuGN/WrTDnQsernrZewpcYZUSMRkTbpvU53LFJE8slY5jh+FwwVoWuRyXQybmOatnJZT5X08mEAQBAEAQIdJ+z6mbFZHTY400pJ7BsCzmKPvNl7kNfgcaNps3epRr9WPrrvVTvOcyFrRnc0bAruljSOFrUM3XSrLUOcveR6kEQIAgCAID0hlkhkEkL3seNzmHBC8ua1yWch6a9zVu1bKdA0Jul2uAf62WyUrNgmcMOLubrVBiMEES/Rv7jV4TVVU6L8TVqdvaW7tVWXYQDCA16ygpa2Msq6eOZvM9uV7ZK+NbsWxylgjlSz23K3cNBrfNl1I+Smdze0357fmrCLFZm6PS5VT4HA/WNVb6lcuGht1pMuijbUsHLEdvcVZRYnA/foVE+DVEWrUzJyICaKSF+pNG6N/Rc0g/NTmua7Vq3Kt7HMWzksp5r0eAgCA6boBK1+j7WA5Mcjmnvz/KzeKNVKhV7zY4K5FpUROxVKBeqV9HdaqB4xqyux1gnIV7TSJJE1ydxmKyJYp3tXvNFdyKEAQBAEBOaOaPVF4mDnAx0jTx5Mb+pvOVBrKxsCWTi7izw/Dn1Trro3vOn0lLDR00cFNG2OKMYa0cizb3ue5XOXU2MUTImIxiWQ+LjXQW+kfUVLwyNneTzDrX2KJ8r0YzeeZ5mQMV79yFHp9Oqpta908DH0zncVjdjmDt5VdPwlisTKupnGY7Ikiq9t2/ottq0ht1zAbTzhsp/pScV3dy/BVU9JLDxJp3l5TV8FRwO17l3krlRSaZQBAatXQU1YwsqoI5Wnke3K6MlfGt2rY5SQRypZ6IpXLhoLQzZdRSvpnn8PtN8/mrGLFZW6PS5Uz4HC/WNbL6Fem0Ku0chaz0MjeRweBnvU5uK06p9V+hVPwSpRbJZfyVpWZTlo0DuraK4PpJ3BsVTjBO4PG7v3dyq8Tp1kjzt3p+i7wWrSKX4btzv2WbSrRwXhgnpiGVcYwM7njmPmq2irVgXK7hUt8Sw5KpMzeJPU5zW0VTRSmOrgkhcOm3GeznWijlZIl2LcyUsEkK5ZG2NddDkEB7U1LUVUgjpoJJXn8LGkleHyNYl3LY6xRPlWzEuXGxaEOLmz3c4A2iBh3/qP8BVFTin+sPUv6PBNc0/T+S7wwxQRNiiY1kbRhrWjAAVK5yqt13miaxrERrUsiGtdLnS2umM9XIGt/C3O1x5gOVdIYXzPysQ5VFTHTszyKcwv98qbzUa8nEhb/pxA7B285WlpaRlO2yb+8xldXPqn3XRE3IRKlkEy0kEEEgjaCDuSx9RbFhtOl1xoNVkzhVQj8MntDsPmq6fDopdW6KWtLi88OjvqTnv6l1tWlFtuQDRMIJj/AE5SASeo8qp56GaHsunI0NNilPPpmsvcpNhwUMsTK+AIDGAvoOHuaWkh2wg4wtre5+cKljAOMYOE3hFsXzRnS+MsbS3Z+o8DDJzucPzcx61RVmHKiq+JPx/Bp8Pxdqokc62Xv/kt0kVPWRYkjjmjcNgcA4FVKOcxdFspeuayRuqIqdSOl0XssrusoIwfyZb+yktrqhv+xDdhlI7/AEMR6LWaN2W0MZP5iT/KOrqh29wbhlI1eBCTgpoKZgZTwxxM5mNA/ZRnPc9buW5MZGxiWalj0fIyNpL3BrRvJOAF5RFVbIelVES6lXvemdHRtdHQYqp92sDxG/Hl+Cs6fDZJLLJonqU9XjMUX0xfUvoUK43CquNQZ6uUyPO4Hc3sHIr2GBkLcrEsZeeokndmkW6/+3GoupwCAIAgCH1CbtWk9ytuGtl9ND/blyR8DvChT0EMutrLyLGmxSog0Rbp3KXS0aYW2u1WTO9VlOzEh4p7HKnnw6aPVNU5GhpcXp5tHLlXn/JYmvDgCCCDuIVeWiKipc+kPpx3SCm9UvNXEBxfSFzOw7QtZSSfEha4wdfD8Gpe3n+9SOUkhhBckLderhbf/UqntZ0Cct7lHlpYpeJpLgrZ4OBxPU+ntewYnpqeTHKCWk/uoLsJiXhcqFmzHpk4movp/JsH7QZMbLc3PXN9F42Onj9Dr8wL4PX2NSp06ucoIgip4Rz4Lj+/8LqzCYU4lVThJjtQ7hREIGvutfcDmrqpJB0c4b3DYp0VPFFwNsVk9XNP/cdc0l2IwQBAEAQBAEAQBAEAQBAEAQBAEAQBASdrv1xtZApZyY/7T+M3u5Pgo09JFNxpqTaavnp+B2ncp0mx3Z1xtsVVPD6J788UHIwDjKzVVE2CVWItzY0U7qiFJFS1yvaa2l1VTtroBmSEYe0De3n+Cn4ZU5HfDd2lbjNEsrElYmqb/t7FEV+ZMIAgCAIAgCAIAgCAIAgCAIAgCA3bTbpLnWsp4s4O17uiOdcKidII1epKpKV1TKjG/k6jDFHBEyKJoDGDVaOYLJucrlVVN4xqRtRrU0QknMaR7I7l8PRQdKtFTFK6qtbMsdxnwje3rb1dSu6LEUX6JV/P8maxHCFRfiwJ+CnEEEgjBBwQrkztlQwh8CAIAgCAIAgCAIAgCAIAgN212yquc3o6aPIHtPPst7VwnqI4Eu9SVS0ktS7KxDqFhstPaKQRQjXe/BkkI2vPks1U1L6h+ZTZ0dGylZlbv7VJTVb0R3KPYln0gNSq9tvYvgIqtsdvubwaqAa/TYdVylQ1U0OjFIlTQwVGr2695SNIbTBbZiyB8paDueQf4V7TVL5U+oylfSMp1sy5BqcVoQBAEAQBAEAQBAEB6QMD5WtJIycbF4e7Klz2xEc5EUvNo0Vtj6ZtRM2WV2zivfgfLCpKjEJ0VUatjT0eFU7mo9yKv5LFBFHBG2OFjY2AbGtGAqtznOW6qXjWNYiI1LISTPYb2Lyej6QH//Z" class="h-8 w-8" alt="支付宝" />
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
              {{ paymentOrderType === 'subscription' ? t('pricing.payment.paymentSuccess') : t('recharge.payment.success') }}
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
              @click="goAfterPayment"
            >
              {{ paymentOrderType === 'subscription' ? t('pricing.payment.viewSubscription') : t('recharge.payment.viewDashboard') }}
            </button>
          </div>
        </div>
      </div>
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
const selectedPayMethod = ref<'wechat' | 'alipay'>('wechat')
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
  // Match by amount_fen (29900 = ¥299)
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
  // Match by amount_fen (58900 = ¥589)
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
  // Match by amount_fen (118900 = ¥1189)
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
    badge: '个人首选',
    audience: '适合：独立开发者 / 个人项目',
    highlight: '总额$600：官方价格体系，无需复杂换算',
    features: [
      '每日 $30 高速配额：专属通道，响应丝滑无卡顿',
      '总额 $600/月：充足算力，轻松覆盖日常开发',
      '永不掉线保障：额度用尽自动转按量，业务 7x24 在线',
      '超高性价比：个人开发首选，一次付费全月无忧',
      '微信扫码支付：便捷安全',
    ],
    cta: '购买开发者版',
    featured: false,
  },
  '58900': {
    badge: '生产环境推荐 / 主推套餐',
    audience: '适合：全职开发 / 生产级应用',
    highlight: '折合 $1=¥0.45，立省 55% 成本！',
    features: [
      '每日 $60 高速配额：生产级优先级，拒绝排队等待',
      '总额 $1,300/月：充足算力储备，应对高频调用',
      '无损计费机制：拒绝倍率陷阱，每一分实打实可用',
      '熔断智能防护：异常流量自动隔离，生产环境稳定如磐石',
      '永不掉线保障：额度用尽自动转按量，业务不中断',
      '微信扫码支付：便捷安全',
    ],
    cta: '立即购买，释放生产力',
    featured: true,
  },
  '118900': {
    badge: '团队专属 / SLA 保障',
    audience: '适合：开发团队 / AI 产品项目',
    highlight: '总额$2,500：顶配储备，支撑海量高并发接入',
    features: [
      '每日 $120 高速配额：团队级 VIP 优先级调度，秒级响应',
      '总额 $2,500/月：顶配算力，满足团队级高频需求',
      '高可用服务保障：极致架构，满足企业级稳定运营标准',
      '无感溢出续航：额度海量储备，确保 AI 工作流永不断供',
      '微信扫码支付：便捷安全',
    ],
    cta: '购买旗舰版',
    featured: false,
  },
}

const decoratedPlans = computed<DecoratedPlan[]>(() => {
  // Target plans by amount_fen: 29900, 58900, 118900
  const targetAmounts = [29900, 58900, 118900]
  const matched = plans.value.filter(p => targetAmounts.includes(p.amount_fen))

  // If none of the target amounts match, fall back to showing all plans with generic decoration
  const source = matched.length > 0 ? matched : plans.value

  return source.map(plan => {
    const key = String(plan.amount_fen)
    const deco = planDecoZh[key] || planDecorations[key]
    if (deco) {
      return { ...deco, plan }
    }
    // Fallback for undecorated plans
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

  // We don't auto-validate here because we need an amount_fen.
  // Validation happens when user clicks buy.
}

async function validatePromoForAmount(code: string, amountFen: number): Promise<boolean> {
  if (!code.trim()) return true // No promo code = ok

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

  // Validate promo code if entered
  const code = promoCode.value.trim()
  if (code) {
    const valid = await validatePromoForAmount(code, plan.amount_fen)
    if (!valid) return
  }

  // Show payment method selection modal
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

  // Validate promo code if entered
  const code = promoCode.value.trim()
  if (code) {
    const valid = await validatePromoForAmount(code, rp.pay_amount_fen)
    if (!valid) return
  }

  // Show payment method selection modal
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

  // Validate promo code if entered
  const code = promoCode.value.trim()
  const amountFen = Math.round(customFinalAmount.value * 100)
  if (code) {
    const valid = await validatePromoForAmount(code, amountFen)
    if (!valid) return
  }

  // Show payment method selection modal
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
</style>