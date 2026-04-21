<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.docs.home') }}</h1>
      <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('tutorial.subtitle') }}</p>
    </div>

    <!-- Full Docs Banner -->
    <a
      href="https://docs.airiver.cn/"
      target="_blank"
      rel="noopener noreferrer"
      class="block rounded-xl bg-gradient-to-r from-primary-500 to-primary-600 p-6 text-white shadow-lg transition-all hover:shadow-xl hover:scale-[1.02] dark:from-primary-600 dark:to-primary-700"
    >
      <div class="flex items-center gap-4">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-lg bg-white/20 backdrop-blur-sm">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25" />
          </svg>
        </div>
        <div class="flex-1">
          <div class="flex items-center gap-2">
            <h3 class="text-lg font-semibold">完整使用教程与文档</h3>
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
            </svg>
          </div>
          <p class="mt-1 text-sm text-white/90">访问 docs.airiver.cn 查看更全面的教程、API 文档和最佳实践</p>
        </div>
      </div>
    </a>

    <section>
      <div class="flex flex-col gap-2 md:flex-row md:items-end md:justify-between">
        <div>
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-1">{{ t('tutorial.docs.downloads') }}</h2>
          <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('tutorial.docs.downloadsDesc') }}</p>
        </div>
      </div>
      <div class="mt-4 grid gap-4 lg:grid-cols-3">
        <div
          v-for="card in downloadCards"
          :key="card.title"
          class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm transition-all hover:-translate-y-0.5 hover:shadow-md dark:border-dark-700 dark:bg-dark-900"
        >
          <div class="flex items-start gap-3">
            <div class="flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-xl bg-primary-50 text-primary-600 dark:bg-primary-900/20 dark:text-primary-400">
              <component :is="card.icon" class="h-5 w-5" />
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <h3 class="text-base font-semibold text-gray-900 dark:text-white">{{ card.title }}</h3>
                <span class="rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600 dark:bg-dark-800 dark:text-dark-300">{{ card.version }}</span>
              </div>
              <p class="mt-1 text-sm leading-6 text-gray-500 dark:text-dark-400">{{ card.desc }}</p>
            </div>
          </div>

          <div class="mt-4 flex flex-wrap gap-2">
            <a
              v-for="action in card.metaLinks"
              :key="action.href"
              :href="action.href"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center rounded-lg border border-gray-200 px-3 py-2 text-sm font-medium text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-400"
            >
              {{ action.label }}
            </a>
          </div>

          <div class="mt-4 space-y-3">
            <div
              v-for="group in card.groups"
              :key="group.title"
              class="rounded-xl bg-gray-50 p-3 dark:bg-dark-800/70"
            >
              <p class="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">{{ group.title }}</p>
              <div class="mt-2 flex flex-wrap gap-2">
                <a
                  v-for="link in group.links"
                  :key="link.href"
                  :href="link.href"
                  rel="noopener noreferrer"
                  class="inline-flex items-center rounded-lg border px-3 py-2 text-sm font-medium transition-colors"
                  :class="link.recommended
                    ? 'border-primary-200 bg-primary-50 text-primary-700 hover:border-primary-300 dark:border-primary-800 dark:bg-primary-900/20 dark:text-primary-300'
                    : 'border-gray-200 text-gray-700 hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-400'"
                >
                  {{ link.label }}
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Advanced -->
    <section>
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">{{ t('tutorial.docs.advanced') }}</h2>
      <div class="grid gap-4 md:grid-cols-2">
        <router-link
          v-for="doc in advancedDocs"
          :key="doc.slug"
          :to="{ path: '/tutorial', query: { doc: doc.slug } }"
          class="group rounded-xl border-2 border-gray-200 bg-white p-5 transition-all hover:border-primary-300 hover:shadow-md dark:border-dark-700 dark:bg-dark-900 dark:hover:border-primary-600"
        >
          <div class="flex items-start gap-3">
            <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-purple-50 text-purple-600 dark:bg-purple-900/20 dark:text-purple-400">
              <component :is="doc.icon" class="h-5 w-5" />
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="font-medium text-gray-900 group-hover:text-primary-600 dark:text-white dark:group-hover:text-primary-400 transition-colors">{{ doc.title }}</h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400 line-clamp-2">{{ doc.desc }}</p>
            </div>
          </div>
        </router-link>
      </div>
    </section>

    <!-- Config Generator (inline) -->
    <section>
      <ConfigGenerator />
    </section>

  </div>
</template>

<script setup lang="ts">
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'
import ConfigGenerator from './ConfigGenerator.vue'
import { tutorialDownloadTools } from '../downloads'

const { t, locale } = useI18n()

// Inline icons
const BotIcon = () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714c0 .597.237 1.17.659 1.591L19.8 15.3M14.25 3.104c.251.023.501.05.75.082M19.8 15.3l-1.57.393A9.065 9.065 0 0112 15a9.065 9.065 0 00-6.23.693L5 14.5m14.8.8l1.402 1.402c1.232 1.232.65 3.318-1.067 3.611A48.309 48.309 0 0112 21c-2.773 0-5.491-.235-8.135-.687-1.718-.293-2.3-2.379-1.067-3.61L5 14.5' })
])

const SparkIcon = () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 00-2.455 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z' })
])

const CodeIcon = () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M17.25 6.75L22.5 12l-5.25 5.25m-10.5 0L1.5 12l5.25-5.25m7.5-3l-4.5 16.5' })
])

const RocketIcon = () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M15.59 14.37a6 6 0 01-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 006.16-12.12A14.98 14.98 0 009.631 8.41m5.96 5.96a14.926 14.926 0 01-5.841 2.58m-.119-8.54a6 6 0 00-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 00-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 01-2.448-2.448 14.9 14.9 0 01.06-.312m-2.24 2.39a4.493 4.493 0 00-1.757 4.306 4.493 4.493 0 004.306-1.758M16.5 9a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0z' })
])

const iconMap = {
  'CC Switch': SparkIcon,
  Codex: CodeIcon,
  'Cherry Studio': BotIcon,
} as const

type DownloadToolTitle = keyof typeof iconMap

const downloadCards = computed(() =>
  tutorialDownloadTools.map((tool) => {
    const metaLinks = [
      tool.officialRepo ? { label: t('tutorial.downloads.officialRepo'), href: tool.officialRepo } : null,
      tool.installGuide ? { label: t('tutorial.downloads.installGuide'), href: tool.installGuide } : null,
      tool.providerDocs ? { label: t('tutorial.downloads.providerDocs'), href: tool.providerDocs } : null,
      tool.releases ? { label: t('tutorial.downloads.releases'), href: tool.releases } : null,
    ].filter((item): item is { label: string; href: string } => Boolean(item))

    return {
      title: tool.title,
      version: tool.version,
      desc: locale.value.startsWith('zh') ? tool.descZh : tool.descEn,
      icon: iconMap[tool.title as DownloadToolTitle] ?? BotIcon,
      metaLinks,
      groups: tool.groups,
    }
  })
)

const advancedDocs = [
  { slug: 'openclaw', title: 'OpenClaw 部署教程', desc: '从零开始部署 OpenClaw Telegram 机器人，并接入 ccoder.me 平台。', icon: RocketIcon },
  { slug: 'opencode', title: 'OpenCode 配置教程', desc: '通过配置文件将 OpenCode 接入 Ai Go Code 平台，快速完成模型配置。', icon: BotIcon },
]
</script>
