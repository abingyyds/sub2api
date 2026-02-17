<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950">
    <!-- Navigation -->
    <nav class="sticky top-0 z-50 border-b border-gray-200/50 bg-white/80 backdrop-blur-lg dark:border-dark-800/50 dark:bg-dark-900/80">
      <div class="mx-auto max-w-4xl px-6 py-4">
        <div class="flex items-center justify-between">
          <router-link to="/home" class="flex items-center gap-3">
            <Logo :size="32" theme="auto" />
            <span class="text-lg font-bold text-gray-900 dark:text-white">Sub<span class="text-primary-600">Router</span><span class="text-primary-600">.ai</span></span>
          </router-link>
          <router-link to="/home" class="text-sm text-gray-500 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">
            &larr; トップページへ戻る
          </router-link>
        </div>
      </div>
    </nav>

    <!-- Content -->
    <main class="mx-auto max-w-4xl px-6 py-16">
      <h1 class="mb-2 text-3xl font-bold text-gray-900 dark:text-white">特定商取引法に基づく表記</h1>
      <p class="mb-12 text-sm text-gray-500 dark:text-dark-400">Legal Disclosure under the Specified Commercial Transactions Act</p>

      <div class="space-y-8">
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">販売事業者</h2>
          <dl class="space-y-4 rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <div>
              <dt class="text-sm text-gray-500 dark:text-dark-400">事業者名</dt>
              <dd class="mt-1 font-medium text-gray-900 dark:text-white">SubRouter</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-dark-400">責任者</dt>
              <dd class="mt-1 font-medium text-gray-900 dark:text-white">HanShi Li</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-dark-400">所在地</dt>
              <dd class="mt-1 text-gray-700 dark:text-dark-300">ご請求があった場合は遅滞なく開示いたします。お問い合わせください。</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-dark-400">メールアドレス</dt>
              <dd class="mt-1 font-medium text-primary-600 dark:text-primary-400">support@subrouter.ai</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-dark-400">販売URL</dt>
              <dd class="mt-1"><a href="https://subrouter.ai" class="text-primary-600 hover:text-primary-700 dark:text-primary-400">https://subrouter.ai</a></dd>
            </div>
          </dl>
        </section>

        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">サービス内容</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <p class="text-gray-700 dark:text-dark-300">SubRouter は、AI API 推論サービスのルーティングおよび課金インフラを提供するソフトウェアサービスです。</p>
          </div>
        </section>

        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">販売価格</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <p class="text-gray-700 dark:text-dark-300">各サービスページまたはコンソールに表示される料金体系に基づき、サブスクリプション制または従量課金制にて提供いたします。</p>
          </div>
        </section>

        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">支払方法</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <p class="text-gray-700 dark:text-dark-300">クレジットカード（Stripe 決済）</p>
          </div>
        </section>

        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">サービス提供時期</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <p class="text-gray-700 dark:text-dark-300">お支払い完了後、直ちにサービスをご利用いただけます。</p>
          </div>
        </section>

        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">返品・返金ポリシー</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <p class="text-gray-700 dark:text-dark-300">デジタルサービスの性質上、原則として返金には対応しておりません。技術的な問題が発生した場合は、カスタマーサポートまでお問い合わせください。</p>
          </div>
        </section>

        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">解約について</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <p class="text-gray-700 dark:text-dark-300">ユーザーコンソールからいつでもサブスクリプションを解約できます。解約後、次の請求サイクルから課金は停止されます。</p>
          </div>
        </section>
      </div>
    </main>

    <!-- Footer -->
    <footer class="border-t border-gray-200 px-6 py-8 dark:border-dark-800">
      <div class="mx-auto max-w-4xl">
        <div class="flex flex-col items-center gap-4">
          <div class="flex flex-wrap items-center justify-center gap-4 text-xs text-gray-400 dark:text-dark-500">
            <router-link to="/legal/tokushoho" class="font-medium text-gray-600 dark:text-dark-300">特定商取引法に基づく表記</router-link>
            <span>|</span>
            <router-link to="/legal/disclosure" class="transition-colors hover:text-gray-600 dark:hover:text-dark-300">商業披露</router-link>
          </div>
          <p class="text-xs text-gray-400 dark:text-dark-500">&copy; {{ currentYear }} SubRouter. All rights reserved.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Logo from '@/components/Logo.vue'

const currentYear = computed(() => new Date().getFullYear())
</script>
