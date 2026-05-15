<template>
  <BaseDialog :show="show" :title="title" width="full" @close="handleClose">
    <div class="space-y-4">
      <div class="rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-sm text-gray-600 dark:border-dark-700 dark:bg-dark-900/60 dark:text-dark-300">
        {{ document.summary }}
      </div>

      <div
        ref="scrollRef"
        class="max-h-[65vh] overflow-y-auto rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-900"
        @scroll="handleScroll"
      >
        <LegalDocumentRenderer :document="document" />
      </div>

      <div class="rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-sm dark:border-dark-700 dark:bg-dark-900/60">
        <div class="flex items-center gap-2 text-gray-600 dark:text-dark-300">
          <Icon :name="accepted ? 'checkCircle' : 'document'" size="sm" :class="accepted ? 'text-emerald-500' : 'text-gray-400'" />
          <span>
            {{ accepted ? acceptedLabel : readPrompt }}
          </span>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex w-full flex-col gap-3 sm:flex-row sm:justify-end">
        <button type="button" class="btn btn-secondary w-full sm:w-auto" @click="handleReject">
          {{ rejectLabel }}
        </button>
        <button type="button" class="btn btn-primary w-full sm:w-auto" :disabled="!accepted" @click="handleAccept">
          {{ acceptLabel }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import LegalDocumentRenderer from './LegalDocumentRenderer.vue'
import type { LegalDocument } from '@/legal/documents'

const props = defineProps<{
  show: boolean
  document: LegalDocument
  acceptLabel: string
  rejectLabel: string
  readPrompt: string
  acceptedLabel: string
}>()

const emit = defineEmits<{
  (e: 'accept'): void
  (e: 'reject'): void
}>()

const scrollRef = ref<HTMLElement | null>(null)
const accepted = ref(false)
const atBottom = ref(false)

const title = computed(() => props.document.title)

function handleClose(): void {
  emit('reject')
}

function handleAccept(): void {
  if (!accepted.value) {
    return
  }
  emit('accept')
}

function handleReject(): void {
  emit('reject')
}

function handleScroll(): void {
  const el = scrollRef.value
  if (!el) return

  const threshold = 24
  atBottom.value = el.scrollTop + el.clientHeight >= el.scrollHeight - threshold
  accepted.value = atBottom.value
}

watch(
  () => props.show,
  async (showing) => {
    if (!showing) {
      accepted.value = false
      atBottom.value = false
      return
    }

    await nextTick()
    const el = scrollRef.value
    if (!el) return
    el.scrollTop = 0
    accepted.value = false
    atBottom.value = false
  }
)
</script>
