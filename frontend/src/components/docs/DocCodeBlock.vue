<template>
  <div class="relative my-4">
    <div class="bg-gray-900 dark:bg-dark-900 rounded-xl overflow-hidden">
      <div v-if="filename || language" class="flex items-center justify-between px-4 py-2 bg-gray-800 dark:bg-dark-800 border-b border-gray-700 dark:border-dark-700">
        <span class="text-xs text-gray-400 font-mono">{{ filename || language }}</span>
        <button
          @click="handleCopy"
          class="flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium rounded-lg transition-colors"
          :class="isCopied
            ? 'bg-green-500/20 text-green-400'
            : 'bg-gray-700 hover:bg-gray-600 text-gray-300 hover:text-white'"
        >
          <svg v-if="isCopied" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>
          <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.666 3.888A2.25 2.25 0 0013.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 01-.75.75H9a.75.75 0 01-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 01-2.25 2.25H6.75A2.25 2.25 0 014.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 011.927-.184" />
          </svg>
          {{ isCopied ? '已复制' : '复制' }}
        </button>
      </div>
      <pre class="p-4 text-sm font-mono text-gray-100 overflow-x-auto leading-relaxed"><code>{{ code }}</code></pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useClipboard } from '@/composables/useClipboard'

const props = defineProps<{
  code: string
  language?: string
  filename?: string
}>()

const { copyToClipboard } = useClipboard()
const isCopied = ref(false)

const handleCopy = async () => {
  const success = await copyToClipboard(props.code, '已复制')
  if (success) {
    isCopied.value = true
    setTimeout(() => { isCopied.value = false }, 2000)
  }
}
</script>
