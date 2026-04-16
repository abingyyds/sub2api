<template>
  <!-- Summit: editorial/magazine layout with serif hero, warm tones -->
  <header class="relative z-20 border-b border-stone-300/60 bg-stone-50/80 px-6 py-4 backdrop-blur-md dark:border-stone-700/40 dark:bg-stone-900/60">
    <nav class="mx-auto flex max-w-6xl items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="h-10 w-10 overflow-hidden rounded-full ring-2 ring-orange-500/30">
          <img :src="props.siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
        </div>
        <div>
          <div class="font-serif text-lg font-bold text-stone-800 dark:text-stone-100">{{ props.siteName }}</div>
          <div class="text-[10px] uppercase tracking-[0.2em] text-stone-500 dark:text-stone-400">SUMMIT · ISSUE {{ props.currentYear }}</div>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <LocaleSwitcher />
        <a v-if="props.docUrl" :href="props.docUrl" target="_blank" rel="noopener noreferrer"
          class="rounded p-2 text-stone-500 hover:bg-stone-200/70 dark:text-stone-300 dark:hover:bg-stone-800"
          :title="t('home.viewDocs')"><Icon name="book" size="md" /></a>
        <button @click="emit('toggle-theme')"
          class="rounded p-2 text-stone-500 hover:bg-stone-200/70 dark:text-stone-300 dark:hover:bg-stone-800"
          :title="props.isDark ? t('home.switchToLight') : t('home.switchToDark')">
          <Icon v-if="props.isDark" name="sun" size="md" />
          <Icon v-else name="moon" size="md" />
        </button>
        <router-link v-if="props.isAuthenticated" :to="props.dashboardPath"
          class="inline-flex items-center gap-1.5 border-2 border-stone-900 px-3 py-1 text-xs font-bold uppercase tracking-wider text-stone-900 transition-colors hover:bg-stone-900 hover:text-stone-50 dark:border-stone-100 dark:text-stone-100 dark:hover:bg-stone-100 dark:hover:text-stone-900">
          {{ t('home.dashboard') }}
        </router-link>
        <router-link v-else to="/login"
          class="inline-flex items-center border-2 border-stone-900 px-3 py-1 text-xs font-bold uppercase tracking-wider text-stone-900 transition-colors hover:bg-stone-900 hover:text-stone-50 dark:border-stone-100 dark:text-stone-100 dark:hover:bg-stone-100 dark:hover:text-stone-900">
          {{ t('home.login') }}
        </router-link>
      </div>
    </nav>
  </header>

  <main class="relative z-10 flex-1 px-6 py-16">
    <div class="mx-auto max-w-5xl">
      <!-- Hero: editorial column layout -->
      <article class="mb-16 border-b-4 border-double border-stone-400/60 pb-16 dark:border-stone-600/60">
        <p class="mb-4 text-xs font-bold uppercase tracking-[0.35em] text-orange-700 dark:text-orange-400">— FEATURE —</p>
        <h1 class="mb-6 font-serif text-5xl font-black leading-[1.05] text-stone-900 dark:text-stone-50 md:text-6xl lg:text-7xl">
          {{ props.heroTitle }}
        </h1>
        <p class="mb-10 max-w-3xl border-l-4 border-orange-500 pl-6 font-serif text-xl italic text-stone-700 dark:text-stone-300">
          {{ props.heroDescription }}
        </p>
        <router-link :to="props.isAuthenticated ? props.dashboardPath : '/login'"
          class="btn inline-flex items-center gap-2 bg-stone-900 px-8 py-3 font-bold uppercase tracking-widest text-stone-50 shadow-md transition-all hover:bg-orange-600 dark:bg-stone-100 dark:text-stone-900 dark:hover:bg-orange-500 dark:hover:text-stone-50">
          {{ props.isAuthenticated ? t('home.goToDashboard') : props.ctaText }}
          <Icon name="arrowRight" size="md" />
        </router-link>
      </article>

      <!-- Feature tags: numbered columns -->
      <div class="mb-16 grid gap-8 md:grid-cols-3">
        <div v-for="(tag, i) in props.featureTags" :key="i" class="border-t-2 border-stone-900 pt-4 dark:border-stone-100">
          <div class="mb-2 font-serif text-3xl font-black text-orange-600 dark:text-orange-400">{{ String(i + 1).padStart(2, '0') }}</div>
          <div class="font-serif text-lg font-bold text-stone-900 dark:text-stone-100">{{ tag }}</div>
        </div>
      </div>

      <div v-if="props.registrationNotice || props.allowSubSiteOpen" class="mb-16 space-y-4">
        <div v-if="props.registrationNotice" class="border-l-4 border-amber-500 bg-amber-50 px-6 py-4 font-serif text-amber-900 dark:bg-amber-950/40 dark:text-amber-200">
          {{ props.registrationNotice }}
        </div>
        <div v-if="props.allowSubSiteOpen" class="border-y-2 border-stone-900 bg-stone-100 px-6 py-6 dark:border-stone-100 dark:bg-stone-800/60">
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p class="font-serif text-lg font-bold text-stone-900 dark:text-stone-50">当前分站支持自助开通下级分站</p>
              <p class="mt-1 text-sm text-stone-600 dark:text-stone-300">售价 ￥{{ props.subSiteOpenPrice }}</p>
            </div>
            <router-link :to="props.isAuthenticated ? '/subsites' : '/login'"
              class="btn border-2 border-stone-900 bg-transparent px-5 py-2 font-bold uppercase tracking-wider text-stone-900 hover:bg-stone-900 hover:text-stone-50 dark:border-stone-100 dark:text-stone-100 dark:hover:bg-stone-100 dark:hover:text-stone-900">
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
</script>
