<template>
  <!-- Terminal: full mono, CLI-inspired, green-on-black aesthetic for Claude Code ops -->
  <header class="relative z-20 border-b border-emerald-500/20 bg-zinc-950/90 px-6 py-3 font-mono backdrop-blur">
    <nav class="mx-auto flex max-w-6xl items-center justify-between">
      <div class="flex items-center gap-3 text-emerald-400">
        <div class="flex h-8 w-8 items-center justify-center overflow-hidden rounded bg-emerald-500/10 ring-1 ring-emerald-500/40">
          <img :src="props.siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
        </div>
        <span class="text-sm font-semibold tracking-wider">{{ props.siteName.toLowerCase() }}@prod:~$</span>
      </div>
      <div class="flex items-center gap-2 text-xs">
        <LocaleSwitcher />
        <a v-if="props.docUrl" :href="props.docUrl" target="_blank" rel="noopener noreferrer"
          class="rounded border border-emerald-500/30 px-2 py-1 text-emerald-300 hover:bg-emerald-500/10"
          :title="t('home.viewDocs')">--docs</a>
        <button @click="emit('toggle-theme')"
          class="rounded border border-emerald-500/30 px-2 py-1 text-emerald-300 hover:bg-emerald-500/10"
          :title="props.isDark ? t('home.switchToLight') : t('home.switchToDark')">
          {{ props.isDark ? '--light' : '--dark' }}
        </button>
        <router-link v-if="props.isAuthenticated" :to="props.dashboardPath"
          class="rounded bg-emerald-500 px-3 py-1 font-bold text-zinc-950 hover:bg-emerald-400">
          exec dashboard
        </router-link>
        <router-link v-else to="/login"
          class="rounded bg-emerald-500 px-3 py-1 font-bold text-zinc-950 hover:bg-emerald-400">
          ./login.sh
        </router-link>
      </div>
    </nav>
  </header>

  <main class="relative z-10 flex-1 px-6 py-12 font-mono">
    <div class="mx-auto max-w-5xl">
      <!-- Hero: CLI output style -->
      <div class="mb-12 rounded-lg border border-emerald-500/30 bg-zinc-950/80 shadow-2xl shadow-emerald-500/10">
        <div class="flex items-center gap-2 border-b border-emerald-500/20 px-4 py-2">
          <span class="h-3 w-3 rounded-full bg-red-500"></span>
          <span class="h-3 w-3 rounded-full bg-yellow-500"></span>
          <span class="h-3 w-3 rounded-full bg-emerald-500"></span>
          <span class="ml-2 text-xs text-zinc-500">~/gateway — zsh</span>
        </div>
        <div class="p-8 md:p-12">
          <div class="mb-4 text-sm text-emerald-500">$ whoami</div>
          <h1 class="mb-6 text-3xl font-bold leading-tight text-emerald-300 md:text-5xl">
            <span class="text-zinc-500">&gt; </span>{{ props.heroTitle }}
          </h1>
          <div class="mb-6 text-sm text-emerald-500">$ cat manifesto.txt</div>
          <p class="mb-10 whitespace-pre-wrap text-base text-zinc-300 md:text-lg">
            {{ props.heroDescription }}
          </p>
          <div class="text-sm text-emerald-500">$ start --interactive</div>
          <router-link :to="props.isAuthenticated ? props.dashboardPath : '/login'"
            class="btn mt-4 inline-flex items-center gap-2 rounded border-2 border-emerald-400 bg-emerald-500/10 px-6 py-2.5 text-base font-bold text-emerald-300 hover:bg-emerald-500 hover:text-zinc-950">
            <span>▶</span>
            {{ props.isAuthenticated ? t('home.goToDashboard') : props.ctaText }}
          </router-link>
        </div>
      </div>

      <!-- Feature tags: terminal boxes -->
      <div class="mb-12 grid gap-3 md:grid-cols-3">
        <div v-for="(tag, i) in props.featureTags" :key="i"
          class="rounded border border-emerald-500/30 bg-zinc-950/60 p-4">
          <div class="mb-1 text-xs text-zinc-500">[module:{{ String(i + 1).padStart(2, '0') }}]</div>
          <div class="text-sm font-bold text-emerald-300">✓ {{ tag }}</div>
        </div>
      </div>

      <div v-if="props.registrationNotice || props.allowSubSiteOpen" class="mb-12 space-y-3 text-sm">
        <div v-if="props.registrationNotice" class="rounded border border-amber-500/40 bg-amber-500/10 px-4 py-3 text-amber-300">
          <span class="mr-2 text-amber-500">[WARN]</span>{{ props.registrationNotice }}
        </div>
        <div v-if="props.allowSubSiteOpen" class="rounded border border-emerald-500/30 bg-zinc-950/60 px-4 py-3 text-emerald-300">
          <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <span class="mr-2 text-emerald-500">[INFO]</span>
              <span>subsite_open_enabled=true · price=¥{{ props.subSiteOpenPrice }}</span>
            </div>
            <router-link :to="props.isAuthenticated ? '/subsites' : '/login'"
              class="btn rounded border border-emerald-400 bg-emerald-500/10 px-3 py-1 text-xs font-bold text-emerald-300 hover:bg-emerald-500 hover:text-zinc-950">
              {{ props.isAuthenticated ? 'cd /subsites' : './login.sh' }}
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
import type { ThemeHomeProps, ThemeHomeEmits } from './types'

const props = defineProps<ThemeHomeProps>()
const emit = defineEmits<ThemeHomeEmits>()
const { t } = useI18n()
</script>
