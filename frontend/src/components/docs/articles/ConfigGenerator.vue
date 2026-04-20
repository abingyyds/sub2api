<template>
  <div class="space-y-8">
    <!-- Title -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.configExport.title') }}</h1>
      <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('tutorial.configExport.subtitle') }}</p>
    </div>

    <!-- API Integration Steps -->
    <GlowCard glow-color="rgb(59, 130, 246)">
      <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft p-6">
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
    </GlowCard>

    <!-- Config Export Section -->
    <GlowCard glow-color="rgb(168, 85, 247)">
      <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft p-6">
        <h2 class="mb-2 text-xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.tools.title') }}</h2>
        <p class="mb-6 text-sm text-gray-500 dark:text-dark-400">{{ t('tutorial.configExport.subtitle') }}</p>

        <!-- Selectors Row -->
        <div class="grid gap-4 md:grid-cols-2 mb-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1.5">{{ t('tutorial.configExport.selectKey') }}</label>
            <select
              v-model="selectedKeyId"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 dark:border-dark-600 dark:bg-dark-800 dark:text-white focus:border-primary-500 focus:ring-1 focus:ring-primary-500"
            >
              <option v-if="loadingKeys" :value="null" disabled>{{ t('common.loading') }}...</option>
              <option v-else-if="apiKeys.length === 0" :value="null" disabled>{{ t('tutorial.configExport.noKeys') }}</option>
              <option v-for="key in apiKeys" :key="key.id" :value="key.id">
                {{ key.name }} ({{ key.key.slice(0, 16) }}...)
              </option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1.5">{{ t('tutorial.configExport.selectModel') }}</label>
            <select
              v-model="selectedModel"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 dark:border-dark-600 dark:bg-dark-800 dark:text-white focus:border-primary-500 focus:ring-1 focus:ring-primary-500"
            >
              <option v-if="loadingModels" value="" disabled>{{ t('common.loading') }}...</option>
              <option v-for="model in selectableModels" :key="model" :value="model">{{ model }}</option>
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
          <div
            v-if="selectedTool === 'cherryStudio'"
            class="mb-4 rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm leading-6 text-amber-900 dark:border-amber-500/30 dark:bg-amber-500/10 dark:text-amber-100"
          >
            {{ t('tutorial.configExport.cherryStudioHint') }}
          </div>
          <div class="bg-gray-900 dark:bg-dark-900 rounded-xl overflow-hidden">
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
            <pre class="p-4 text-sm font-mono text-gray-100 overflow-x-auto"><code>{{ generatedConfig }}</code></pre>
          </div>

          <div class="mt-4 flex flex-wrap gap-3">
            <button
              v-if="canDownloadConfig"
              @click="downloadGeneratedConfig"
              class="inline-flex items-center rounded-lg border border-gray-200 px-3 py-2 text-sm font-medium text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-400"
            >
              {{ t('tutorial.configExport.downloadConfig') }}
            </button>
            <button
              v-if="canImportToCherryStudio"
              @click="importToCherryStudio"
              class="inline-flex items-center rounded-lg bg-rose-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-rose-700"
            >
              {{ t('tutorial.configExport.importToCherryStudio') }}
            </button>
            <button
              v-if="canImportToCcSwitch"
              @click="importToCcSwitch"
              class="inline-flex items-center rounded-lg bg-primary-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
            >
              {{ t('tutorial.configExport.importToCcSwitch') }}
            </button>
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
    </GlowCard>

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
            <li>&bull; {{ t('tutorial.tips.tip4') }}</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { keysAPI } from '@/api/keys'
import { getModelPlaza } from '@/api/model-plaza'
import { GlowCard } from '@/components/animations'
import { useClipboard } from '@/composables/useClipboard'
import { useAppStore } from '@/stores/app'
import type { ApiKey } from '@/types'

type ToolId =
  | 'claudeCode'
  | 'geminiCli'
  | 'codexCli'
  | 'cherryStudio'
  | 'openclaw'
  | 'opencode'
  | 'cursor'
  | 'api'
  | 'openaiApi'
  | 'geminiApi'
  | 'python'
  | 'anthropic'
  | 'geminiPython'

type ModelFamily = 'anthropic' | 'openai' | 'gemini' | 'unknown'

const { t } = useI18n()
const { copyToClipboard } = useClipboard()
const appStore = useAppStore()

const apiBaseUrl = computed(() => window.location.origin)
const selectedTool = ref<ToolId>('claudeCode')
const selectedKeyId = ref<number | null>(null)
const selectedModel = ref('')
const apiKeys = ref<ApiKey[]>([])
const loadingKeys = ref(true)
const loadingModels = ref(true)
const copied = ref(false)
const availableModels = ref<string[]>([])

const tools: Array<{ id: ToolId; name: string }> = [
  { id: 'claudeCode', name: 'Claude Code' },
  { id: 'geminiCli', name: 'Gemini CLI' },
  { id: 'codexCli', name: 'Codex CLI' },
  { id: 'cherryStudio', name: 'Cherry Studio' },
  { id: 'openclaw', name: 'OpenClaw' },
  { id: 'opencode', name: 'OpenCode' },
  { id: 'cursor', name: 'Cursor' },
  { id: 'api', name: 'Anthropic cURL' },
  { id: 'openaiApi', name: 'OpenAI cURL' },
  { id: 'geminiApi', name: 'Gemini cURL' },
  { id: 'python', name: 'Python (OpenAI)' },
  { id: 'anthropic', name: 'Python (Anthropic)' },
  { id: 'geminiPython', name: 'Python (Gemini)' },
]

const selectedKey = computed(() => apiKeys.value.find(k => k.id === selectedKeyId.value) ?? null)
const currentApiKey = computed(() => selectedKey.value?.key || 'sk-your-api-key')

function detectModelFamily(model: string): ModelFamily {
  const lowerModel = model.toLowerCase()
  if (lowerModel.includes('claude')) {
    return 'anthropic'
  }
  if (lowerModel.includes('gemini')) {
    return 'gemini'
  }
  if (['gpt', 'chatgpt', 'codex', 'o1', 'o3', 'o4'].some(keyword => lowerModel.includes(keyword))) {
    return 'openai'
  }
  return 'unknown'
}

function preferredFamiliesForTool(toolId: ToolId): ModelFamily[] | null {
  switch (toolId) {
    case 'claudeCode':
    case 'api':
    case 'anthropic':
      return ['anthropic']
    case 'geminiCli':
    case 'geminiApi':
    case 'geminiPython':
      return ['gemini']
    case 'codexCli':
    case 'cursor':
    case 'openaiApi':
    case 'python':
      return ['openai']
    default:
      return null
  }
}

function pickDefaultModel(models: string[], toolId: ToolId) {
  if (!models.length) {
    return ''
  }

  const preferredFamily = preferredFamiliesForTool(toolId)?.[0]

  if (preferredFamily === 'anthropic') {
    return models.find(model => model.includes('claude-sonnet')) || models.find(model => model.includes('claude')) || models[0]
  }
  if (preferredFamily === 'gemini') {
    return models.find(model => model.includes('gemini-2.5-pro')) || models.find(model => model.includes('gemini')) || models[0]
  }
  if (preferredFamily === 'openai') {
    return models.find(model => model.includes('gpt-5.3-codex')) || models.find(model => model.includes('gpt-5')) || models[0]
  }

  return models[0]
}

const selectableModels = computed(() => {
  const preferredFamilies = preferredFamiliesForTool(selectedTool.value)
  if (!preferredFamilies) {
    return availableModels.value
  }

  const filtered = availableModels.value.filter(model => preferredFamilies.includes(detectModelFamily(model)))
  return filtered.length ? filtered : availableModels.value
})

watch(
  [selectableModels, selectedTool],
  ([models, tool]) => {
    if (!models.length) {
      selectedModel.value = ''
      return
    }

    if (!models.includes(selectedModel.value)) {
      selectedModel.value = pickDefaultModel(models, tool)
    }
  },
  { immediate: true }
)

function getCherryStudioProvider(model: string, baseUrl: string) {
  const family = detectModelFamily(model)

  if (family === 'gemini') {
    return {
      providerType: 'Gemini',
      type: 'gemini',
      providerId: 'ccoder-me-gemini',
      providerName: 'cCoder.me (Gemini)',
      baseUrl: `${baseUrl}/v1beta`,
      docsEntry: '选择 Gemini 服务商或自定义 Gemini Provider',
    }
  }

  if (family === 'openai') {
    return {
      providerType: 'OpenAI',
      type: 'openai',
      providerId: 'ccoder-me-openai',
      providerName: 'cCoder.me (OpenAI)',
      baseUrl: `${baseUrl}/v1`,
      docsEntry: '选择 OpenAI 兼容服务商',
    }
  }

  return {
    providerType: 'Anthropic',
    type: 'anthropic',
    providerId: 'ccoder-me-anthropic',
    providerName: 'cCoder.me (Anthropic)',
    baseUrl,
    docsEntry: '选择 Anthropic 服务商',
  }
}

const configFilePath = computed(() => {
  switch (selectedTool.value) {
    case 'claudeCode': return '~/.claude/settings.json'
    case 'geminiCli': return '~/.gemini/.env + settings.json'
    case 'codexCli': return '~/.codex/config.toml + auth.json'
    case 'cherryStudio': return 'Cherry Studio -> 设置 -> 模型服务 -> 自定义服务商'
    case 'openclaw': return '~/.openclaw/openclaw.json'
    case 'opencode': return '~/.config/opencode/opencode.json'
    case 'cursor': return 'Cursor Settings'
    case 'api': return 'Terminal (Anthropic API)'
    case 'openaiApi': return 'Terminal (OpenAI API)'
    case 'geminiApi': return 'Terminal (Gemini API)'
    case 'python': return 'main.py'
    case 'anthropic': return 'main.py'
    case 'geminiPython': return 'main.py'
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
          ANTHROPIC_BASE_URL: `${base}/`,
          ANTHROPIC_MODEL: model
        }
      }, null, 2)

    case 'geminiCli':
      return `# 1. ${t('tutorial.geminiCli.step1')}
npm install -g @google/gemini-cli

# 2. ${t('tutorial.geminiCli.step2')} ~/.gemini/.env
GOOGLE_GEMINI_BASE_URL=${base}/v1beta/
GEMINI_API_KEY=${key}
GEMINI_MODEL=${model}

# 3. ${t('tutorial.geminiCli.step3')} ~/.gemini/settings.json
${JSON.stringify({
  ide: { enabled: true },
  security: {
    auth: { selectedType: 'gemini-api-key' }
  }
}, null, 2)}

# ${t('tutorial.geminiCli.authTip')}`

    case 'codexCli':
      return `# ~/.codex/config.toml
model_provider = "ccoder.me"
model = "${model}"
model_reasoning_effort = "high"
disable_response_storage = true
preferred_auth_method = "apikey"

[model_providers.ccoder.me]
name = "ccoder.me"
base_url = "${base}"
wire_api = "responses"
requires_openai_auth = true

# ~/.codex/auth.json
{
  "OPENAI_API_KEY": "${key}"
}`

    case 'cherryStudio': {
      const provider = getCherryStudioProvider(model, base)
      return `# Cherry Studio 自定义服务商参数
提供商名称：${provider.providerName}
服务商类型：${provider.providerType}
API 地址：${provider.baseUrl}
API Key：${key}
默认模型：${model}

# 配置步骤
1. 打开 Cherry Studio -> 设置 -> 模型服务
2. ${provider.docsEntry}
3. 填入上面的 API 地址与 API Key
4. 将模型 ${model} 添加到模型列表并保存`
    }

    case 'openclaw': {
      let provider = 'anthropic'
      let api = 'anthropic-messages'
      const lowerModel = model.toLowerCase()
      if (lowerModel.includes('gpt') || lowerModel.includes('o1') || lowerModel.includes('o3') || lowerModel.includes('codex')) {
        provider = 'openai'
        api = 'openai-chat'
      } else if (lowerModel.includes('gemini')) {
        provider = 'google'
        api = 'google-chat'
      }
      return JSON.stringify({
        provider,
        base_url: `${base}/`,
        api,
        api_key: key,
        model: { id: model, name: model }
      }, null, 2)
    }

    case 'opencode':
      return JSON.stringify({
        provider: {
          openai: {
            options: { baseURL: `${base}/v1`, apiKey: key },
            models: { [model]: { name: model, options: { store: false } } }
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

    case 'openaiApi':
      return `curl ${base}/v1/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer ${key}" \\
  -d '{
    "model": "${model}",
    "max_tokens": 1024,
    "messages": [
      {"role": "user", "content": "Hello"}
    ]
  }'`

    case 'geminiApi':
      return `curl "${base}/v1beta/models/${model}:generateContent?key=${key}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "contents": [
      {"parts": [{"text": "Hello, Gemini"}]}
    ]
  }'

# Streaming:
curl "${base}/v1beta/models/${model}:streamGenerateContent?alt=sse&key=${key}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "contents": [
      {"parts": [{"text": "Hello, Gemini"}]}
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

    case 'geminiPython':
      return `from google import genai

client = genai.Client(
    api_key="${key}",
    http_options={"api_version": "v1beta", "base_url": "${base}"}
)

response = client.models.generate_content(
    model="${model}",
    contents="Hello, Gemini"
)
print(response.text)

# Image Generation:
# response = client.models.generate_content(
#     model="gemini-3.1-flash-image",
#     contents="Generate a cute cat illustration",
#     config={"response_modalities": ["IMAGE", "TEXT"]}
# )`

    default:
      return ''
  }
})

const canDownloadConfig = computed(() => Boolean(selectedKeyId.value && generatedConfig.value))
const canImportToCherryStudio = computed(() => selectedTool.value === 'cherryStudio' && Boolean(selectedKey.value && selectedModel.value))

const canImportToCcSwitch = computed(() => {
  const platform = selectedKey.value?.group?.platform
  if (!platform) {
    return false
  }

  if (selectedTool.value === 'claudeCode') {
    return platform === 'anthropic' || platform === 'antigravity'
  }
  if (selectedTool.value === 'geminiCli') {
    return platform === 'gemini' || platform === 'antigravity'
  }
  if (selectedTool.value === 'codexCli') {
    return platform === 'openai'
  }
  return false
})

function getDownloadFileName(toolId: ToolId) {
  switch (toolId) {
    case 'claudeCode': return 'claude-settings.json'
    case 'geminiCli': return 'gemini-cli-config.txt'
    case 'codexCli': return 'codex-config.txt'
    case 'cherryStudio': return 'cherry-studio-provider.txt'
    case 'openclaw': return 'openclaw.json'
    case 'opencode': return 'opencode.json'
    case 'cursor': return 'cursor-settings.txt'
    case 'api':
    case 'openaiApi':
    case 'geminiApi':
      return 'request-example.sh'
    case 'python':
    case 'anthropic':
    case 'geminiPython':
      return 'main.py'
    default:
      return 'config.txt'
  }
}

function downloadGeneratedConfig() {
  if (!generatedConfig.value) {
    return
  }

  const blob = new Blob([generatedConfig.value], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = getDownloadFileName(selectedTool.value)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

function getCherryStudioImportUrl() {
  if (!selectedKey.value || !selectedModel.value) {
    return null
  }

  const provider = getCherryStudioProvider(selectedModel.value, apiBaseUrl.value)
  const payload = JSON.stringify({
    id: provider.providerId,
    name: provider.providerName,
    type: provider.type,
    apiKey: selectedKey.value.key,
    baseUrl: provider.baseUrl,
  })

  const data = btoa(payload)
    .replace(/\+/g, '_')
    .replace(/\//g, '-')

  const params = new URLSearchParams({
    v: '1',
    data,
  })

  return `cherrystudio://providers/api-keys?${params.toString()}`
}

function importToCherryStudio() {
  const importUrl = getCherryStudioImportUrl()
  if (!importUrl) {
    return
  }

  try {
    window.open(importUrl, '_self')
    setTimeout(() => {
      if (document.hasFocus()) {
        appStore.showError(t('tutorial.configExport.cherryStudioNotInstalled'))
      }
    }, 120)
  } catch {
    appStore.showError(t('tutorial.configExport.cherryStudioNotInstalled'))
  }
}

function importToCcSwitch() {
  if (!canImportToCcSwitch.value || !selectedKey.value) {
    return
  }

  const baseUrl = apiBaseUrl.value
  const platform = selectedKey.value.group?.platform
  const app =
    selectedTool.value === 'codexCli'
      ? 'codex'
      : selectedTool.value === 'geminiCli'
        ? 'gemini'
        : 'claude'

  const endpoint = platform === 'antigravity' ? `${baseUrl}/antigravity` : baseUrl
  const usageScript = `({
    request: {
      url: "{{baseUrl}}/v1/usage",
      method: "GET",
      headers: { "Authorization": "Bearer {{apiKey}}" }
    },
    extractor: function(response) {
      return {
        isValid: response.is_active || true,
        remaining: response.balance,
        unit: "USD"
      };
    }
  })`

  const params = new URLSearchParams({
    resource: 'provider',
    app,
    name: 'ccoder.me',
    homepage: baseUrl,
    endpoint,
    apiKey: selectedKey.value.key,
    configFormat: 'json',
    usageEnabled: 'true',
    usageScript: btoa(usageScript),
    usageAutoInterval: '30'
  })

  try {
    window.open(`ccswitch://v1/import?${params.toString()}`, '_self')
    setTimeout(() => {
      if (document.hasFocus()) {
        appStore.showError(t('keys.ccSwitchNotInstalled'))
      }
    }, 120)
  } catch {
    appStore.showError(t('keys.ccSwitchNotInstalled'))
  }
}

onMounted(async () => {
  const keysPromise = keysAPI.list(1, 100)
    .then(resp => {
      apiKeys.value = (resp.items || []).filter((key: ApiKey) => key.status === 'active')
      if (apiKeys.value.length > 0) {
        selectedKeyId.value = apiKeys.value[0].id
      }
    })
    .catch(() => {})
    .finally(() => {
      loadingKeys.value = false
    })

  const modelsPromise = getModelPlaza()
    .then(groups => {
      const modelSet = new Set<string>()
      for (const group of groups) {
        const models = Array.isArray(group.models) ? group.models : []
        for (const model of models) {
          modelSet.add(model)
        }
      }
      availableModels.value = Array.from(modelSet).sort()
      if (!selectedModel.value) {
        selectedModel.value = pickDefaultModel(availableModels.value, selectedTool.value)
      }
    })
    .catch(() => {})
    .finally(() => {
      loadingModels.value = false
    })

  await Promise.all([keysPromise, modelsPromise])
})

const copyConfig = async () => {
  if (!generatedConfig.value) {
    return
  }

  const success = await copyToClipboard(generatedConfig.value, t('keys.useKeyModal.copied'))
  if (success) {
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  }
}
</script>
