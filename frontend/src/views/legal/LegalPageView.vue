<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950">
    <!-- Navigation -->
    <nav class="sticky top-0 z-50 border-b border-gray-200/50 bg-white/80 backdrop-blur-lg dark:border-dark-800/50 dark:bg-dark-900/80">
      <div class="mx-auto max-w-5xl px-6 py-4">
        <div class="flex items-center justify-between">
          <router-link to="/home" class="flex items-center gap-3">
            <Logo :size="32" theme="auto" />
            <span class="text-lg font-bold text-gray-900 dark:text-white">cCoder.me</span>
          </router-link>
          <router-link to="/home" class="text-sm text-gray-500 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">
            &larr; Back to Home
          </router-link>
        </div>
      </div>
    </nav>

    <!-- Content -->
    <main class="mx-auto max-w-5xl px-6 py-12">
      <div class="mb-8 flex flex-wrap gap-3">
        <router-link
          to="/legal/terms"
          class="rounded-full border px-4 py-2 text-sm font-medium transition-colors"
          :class="document.kind === 'terms' ? activeTabClasses : inactiveTabClasses"
        >
          User Agreement
        </router-link>
        <router-link
          to="/legal/privacy"
          class="rounded-full border px-4 py-2 text-sm font-medium transition-colors"
          :class="document.kind === 'privacy' ? activeTabClasses : inactiveTabClasses"
        >
          Privacy Policy
        </router-link>
      </div>

      <LegalDocumentRenderer :document="document" />
    </main>

    <!-- Footer -->
    <footer class="border-t border-gray-200 px-6 py-8 dark:border-dark-800">
      <div class="mx-auto max-w-5xl">
        <div class="flex flex-col items-center gap-4">
          <div class="flex flex-wrap items-center justify-center gap-4 text-xs text-gray-400 dark:text-dark-500">
            <router-link to="/legal/terms" class="transition-colors hover:text-gray-600 dark:hover:text-dark-300">
              User Agreement
            </router-link>
            <span>|</span>
            <router-link to="/legal/privacy" class="transition-colors hover:text-gray-600 dark:hover:text-dark-300">
              Privacy Policy
            </router-link>
          </div>
          <p class="text-xs text-gray-400 dark:text-dark-500">
            &copy; {{ currentYear }} cCoder.me. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Logo from '@/components/Logo.vue'
import LegalDocumentRenderer from '@/components/legal/LegalDocumentRenderer.vue'
import { getLegalDocument, type LegalDocumentKind } from '@/legal/documents'

const props = defineProps<{
  kind: LegalDocumentKind
}>()

const currentYear = computed(() => new Date().getFullYear())
const document = computed(() => getLegalDocument(props.kind))

const activeTabClasses = 'border-primary-500 bg-primary-50 text-primary-700 dark:border-primary-500 dark:bg-primary-900/30 dark:text-primary-300'
const inactiveTabClasses = 'border-gray-200 bg-white text-gray-600 hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-300 dark:hover:border-primary-700 dark:hover:text-primary-300'
</script>
