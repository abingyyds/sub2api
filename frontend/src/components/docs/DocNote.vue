<template>
  <div
    class="my-4 rounded-lg p-4"
    :class="typeClasses"
  >
    <div class="flex items-start gap-3">
      <!-- Info icon -->
      <svg v-if="type === 'info'" class="mt-0.5 h-5 w-5 flex-shrink-0" :class="iconClass" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z" />
      </svg>
      <!-- Warning icon -->
      <svg v-else-if="type === 'warning'" class="mt-0.5 h-5 w-5 flex-shrink-0" :class="iconClass" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
      </svg>
      <!-- Tip icon -->
      <svg v-else class="mt-0.5 h-5 w-5 flex-shrink-0" :class="iconClass" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 18v-5.25m0 0a6.01 6.01 0 001.5-.189m-1.5.189a6.01 6.01 0 01-1.5-.189m3.75 7.478a12.06 12.06 0 01-4.5 0m3.75 2.383a14.406 14.406 0 01-3 0M14.25 18v-.192c0-.983.658-1.823 1.508-2.316a7.5 7.5 0 10-7.517 0c.85.493 1.509 1.333 1.509 2.316V18" />
      </svg>
      <div class="flex-1 text-sm" :class="textClass">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  type?: 'info' | 'warning' | 'tip'
}>(), {
  type: 'info'
})

const typeClasses = computed(() => {
  switch (props.type) {
    case 'warning': return 'bg-yellow-50 dark:bg-yellow-900/20'
    case 'tip': return 'bg-green-50 dark:bg-green-900/20'
    default: return 'bg-blue-50 dark:bg-blue-900/20'
  }
})

const iconClass = computed(() => {
  switch (props.type) {
    case 'warning': return 'text-yellow-600 dark:text-yellow-400'
    case 'tip': return 'text-green-600 dark:text-green-400'
    default: return 'text-blue-600 dark:text-blue-400'
  }
})

const textClass = computed(() => {
  switch (props.type) {
    case 'warning': return 'text-yellow-800 dark:text-yellow-300'
    case 'tip': return 'text-green-800 dark:text-green-300'
    default: return 'text-blue-800 dark:text-blue-300'
  }
})
</script>
