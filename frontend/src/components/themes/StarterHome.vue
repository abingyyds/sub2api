<template>
  <!-- Header -->
  <header class="relative z-20 px-6 py-4">
    <nav class="mx-auto flex max-w-6xl items-center justify-between">
      <div class="flex items-center">
        <div class="h-10 w-10 overflow-hidden rounded-xl shadow-md">
          <img :src="props.siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
        </div>
      </div>
      <div class="flex items-center gap-3">
        <LocaleSwitcher />
        <a
          v-if="props.docUrl"
          :href="props.docUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
          :title="t('home.viewDocs')"
        >
          <Icon name="book" size="md" />
        </a>
        <button
          @click="emit('toggle-theme')"
          class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
          :title="props.isDark ? t('home.switchToLight') : t('home.switchToDark')"
        >
          <Icon v-if="props.isDark" name="sun" size="md" />
          <Icon v-else name="moon" size="md" />
        </button>
        <router-link
          v-if="props.isAuthenticated"
          :to="props.dashboardPath"
          class="inline-flex items-center gap-1.5 rounded-full bg-gray-900 py-1 pl-1 pr-2.5 transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
        >
          <span class="flex h-5 w-5 items-center justify-center rounded-full bg-gradient-to-br from-primary-400 to-primary-600 text-[10px] font-semibold text-white">
            {{ props.userInitial }}
          </span>
          <span class="text-xs font-medium text-white">{{ t('home.dashboard') }}</span>
        </router-link>
        <router-link
          v-else
          to="/login"
          class="inline-flex items-center rounded-full bg-gray-900 px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
        >
          {{ t('home.login') }}
        </router-link>
      </div>
    </nav>
  </header>

  <main class="relative z-10 flex-1 px-6 py-16">
    <div class="mx-auto max-w-6xl">
      <!-- Hero: Left text + right terminal animation -->
      <div class="mb-12 flex flex-col items-center justify-between gap-12 lg:flex-row lg:gap-16">
        <div class="flex-1 text-center lg:text-left">
          <h1 class="mb-4 text-4xl font-bold text-gray-900 dark:text-white md:text-5xl lg:text-6xl">
            {{ props.heroTitle }}
          </h1>
          <p class="mb-8 text-lg text-gray-600 dark:text-dark-300 md:text-xl">
            {{ props.heroDescription }}
          </p>
          <div>
            <router-link
              :to="props.isAuthenticated ? props.dashboardPath : '/login'"
              class="btn btn-primary px-8 py-3 text-base shadow-lg shadow-primary-500/30"
            >
              {{ props.isAuthenticated ? t('home.goToDashboard') : props.ctaText }}
              <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
            </router-link>
          </div>
        </div>
        <div class="flex flex-1 justify-center lg:justify-end">
          <div class="terminal-window">
            <div class="terminal-header">
              <div class="terminal-buttons">
                <span class="btn-close"></span>
                <span class="btn-minimize"></span>
                <span class="btn-maximize"></span>
              </div>
              <span class="terminal-title">terminal</span>
            </div>
            <div class="terminal-body">
              <div class="code-line line-1">
                <span class="code-prompt">$</span>
                <span class="code-cmd">curl</span>
                <span class="code-flag">-X POST</span>
                <span class="code-url">/v1/messages</span>
              </div>
              <div class="code-line line-2">
                <span class="code-comment"># Routing to upstream...</span>
              </div>
              <div class="code-line line-3">
                <span class="code-success">200 OK</span>
                <span class="code-response">{ "content": "Hello!" }</span>
              </div>
              <div class="code-line line-4">
                <span class="code-prompt">$</span>
                <span class="cursor"></span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Feature Tags -->
      <div class="mb-12 flex flex-wrap items-center justify-center gap-4 md:gap-6">
        <div
          v-for="(tag, i) in props.featureTags"
          :key="i"
          class="inline-flex items-center gap-2.5 rounded-full border border-gray-200/50 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/80"
        >
          <Icon :name="tagIcons[i] ?? 'check'" size="sm" class="text-primary-500" />
          <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ tag }}</span>
        </div>
      </div>

      <div v-if="props.registrationNotice || props.allowSubSiteOpen" class="mb-12 grid gap-4 lg:grid-cols-2">
        <div v-if="props.registrationNotice" class="rounded-2xl border border-amber-200 bg-amber-50 px-5 py-4 text-sm text-amber-700 dark:border-amber-900/40 dark:bg-amber-950/20 dark:text-amber-300">
          {{ props.registrationNotice }}
        </div>
        <div v-if="props.allowSubSiteOpen" class="rounded-2xl border border-primary-200 bg-white/70 px-5 py-4 shadow-sm backdrop-blur-sm dark:border-primary-900/30 dark:bg-dark-800/60">
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p class="text-sm font-semibold text-gray-900 dark:text-white">当前分站支持自助开通下级分站</p>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">售价 ￥{{ props.subSiteOpenPrice }}，站长可自行设置模板和售卖价格。</p>
            </div>
            <router-link :to="props.isAuthenticated ? '/subsites' : '/login'" class="btn btn-primary whitespace-nowrap">
              {{ props.isAuthenticated ? '进入分站中心' : '登录后开通' }}
            </router-link>
          </div>
        </div>
      </div>

      <slot name="body" />
    </div>
  </main>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import type { ThemeHomeProps, ThemeHomeEmits } from './types'

const props = defineProps<ThemeHomeProps>()
const emit = defineEmits<ThemeHomeEmits>()
const { t } = useI18n()

const tagIcons: Array<'swap' | 'shield' | 'chart'> = ['swap', 'shield', 'chart']
</script>

<style scoped>
.terminal-window {
  width: 420px;
  background: linear-gradient(145deg, #1e293b 0%, #0f172a 100%);
  border-radius: 14px;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
  transform: perspective(1000px) rotateX(2deg) rotateY(-2deg);
  transition: transform 0.3s ease;
}
.terminal-window:hover {
  transform: perspective(1000px) rotateX(0) rotateY(0) translateY(-4px);
}
.terminal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: rgba(30, 41, 59, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}
.terminal-buttons {
  display: flex;
  gap: 8px;
}
.terminal-buttons span {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}
.btn-close { background: #ef4444; }
.btn-minimize { background: #eab308; }
.btn-maximize { background: #22c55e; }
.terminal-title {
  flex: 1;
  text-align: center;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #64748b;
  margin-right: 52px;
}
.terminal-body {
  padding: 20px 24px;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 14px;
  line-height: 2;
}
.code-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  opacity: 0;
  animation: line-appear 0.5s ease forwards;
}
.line-1 { animation-delay: 0.3s; }
.line-2 { animation-delay: 1s; }
.line-3 { animation-delay: 1.8s; }
.line-4 { animation-delay: 2.5s; }
@keyframes line-appear {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}
.code-prompt { color: #22c55e; font-weight: bold; }
.code-cmd { color: #38bdf8; }
.code-flag { color: #a78bfa; }
.code-url { color: #14b8a6; }
.code-comment { color: #64748b; font-style: italic; }
.code-success { color: #22c55e; background: rgba(34, 197, 94, 0.15); padding: 2px 8px; border-radius: 4px; font-weight: 600; }
.code-response { color: #fbbf24; }
.cursor { display: inline-block; width: 8px; height: 16px; background: #22c55e; animation: blink 1s step-end infinite; }
@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}
</style>
