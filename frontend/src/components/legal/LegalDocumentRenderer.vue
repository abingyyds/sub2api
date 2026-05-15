<template>
  <div class="space-y-8">
    <section>
      <div class="rounded-2xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-800">
        <p class="text-sm font-medium uppercase tracking-wide text-primary-600 dark:text-primary-400">
          {{ document.eyebrow }}
        </p>
        <h1 class="mt-2 text-3xl font-bold text-gray-900 dark:text-white">
          {{ document.title }}
        </h1>
        <p class="mt-4 text-gray-700 dark:text-dark-300">
          {{ document.summary }}
        </p>
        <p class="mt-3 text-sm text-gray-500 dark:text-dark-400">
          {{ document.readingNote }}
        </p>
        <p class="mt-3 text-xs text-gray-400 dark:text-dark-500">
          Last updated: {{ document.lastUpdated }} &middot; Version {{ document.version }}
        </p>
      </div>
    </section>

    <section
      v-for="section in document.sections"
      :key="section.title"
      class="rounded-2xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-800"
    >
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
        {{ section.title }}
      </h2>

      <div v-if="section.paragraphs?.length" class="mt-4 space-y-3 text-sm leading-7 text-gray-700 dark:text-dark-300">
        <p v-for="paragraph in section.paragraphs" :key="paragraph">
          {{ paragraph }}
        </p>
      </div>

      <ul
        v-if="section.bullets?.length"
        class="mt-4 list-disc space-y-2 pl-5 text-sm leading-7 text-gray-700 dark:text-dark-300"
      >
        <li v-for="bullet in section.bullets" :key="bullet">
          {{ bullet }}
        </li>
      </ul>

      <div v-if="section.subsections?.length" class="mt-6 space-y-6">
        <article
          v-for="subsection in section.subsections"
          :key="subsection.title"
          class="rounded-xl border border-gray-100 bg-gray-50 p-5 dark:border-dark-700 dark:bg-dark-900/60"
        >
          <h3 class="text-base font-semibold text-gray-900 dark:text-white">
            {{ subsection.title }}
          </h3>

          <div
            v-if="subsection.paragraphs?.length"
            class="mt-3 space-y-3 text-sm leading-7 text-gray-700 dark:text-dark-300"
          >
            <p v-for="paragraph in subsection.paragraphs" :key="paragraph">
              {{ paragraph }}
            </p>
          </div>

          <ul
            v-if="subsection.bullets?.length"
            class="mt-3 list-disc space-y-2 pl-5 text-sm leading-7 text-gray-700 dark:text-dark-300"
          >
            <li v-for="bullet in subsection.bullets" :key="bullet">
              {{ bullet }}
            </li>
          </ul>
        </article>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { LegalDocument } from '@/legal/documents'

defineProps<{
  document: LegalDocument
}>()
</script>
