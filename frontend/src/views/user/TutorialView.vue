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

      <!-- Claude Code -->
      <div class="card p-6">
        <h2 class="mb-4 text-xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.tools.claudeCode.title') }}</h2>
        <p class="mb-4 text-sm text-gray-600 dark:text-dark-300">{{ t('tutorial.tools.claudeCode.desc') }}</p>
        <div>
          <h3 class="mb-2 text-sm font-medium text-gray-700 dark:text-dark-300">{{ t('tutorial.tools.claudeCode.file') }}</h3>
          <div class="rounded-lg bg-gray-100 p-3 font-mono text-sm dark:bg-dark-700">
            <pre class="overflow-x-auto text-gray-800 dark:text-dark-200">{{ claudeCodeConfig }}</pre>
          </div>
        </div>
      </div>

      <!-- OpenClaw -->
      <div class="card p-6">
        <h2 class="mb-4 text-xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.tools.openclaw.title') }}</h2>
        <p class="mb-4 text-sm text-gray-600 dark:text-dark-300">{{ t('tutorial.tools.openclaw.desc') }}</p>
        <div>
          <h3 class="mb-2 text-sm font-medium text-gray-700 dark:text-dark-300">{{ t('tutorial.tools.openclaw.file') }}</h3>
          <div class="rounded-lg bg-gray-100 p-3 font-mono text-sm dark:bg-dark-700">
            <pre class="overflow-x-auto text-gray-800 dark:text-dark-200">{{ openclawConfig }}</pre>
          </div>
        </div>
      </div>

      <!-- OpenCode -->
      <div class="card p-6">
        <h2 class="mb-4 text-xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.tools.opencode.title') }}</h2>
        <p class="mb-4 text-sm text-gray-600 dark:text-dark-300">{{ t('tutorial.tools.opencode.desc') }}</p>
        <div>
          <h3 class="mb-2 text-sm font-medium text-gray-700 dark:text-dark-300">{{ t('tutorial.tools.opencode.file') }}</h3>
          <div class="rounded-lg bg-gray-100 p-3 font-mono text-sm dark:bg-dark-700">
            <pre class="overflow-x-auto text-gray-800 dark:text-dark-200">{{ opencodeConfig }}</pre>
          </div>
        </div>
      </div>

      <!-- Cursor -->
      <div class="card p-6">
        <h2 class="mb-4 text-xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.tools.cursor.title') }}</h2>
        <p class="mb-4 text-sm text-gray-600 dark:text-dark-300">{{ t('tutorial.tools.cursor.desc') }}</p>
        <div class="space-y-2">
          <div class="rounded-lg bg-gray-100 p-3 dark:bg-dark-700">
            <p class="text-sm text-gray-600 dark:text-dark-300">{{ t('tutorial.tools.cursor.apiKey') }}</p>
            <code class="text-sm font-mono text-gray-800 dark:text-dark-200">sk-your-api-key</code>
          </div>
          <div class="rounded-lg bg-gray-100 p-3 dark:bg-dark-700">
            <p class="text-sm text-gray-600 dark:text-dark-300">{{ t('tutorial.tools.cursor.baseUrl') }}</p>
            <code class="text-sm font-mono text-gray-800 dark:text-dark-200">{{ apiBaseUrl }}/v1</code>
          </div>
        </div>
      </div>

      <!-- API Usage -->
      <div class="card p-6">
        <h2 class="mb-4 text-xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.apiUsage.title') }}</h2>
        <div class="space-y-4">
          <div>
            <h3 class="mb-2 font-medium text-gray-900 dark:text-white">{{ t('tutorial.apiUsage.endpoint') }}</h3>
            <div class="rounded-lg bg-gray-100 p-3 font-mono text-sm dark:bg-dark-700">
              <code class="text-gray-800 dark:text-dark-200">{{ apiBaseUrl }}/v1/messages</code>
            </div>
          </div>
          <div>
            <h3 class="mb-2 font-medium text-gray-900 dark:text-white">{{ t('tutorial.apiUsage.curlExample') }}</h3>
            <div class="rounded-lg bg-gray-100 p-3 font-mono text-sm dark:bg-dark-700">
              <pre class="overflow-x-auto text-gray-800 dark:text-dark-200">{{ curlExample }}</pre>
            </div>
          </div>
        </div>
      </div>

      <!-- Python SDK -->
      <div class="card p-6">
        <h2 class="mb-4 text-xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.pythonSdk.title') }}</h2>
        <div class="space-y-4">
          <div>
            <h3 class="mb-2 font-medium text-gray-900 dark:text-white">{{ t('tutorial.pythonSdk.install') }}</h3>
            <div class="rounded-lg bg-gray-100 p-3 font-mono text-sm dark:bg-dark-700">
              <code class="text-gray-800 dark:text-dark-200">pip install anthropic</code>
            </div>
          </div>
          <div>
            <h3 class="mb-2 font-medium text-gray-900 dark:text-white">{{ t('tutorial.pythonSdk.example') }}</h3>
            <div class="rounded-lg bg-gray-100 p-3 font-mono text-sm dark:bg-dark-700">
              <pre class="overflow-x-auto text-gray-800 dark:text-dark-200">{{ pythonExample }}</pre>
            </div>
          </div>
        </div>
      </div>

      <!-- OpenAI Compatible -->
      <div class="card p-6">
        <h2 class="mb-4 text-xl font-bold text-gray-900 dark:text-white">{{ t('tutorial.openai.title') }}</h2>
        <div class="space-y-4">
          <p class="text-gray-600 dark:text-dark-300">{{ t('tutorial.openai.description') }}</p>
          <div>
            <h3 class="mb-2 font-medium text-gray-900 dark:text-white">{{ t('tutorial.openai.endpoint') }}</h3>
            <div class="rounded-lg bg-gray-100 p-3 font-mono text-sm dark:bg-dark-700">
              <code class="text-gray-800 dark:text-dark-200">{{ apiBaseUrl }}/v1/chat/completions</code>
            </div>
          </div>
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
              <li>• {{ t('tutorial.tips.tip1') }}</li>
              <li>• {{ t('tutorial.tips.tip2') }}</li>
              <li>• {{ t('tutorial.tips.tip3') }}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'

const { t } = useI18n()

const apiBaseUrl = computed(() => window.location.origin)

const claudeCodeConfig = computed(() => `{
  "env": {
    "ANTHROPIC_API_KEY": "sk-your-api-key",
    "ANTHROPIC_BASE_URL": "${apiBaseUrl.value}/",
    "ANTHROPIC_MODEL": "claude-sonnet-4-20250514"
  }
}`)

const openclawConfig = computed(() => `{
  "provider": "anthropic",
  "base_url": "${apiBaseUrl.value}/",
  "api": "anthropic-messages",
  "api_key": "sk-your-api-key",
  "model": {
    "id": "claude-haiku-4-5-20251001",
    "name": "claude-haiku-4-5-20251001"
  }
}`)

const opencodeConfig = computed(() => `{
  "provider": {
    "openai": {
      "options": {
        "baseURL": "${apiBaseUrl.value}/v1",
        "apiKey": "sk-your-api-key"
      },
      "models": {
        "gpt-4o": {
          "name": "GPT-4o",
          "options": {
            "store": false
          }
        }
      }
    }
  },
  "$schema": "https://opencode.ai/config.json"
}`)

const curlExample = computed(() => `curl ${apiBaseUrl.value}/v1/messages \\
  -H "Content-Type: application/json" \\
  -H "x-api-key: sk-your-api-key" \\
  -H "anthropic-version: 2023-06-01" \\
  -d '{
    "model": "claude-sonnet-4-20250514",
    "max_tokens": 1024,
    "messages": [
      {"role": "user", "content": "Hello, Claude"}
    ]
  }'`)

const pythonExample = computed(() => `from anthropic import Anthropic

client = Anthropic(
    api_key="sk-your-api-key",
    base_url="${apiBaseUrl.value}"
)

message = client.messages.create(
    model="claude-sonnet-4-20250514",
    max_tokens=1024,
    messages=[
        {"role": "user", "content": "Hello, Claude"}
    ]
)
print(message.content)`)
</script>
