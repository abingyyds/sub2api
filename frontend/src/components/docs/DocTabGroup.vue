<template>
  <div class="my-4">
    <!-- Tab Headers -->
    <div class="flex border-b border-gray-200 dark:border-dark-700">
      <button
        v-for="tab in normalizedTabs"
        :key="tab.id"
        @click="activeTab = tab.id"
        class="px-4 py-2 text-sm font-medium border-b-2 transition-colors -mb-px"
        :class="activeTab === tab.id
          ? 'border-primary-500 text-primary-600 dark:text-primary-400'
          : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-dark-400 dark:hover:text-dark-200'"
      >
        {{ tab.label }}
      </button>
    </div>
    <!-- Tab Content -->
    <div class="mt-3">
      <template v-for="tab in normalizedTabs" :key="tab.id">
        <div v-show="activeTab === tab.id">
          <slot :name="tab.id" />
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  tabs: (string | { id: string; label: string })[]
  defaultTab?: string
}>()

const normalizedTabs = computed(() =>
  props.tabs.map(tab =>
    typeof tab === 'string' ? { id: tab, label: tab } : tab
  )
)

const activeTab = ref(props.defaultTab || normalizedTabs.value[0]?.id)
</script>
