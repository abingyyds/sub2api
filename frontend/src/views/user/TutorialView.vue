<template>
  <AppLayout>
    <div class="mx-auto max-w-4xl space-y-8">
      <!-- Title -->
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.title') }}</h1>
        <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('tutorial.subtitle') }}</p>
      </div>

      <!-- API Integration Steps -->
      <div class="card p-6">
        <h2 class="mb-4 text-xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.integration.title') }}</h2>
        <div class="space-y-4">
          <div class="flex items-start gap-4">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">1</div>
            <div class="flex-1">
              <p class="font-medium text-gray-900 dark:text-white">{{ t('tutorial.integration.step1') }}</p>
              <p class="mt-1 text-sm text-gray-600 dark:text-dark-300">{{ t('tutorial.integration.step1Desc') }}</p>
            </div>
          </div>
          <div class="flex items-start gap-4">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">2</div>
            <div class="flex-1">
              <p class="font-medium text-gray-900 dark:text-white">{{ t('tutorial.integration.step2') }}</p>
              <p class="mt-1 text-sm text-gray-600 dark:text-dark-300">{{ t('tutorial.integration.step2Desc') }}</p>
            </div>
          </div>
          <div class="flex items-start gap-4">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">3</div>
            <div class="flex-1">
              <p class="font-medium text-gray-900 dark:text-white">{{ t('tutorial.integration.step3') }}</p>
              <p class="mt-1 text-sm text-gray-600 dark:text-dark-300">{{ t('tutorial.integration.step3Desc') }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Config Export Section -->
      <div class="card p-6">
        <h2 class="mb-2 text-xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.configExport.title') }}</h2>
        <p class="mb-6 text-sm text-gray-500 dark:text-dark-400">{{ t('tutorial.configExport.subtitle') }}</p>

        <!-- Selectors Row -->
        <div class="grid gap-4 md:grid-cols-2 mb-6">
          <!-- API Key Selector -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1.5">{{ t('tutorial.configExport.selectKey') }}</label>
            <select
              v-model="selectedKeyId"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 dark:border-dark-600 dark:bg-dark-800 dark:text-white focus:border-primary-500 focus:ring-1 focus:ring-primary-500"
            >
              <option v-if="loadingKeys" :value="null" disabled>{{ t('common.loading') }}...</option>
              <option v-else-if="apiKeys.length === 0" :value="null" disabled>{{ t('tutorial.configExport.noKeys') }}</option>
              <option
                v-for="key in apiKeys"
                :key="key.id"
                :value="key.id"
              >
                {{ key.name }} ({{ key.key.slice(0, 16) }}...)
              </option>
            </select>
          </div>

          <!-- Model Selector -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1.5">{{ t('tutorial.configExport.selectModel') }}</label>
            <select
              v-model="selectedModel"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 dark:border-dark-600 dark:bg-dark-800 dark:text-white focus:border-primary-500 focus:ring-1 focus:ring-primary-500"
            >
              <option v-for="model in popularModels" :key="model" :value="model">{{ model }}</option>
            </select>
          </div>
        </div>

        <!-- Tool Selection Tabs -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-2">{{ t('tutorial.configExport.selectTool') }}</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="tool in tools"
              :key="tool.id"
              @click="selectedTool = tool.id"
              class="rounded-lg px-3 py-1.5 text-sm font-medium transition-all"
              :class="selectedTool === tool.id
                ? 'bg-primary-600 text-white shadow-sm'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-dark-700 dark:text-dark-300 dark:hover:bg-dark-600'"
            >
              {{ tool.name }}
            </button>
          </div>
        </div>

        <!-- Generated Config Display -->
        <div v-if="selectedKeyId && generatedConfig" class="relative">
          <div class="bg-gray-900 dark:bg-dark-900 rounded-xl overflow-hidden">
            <!-- Code Header -->
            <div class="flex items-center justify-between px-4 py-2 bg-gray-800 dark:bg-dark-800 border-b border-gray-700 dark:border-dark-700">
              <span class="text-xs text-gray-400 font-mono">{{ configFilePath }}</span>
              <button
                @click="copyConfig"
                class="flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium rounded-lg transition-colors"
                :class="copied
                  ? 'bg-green-500/20 text-green-400'
                  : 'bg-gray-700 hover:bg-gray-600 text-gray-300 hover:text-white'"
              >
                <svg v-if="copied" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
                <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15.666 3.888A2.25 2.25 0 0013.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 01-.75.75H9a.75.75 0 01-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 01-2.25 2.25H6.75A2.25 2.25 0 014.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 011.927-.184" />
                </svg>
                {{ copied ? t('keys.useKeyModal.copied') : t('keys.useKeyModal.copy') }}
              </button>
            </div>
            <!-- Code Content -->
            <pre class="p-4 text-sm font-mono text-gray-100 overflow-x-auto"><code>{{ generatedConfig }}</code></pre>
          </div>
        </div>

        <!-- No key selected hint -->
        <div v-else-if="!selectedKeyId && apiKeys.length > 0" class="rounded-lg border-2 border-dashed border-gray-300 dark:border-dark-600 p-8 text-center">
          <svg class="mx-auto h-10 w-10 text-gray-400 dark:text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 5.25a3 3 0 013 3m3 0a6 6 0 01-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 17.25H8.25v2.25H6v2.25H2.25v-2.818c0-.597.237-1.17.659-1.591l6.499-6.499c.404-.404.527-1 .43-1.563A6 6 0 1121.75 8.25z" />
          </svg>
          <p class="mt-3 text-sm text-gray-500 dark:text-dark-400">{{ t('tutorial.configExport.selectKeyHint') }}</p>
        </div>

        <!-- No keys at all -->
        <div v-else-if="apiKeys.length === 0 && !loadingKeys" class="rounded-lg border-2 border-dashed border-gray-300 dark:border-dark-600 p-8 text-center">
          <svg class="mx-auto h-10 w-10 text-gray-400 dark:text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="mt-3 text-sm text-gray-500 dark:text-dark-400">{{ t('tutorial.configExport.noKeysHint') }}</p>
          <router-link to="/keys" class="mt-2 inline-block text-sm font-medium text-primary-600 hover:text-primary-500">
            {{ t('tutorial.configExport.createKey') }} &rarr;
          </router-link>
        </div>
      </div>

      <!-- Tips -->
      <div class="rounded-lg bg-blue-50 p-4 dark:bg-blue-900/20">
        <div class="flex items-start gap-3">
          <svg class="mt-0.5 h-5 w-5 flex-shrink-0 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z" />
          </svg>
          <div class="flex-1">
            <p class="font-medium text-blue-900 dark:text-blue-200">{{ t('tutorial.tips.title') }}</p>
            <ul class="mt-2 space-y-1 text-sm text-blue-800 dark:text-blue-300">
              <li>&bull; {{ t('tutorial.tips.tip1') }}</li>
              <li>&bull; {{ t('tutorial.tips.tip2') }}</li>
              <li>&bull; {{ t('tutorial.tips.tip3') }}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import { keysAPI } from '@/api/keys'
import { useClipboard } from '@/composables/useClipboard'
import type { ApiKey } from '@/types'

const { t } = useI18n()
const { copyToClipboard } = useClipboard()

const apiBaseUrl = computed(() => window.location.origin)
const selectedTool = ref('claudeCode')
const selectedKeyId = ref<number | null>(null)
const selectedModel = ref('claude-sonnet-4-6')
const apiKeys = ref<ApiKey[]>([])
const loadingKeys = ref(true)
const copied = ref(false)

const popularModels = [
  'claude-sonnet-4-6',
  'claude-opus-4-6',
  'claude-haiku-4-5-20251001',
  'claude-sonnet-4-5-20250929',
  'claude-opus-4-5-20251101',
  'claude-sonnet-4-20250514',
  'gpt-5.2-codex',
  'gpt-5.1-codex',
  'gpt-5.2',
  'gpt-5.1',
  'gemini-2.5-pro',
  'gemini-2.5-flash',
  'gemini-2.0-flash',
]

const tools = computed(() => [
  { id: 'claudeCode', name: 'Claude Code' },
  { id: 'openclaw', name: 'OpenClaw' },
  { id: 'opencode', name: 'OpenCode' },
  { id: 'cursor', name: 'Cursor' },
  { id: 'api', name: 'cURL' },
  { id: 'python', name: 'Python SDK' },
  { id: 'anthropic', name: 'Anthropic SDK' },
])

const selectedKey = computed(() => apiKeys.value.find(k => k.id === selectedKeyId.value))
const currentApiKey = computed(() => selectedKey.value?.key || 'sk-your-api-key')

const configFilePath = computed(() => {
  switch (selectedTool.value) {
    case 'claudeCode': return '~/.claude/settings.json'
    case 'openclaw': return '~/.openclaw/openclaw.json'
    case 'opencode': return '~/.config/opencode/opencode.json'
    case 'cursor': return 'Cursor Settings'
    case 'api': return 'Terminal'
    case 'python': return 'main.py'
    case 'anthropic': return 'main.py'
    default: return ''
  }
})

const generatedConfig = computed(() => {
  const base = apiBaseUrl.value
  const key = currentApiKey.value
  const model = selectedModel.value

  switch (selectedTool.value) {
    case 'claudeCode':
      return JSON.stringify({
        env: {
          ANTHROPIC_API_KEY: key,
          ANTHROPIC_BASE_URL: base + '/',
          ANTHROPIC_MODEL: model
        }
      }, null, 2)

    case 'openclaw':
      return JSON.stringify({
        provider: 'anthropic',
        base_url: base + '/',
        api: 'anthropic-messages',
        api_key: key,
        model: {
          id: model,
          name: model
        }
      }, null, 2)

    case 'opencode':
      return JSON.stringify({
        provider: {
          openai: {
            options: {
              baseURL: base + '/v1',
              apiKey: key
            },
            models: {
              [model]: {
                name: model,
                options: { store: false }
              }
            }
          }
        },
        $schema: 'https://opencode.ai/config.json'
      }, null, 2)

    case 'cursor':
      return `API Key: ${key}\nBase URL: ${base}/v1\nModel: ${model}`

    case 'api':
      return `curl ${base}/v1/messages \\
  -H "Content-Type: application/json" \\
  -H "x-api-key: ${key}" \\
  -H "anthropic-version: 2023-06-01" \\
  -d '{
    "model": "${model}",
    "max_tokens": 1024,
    "messages": [
      {"role": "user", "content": "Hello, Claude"}
    ]
  }'`

    case 'python':
      return `from openai import OpenAI

client = OpenAI(
    api_key="${key}",
    base_url="${base}/v1"
)

response = client.chat.completions.create(
    model="${model}",
    max_tokens=1024,
    messages=[
        {"role": "user", "content": "Hello"}
    ]
)
print(response.choices[0].message.content)`

    case 'anthropic':
      return `from anthropic import Anthropic

client = Anthropic(
    api_key="${key}",
    base_url="${base}"
)

message = client.messages.create(
    model="${model}",
    max_tokens=1024,
    messages=[
        {"role": "user", "content": "Hello, Claude"}
    ]
)
print(message.content)`

    default:
      return ''
  }
})

onMounted(async () => {
  try {
    const resp = await keysAPI.list(1, 100)
    apiKeys.value = (resp.data || []).filter((k: ApiKey) => k.status === 'active')
    if (apiKeys.value.length > 0) {
      selectedKeyId.value = apiKeys.value[0].id
    }
  } catch {
    // silently fail
  } finally {
    loadingKeys.value = false
  }
})

const copyConfig = async () => {
  if (!generatedConfig.value) return
  const success = await copyToClipboard(generatedConfig.value, t('keys.useKeyModal.copied'))
  if (success) {
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  }
}
</script>
