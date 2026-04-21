<template>
  <div class="mb-6 overflow-hidden rounded-2xl border border-amber-200 bg-gradient-to-r from-amber-50 via-white to-sky-50 shadow-sm dark:border-amber-500/30 dark:from-amber-500/10 dark:via-dark-900 dark:to-sky-500/10">
    <div class="p-5">
      <div class="flex items-start gap-3">
        <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-xl bg-amber-500 text-white shadow-sm">
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 6h9.75M10.5 12h9.75M10.5 18h9.75M3.75 6h.008v.008H3.75V6zm0 6h.008v.008H3.75V12zm0 6h.008v.008H3.75V18z" />
          </svg>
        </div>

        <div class="min-w-0 flex-1">
          <div class="flex flex-wrap items-center gap-2">
            <h2 class="text-base font-semibold text-gray-900 dark:text-white">路线通知</h2>
            <span class="inline-flex items-center rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-800 dark:bg-amber-500/20 dark:text-amber-200">
              API 调用链接
            </span>
          </div>

          <p class="mt-1 text-sm leading-6 text-gray-600 dark:text-dark-300">
            路线根据网络情况选择。API 调用时，在所选路线后追加对应接口后缀即可。
          </p>

          <div class="mt-4 grid gap-3 md:grid-cols-2">
            <a
              v-for="route in routes"
              :key="route.url"
              :href="route.url"
              target="_blank"
              rel="noopener noreferrer"
              class="rounded-xl border border-gray-200 bg-white/80 p-4 transition-all hover:-translate-y-0.5 hover:border-primary-300 hover:shadow-sm dark:border-dark-700 dark:bg-dark-900/80 dark:hover:border-primary-700"
            >
              <div class="flex items-center gap-2">
                <p class="text-sm font-medium text-gray-900 dark:text-white">{{ route.label }}</p>
                <span
                  v-if="route.recommended"
                  class="inline-flex items-center rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700 dark:bg-green-500/20 dark:text-green-300"
                >
                  推荐
                </span>
              </div>
              <p class="mt-2 break-all font-mono text-sm text-primary-700 dark:text-primary-300">{{ route.url }}</p>
            </a>
          </div>

          <div class="mt-4 rounded-xl border border-blue-200 bg-blue-50/80 p-4 dark:border-blue-500/30 dark:bg-blue-500/10">
            <p class="text-sm font-medium text-blue-900 dark:text-blue-200">API 调用链接示例</p>
            <ul class="mt-2 space-y-1 text-sm text-blue-800 dark:text-blue-300">
              <li><strong>Claude / OpenAI：</strong><code class="rounded bg-white px-1.5 py-0.5 text-xs dark:bg-dark-800">所选路线 + /v1</code></li>
              <li><strong>Gemini：</strong><code class="rounded bg-white px-1.5 py-0.5 text-xs dark:bg-dark-800">所选路线 + /v1beta</code></li>
              <li><strong>示例：</strong><code class="rounded bg-white px-1.5 py-0.5 text-xs dark:bg-dark-800">https://airiver.cn/v1</code>、<code class="rounded bg-white px-1.5 py-0.5 text-xs dark:bg-dark-800">https://api.airiver.cn/v1beta</code></li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const routes = [
  { label: '全球优化路线', url: 'https://ccoder.me', recommended: false },
  { label: '大陆优化路线', url: 'https://www.ccoder.me', recommended: false },
  { label: '大陆新路线 1', url: 'https://airiver.cn', recommended: true },
  { label: '大陆新路线 2', url: 'https://api.airiver.cn', recommended: true },
]
</script>
