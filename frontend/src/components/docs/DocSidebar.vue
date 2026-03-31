<template>
  <nav class="space-y-6">
    <div v-for="category in categories" :key="category.titleKey">
      <h3 class="px-3 text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-dark-400">
        {{ t(category.titleKey) }}
      </h3>
      <ul class="mt-2 space-y-0.5">
        <li v-for="item in category.items" :key="item.slug">
          <router-link
            :to="{ path: '/tutorial', query: { doc: item.slug } }"
            class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
            :class="currentDoc === item.slug
              ? 'bg-primary-50 text-primary-700 font-medium dark:bg-primary-900/20 dark:text-primary-400'
              : 'text-gray-700 hover:bg-gray-100 dark:text-dark-300 dark:hover:bg-dark-700'"
          >
            {{ t(item.titleKey) }}
          </router-link>
        </li>
      </ul>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps<{
  currentDoc: string
}>()

const categories = [
  {
    titleKey: 'tutorial.docs.quickStart',
    items: [
      { slug: 'nodejs', titleKey: 'tutorial.docs.nodejs' },
      { slug: 'claude-code', titleKey: 'tutorial.docs.claudeCode' },
      { slug: 'gemini-cli', titleKey: 'tutorial.docs.geminiCli' },
      { slug: 'codex', titleKey: 'tutorial.docs.codex' },
    ]
  },
  {
    titleKey: 'tutorial.docs.advanced',
    items: [
      { slug: 'openclaw', titleKey: 'tutorial.docs.openclaw' },
      { slug: 'opencode', titleKey: 'tutorial.docs.opencode' },
    ]
  },
]
</script>
