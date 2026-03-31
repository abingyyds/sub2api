<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-6xl">
        <!-- Mobile doc nav toggle -->
        <div class="lg:hidden mb-4">
          <button
            @click="showMobileSidebar = !showMobileSidebar"
            class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 dark:text-dark-300 dark:bg-dark-700 dark:hover:bg-dark-600 transition-colors"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
            </svg>
            {{ t('tutorial.docs.navigation') }}
          </button>
          <!-- Mobile sidebar drawer -->
          <div v-if="showMobileSidebar" class="mt-2 rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-900">
            <DocSidebar :current-doc="currentDoc" @click="showMobileSidebar = false" />
          </div>
        </div>

        <div class="flex gap-8">
          <!-- Desktop sidebar -->
          <div class="hidden lg:block w-56 flex-shrink-0">
            <div class="sticky top-6">
              <div class="rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-900">
                <DocSidebar :current-doc="currentDoc" />
              </div>
            </div>
          </div>

          <!-- Content area -->
          <div class="flex-1 min-w-0">
            <SlideIn direction="up" :delay="100">
              <component :is="currentComponent" :key="currentDoc" />
            </SlideIn>
          </div>
        </div>
      </div>
    </FadeIn>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, ref, watch, type Component } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import { FadeIn, SlideIn } from '@/components/animations'
import DocSidebar from '@/components/docs/DocSidebar.vue'

// Article components
import DocsHome from '@/components/docs/articles/DocsHome.vue'
import QuickStartNodejs from '@/components/docs/articles/QuickStartNodejs.vue'
import QuickStartClaudeCode from '@/components/docs/articles/QuickStartClaudeCode.vue'
import QuickStartGeminiCli from '@/components/docs/articles/QuickStartGeminiCli.vue'
import QuickStartCodex from '@/components/docs/articles/QuickStartCodex.vue'
import AdvancedOpenClaw from '@/components/docs/articles/AdvancedOpenClaw.vue'
import AdvancedOpenCode from '@/components/docs/articles/AdvancedOpenCode.vue'

const { t } = useI18n()
const route = useRoute()
const showMobileSidebar = ref(false)

const docComponents: Record<string, Component> = {
  'home': DocsHome,
  'nodejs': QuickStartNodejs,
  'claude-code': QuickStartClaudeCode,
  'gemini-cli': QuickStartGeminiCli,
  'codex': QuickStartCodex,
  'openclaw': AdvancedOpenClaw,
  'opencode': AdvancedOpenCode,
}

const currentDoc = computed(() => (route.query.doc as string) || 'home')
const currentComponent = computed(() => docComponents[currentDoc.value] || DocsHome)

// Close mobile sidebar on navigation
watch(() => route.query.doc, () => {
  showMobileSidebar.value = false
})
</script>
