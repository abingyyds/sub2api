<template>
  <!-- Aurora: floating header + centered vertical hero with animated aurora backdrop -->
  <header class="relative z-20 mx-auto mt-4 w-[min(1100px,95%)] rounded-3xl border border-white/40 bg-white/50 px-6 py-3 shadow-xl shadow-indigo-500/10 backdrop-blur-2xl dark:border-white/10 dark:bg-slate-900/50">
    <nav class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="h-10 w-10 overflow-hidden rounded-2xl ring-2 ring-sky-400/40">
          <img :src="props.siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
        </div>
        <span class="hidden font-semibold tracking-wide text-slate-700 dark:text-slate-100 sm:inline">{{ props.siteName }}</span>
      </div>
      <div class="flex items-center gap-2">
        <LocaleSwitcher />
        <a v-if="props.docUrl" :href="props.docUrl" target="_blank" rel="noopener noreferrer"
          class="rounded-xl p-2 text-slate-500 hover:bg-white/60 dark:text-slate-300 dark:hover:bg-white/5"
          :title="t('home.viewDocs')"><Icon name="book" size="md" /></a>
        <button @click="emit('toggle-theme')"
          class="rounded-xl p-2 text-slate-500 hover:bg-white/60 dark:text-slate-300 dark:hover:bg-white/5"
          :title="props.isDark ? t('home.switchToLight') : t('home.switchToDark')">
          <Icon v-if="props.isDark" name="sun" size="md" />
          <Icon v-else name="moon" size="md" />
        </button>
        <router-link v-if="props.isAuthenticated" :to="props.dashboardPath"
          class="inline-flex items-center gap-1.5 rounded-full bg-gradient-to-r from-sky-500 to-indigo-500 px-3 py-1 text-xs font-semibold text-white shadow-lg shadow-sky-500/30">
          <span>{{ props.userInitial }}</span>
          <span>{{ t('home.dashboard') }}</span>
        </router-link>
        <router-link v-else to="/login"
          class="inline-flex items-center rounded-full bg-gradient-to-r from-sky-500 to-indigo-500 px-4 py-1.5 text-xs font-semibold text-white shadow-lg shadow-sky-500/30">
          {{ t('home.login') }}
        </router-link>
      </div>
    </nav>
  </header>

  <main class="relative z-10 flex-1 px-6 pb-16 pt-24">
    <div class="mx-auto max-w-5xl">
      <!-- Hero: centered vertical with badge ring -->
      <div class="mb-16 text-center">
        <div class="mx-auto mb-8 inline-flex items-center gap-2 rounded-full border border-sky-200/80 bg-white/60 px-4 py-1.5 text-xs font-medium text-sky-700 backdrop-blur-xl dark:border-sky-900/50 dark:bg-sky-950/40 dark:text-sky-300">
          <span class="relative flex h-2 w-2">
            <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-sky-400 opacity-75"></span>
            <span class="relative inline-flex h-2 w-2 rounded-full bg-sky-500"></span>
          </span>
          Aurora Theme · Live Gateway
        </div>
        <h1 class="mb-6 bg-gradient-to-r from-sky-600 via-indigo-500 to-cyan-500 bg-clip-text text-5xl font-bold leading-tight text-transparent md:text-6xl lg:text-7xl">
          {{ props.heroTitle }}
        </h1>
        <p class="mx-auto mb-10 max-w-2xl text-lg text-slate-600 dark:text-slate-300 md:text-xl">
          {{ props.heroDescription }}
        </p>
        <router-link :to="props.isAuthenticated ? props.dashboardPath : '/login'"
          class="btn inline-flex items-center gap-2 rounded-full bg-gradient-to-r from-sky-500 via-indigo-500 to-cyan-500 px-10 py-3.5 text-base font-semibold text-white shadow-2xl shadow-sky-500/40 transition-transform hover:scale-105">
          {{ props.isAuthenticated ? t('home.goToDashboard') : props.ctaText }}
          <Icon name="arrowRight" size="md" />
        </router-link>
      </div>

      <!-- Feature Tags: aurora pills -->
      <div class="mb-16 grid gap-4 md:grid-cols-3">
        <div v-for="(tag, i) in props.featureTags" :key="i"
          class="group relative overflow-hidden rounded-2xl border border-white/60 bg-white/60 p-5 backdrop-blur-xl transition-all hover:-translate-y-1 hover:shadow-2xl dark:border-white/10 dark:bg-slate-900/60">
          <div class="absolute -right-8 -top-8 h-24 w-24 rounded-full bg-gradient-to-br from-sky-400/30 to-indigo-400/30 blur-2xl transition-opacity group-hover:opacity-80"></div>
          <div class="relative flex items-center gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-sky-500 to-indigo-500 shadow-lg shadow-sky-500/30">
              <Icon :name="tagIcons[i] ?? 'check'" size="md" class="text-white" />
            </div>
            <span class="font-semibold text-slate-700 dark:text-slate-100">{{ tag }}</span>
          </div>
        </div>
      </div>

      <div v-if="props.registrationNotice || props.allowSubSiteOpen" class="mb-16 grid gap-4 lg:grid-cols-2">
        <div v-if="props.registrationNotice" class="rounded-2xl border border-amber-200/60 bg-amber-50/80 px-5 py-4 text-sm text-amber-700 backdrop-blur-xl dark:border-amber-900/40 dark:bg-amber-950/30 dark:text-amber-300">
          {{ props.registrationNotice }}
        </div>
        <div v-if="props.allowSubSiteOpen" class="rounded-2xl border border-sky-200/60 bg-white/60 px-5 py-4 shadow-lg backdrop-blur-xl dark:border-sky-900/30 dark:bg-slate-900/60">
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p class="text-sm font-semibold text-slate-900 dark:text-white">当前分站支持自助开通下级分站</p>
              <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">售价 ￥{{ props.subSiteOpenPrice }}</p>
            </div>
            <router-link :to="props.isAuthenticated ? '/subsites' : '/login'"
              class="btn rounded-full bg-gradient-to-r from-sky-500 to-indigo-500 px-5 py-2 text-xs font-semibold text-white shadow-md">
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
