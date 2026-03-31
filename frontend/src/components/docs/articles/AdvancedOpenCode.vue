<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">OpenCode 配置教程</h1>
      <p class="mt-2 text-gray-500 dark:text-dark-400">通过配置文件将 OpenCode 接入 Ai Go Code 平台，快速完成 Claude、GPT 和 Gemini 模型配置。</p>
    </div>

    <DocNote type="info">
      OpenCode 支持接入第三方 AI 服务。如果你想在 OpenCode 里使用 Ai Go Code 平台提供的 Claude、GPT 或 Gemini 模型，可以直接通过配置文件完成接入，整个过程非常快。
    </DocNote>

    <!-- 重要说明 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">重要说明</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">不同模型族使用的 API 端点后缀不同，请先确认：</p>
      <ul class="text-sm text-gray-600 dark:text-dark-300 space-y-1">
        <li><strong>Anthropic（Claude 系列）：</strong><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">https://ccoder.me/v1</code></li>
        <li><strong>OpenAI（GPT 系列）：</strong><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">https://ccoder.me/v1</code></li>
        <li><strong>Google（Gemini 系列）：</strong><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">https://ccoder.me/v1beta</code></li>
      </ul>
      <DocNote type="warning">
        必须使用原生 SDK：Anthropic 使用 <code class="bg-yellow-100 dark:bg-yellow-800/30 px-1.5 py-0.5 rounded text-xs font-mono">@ai-sdk/anthropic</code>，OpenAI 使用 OpenCode 内置 OpenAI provider，Gemini 使用 <code class="bg-yellow-100 dark:bg-yellow-800/30 px-1.5 py-0.5 rounded text-xs font-mono">@ai-sdk/google</code>。
      </DocNote>
    </section>

    <!-- 准备工作 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">准备工作</h2>
      <ul class="text-sm text-gray-600 dark:text-dark-300 space-y-2">
        <li>一个 Ai Go Code 账号</li>
        <li>已创建好的 API Key（在 <router-link to="/keys" class="text-primary-600 hover:underline">控制台</router-link> 创建）</li>
        <li>已安装并可正常启动的 OpenCode</li>
      </ul>
    </section>

    <!-- Step 1 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第一步：找到配置文件</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">OpenCode 主要会用到下面两个文件：</p>
      <ul class="text-sm text-gray-600 dark:text-dark-300 space-y-1">
        <li><strong>模型配置：</strong><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">~/.config/opencode/opencode.jsonc</code></li>
        <li><strong>认证配置：</strong><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">~/.local/share/opencode/auth.json</code></li>
      </ul>
    </section>

    <!-- Step 2 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第二步：配置模型提供商</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">编辑 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">~/.config/opencode/opencode.jsonc</code>，将 provider 配置指向 Ai Go Code：</p>
      <DocCodeBlock :code='`{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "anthropic": {
      "options": {
        "baseURL": "https://ccoder.me/v1",
        "apiKey": "sk-your-api-key-here"
      },
      "npm": "@ai-sdk/anthropic"
    },
    "openai": {
      "options": {
        "baseURL": "https://ccoder.me/v1",
        "apiKey": "sk-your-api-key-here"
      },
      "models": {
        "gpt-5.3-codex": {
          "name": "GPT-5.3 Codex",
          "limit": { "context": 400000, "output": 128000 },
          "options": { "store": false },
          "variants": { "low": {}, "medium": {}, "high": {}, "xhigh": {} }
        }
      }
    },
    "gemini": {
      "options": {
        "baseURL": "https://ccoder.me/v1beta",
        "apiKey": "sk-your-api-key-here"
      },
      "npm": "@ai-sdk/google",
      "models": {
        "gemini-3-flash-preview": {
          "name": "Gemini 3 Flash Preview",
          "limit": { "context": 1048576, "output": 65536 },
          "modalities": { "input": ["text", "image", "pdf"], "output": ["text"] }
        },
        "gemini-3.1-pro-preview": {
          "name": "Gemini 3.1 Pro Preview",
          "limit": { "context": 1048576, "output": 65536 },
          "modalities": { "input": ["text", "image", "pdf"], "output": ["text"] },
          "options": { "thinking": { "budgetTokens": 24576, "type": "enabled" } }
        }
      }
    }
  },
  "agent": {
    "build": { "options": { "store": false } },
    "plan": { "options": { "store": false } }
  }
}`' filename="~/.config/opencode/opencode.jsonc" />
    </section>

    <!-- Step 3 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第三步：确认关键配置</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">上面这份配置里，最容易出错的是下面几项：</p>
      <ul class="text-sm text-gray-600 dark:text-dark-300 space-y-2">
        <li>必须使用原生 SDK，不能随便替换成其他兼容层</li>
        <li>Claude 和 GPT 系列使用 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">/v1</code></li>
        <li>Gemini 系列使用 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">/v1beta</code></li>
        <li>模型名称必须与 Ai Go Code 支持的模型 ID 完全一致</li>
        <li><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">apiKey</code> 需要替换成你自己的真实 API Key</li>
      </ul>
    </section>

    <!-- Step 4 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第四步：处理 auth.json</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">如果你已经把 API Key 写进了 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">opencode.jsonc</code> 的 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">options.apiKey</code>，那么 auth.json 可以保持空对象：</p>
      <DocCodeBlock code="{}" filename="~/.local/share/opencode/auth.json" />
      <DocNote type="tip">
        OpenCode 支持两种 API Key 配置方式：一是写在 opencode.jsonc 的 options.apiKey 中，二是写在 auth.json 中。两者任选其一即可，不需要重复配置。为了维护更简单，推荐直接写在 opencode.jsonc 中。
      </DocNote>
    </section>

    <!-- 验证 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">验证配置是否成功</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">完成配置后，重新打开 OpenCode，确认是否能看到你刚刚配置的模型。建议按下面顺序检查：</p>
      <ol class="text-sm text-gray-600 dark:text-dark-300 space-y-2 list-decimal list-inside">
        <li>模型列表中是否出现 Claude、GPT 或 Gemini 模型</li>
        <li>选择任意一个模型发送测试消息</li>
        <li>确认请求能正常返回结果，没有鉴权或 endpoint 报错</li>
      </ol>
    </section>

    <!-- 常见问题 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">常见问题</h2>
      <div class="space-y-4">
        <div>
          <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-1">为什么 Gemini 不能用 /v1？</h3>
          <p class="text-sm text-gray-600 dark:text-dark-300">因为 Gemini 需要使用 Google 对应的接口路径，所以在 Ai Go Code 平台接入 Gemini 时也要使用 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">https://ccoder.me/v1beta</code>。</p>
        </div>
        <div>
          <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-1">为什么我已经配置了 auth.json，还要写 apiKey？</h3>
          <p class="text-sm text-gray-600 dark:text-dark-300">不需要同时写。auth.json 和 opencode.jsonc 二选一即可，只保留一种方式最稳妥。</p>
        </div>
        <div>
          <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-1">看不到模型怎么办？</h3>
          <p class="text-sm text-gray-600 dark:text-dark-300">优先检查：配置文件路径是否正确、JSON/JSONC 格式是否有语法错误、模型 ID 是否与平台支持的名称完全一致、baseURL 是否填成了对应模型族需要的地址、API Key 是否有效。</p>
        </div>
      </div>
    </section>

    <!-- Prev / Next -->
    <div class="flex justify-start border-t border-gray-200 dark:border-dark-700 pt-6">
      <router-link :to="{ path: '/tutorial', query: { doc: 'openclaw' } }" class="text-sm font-medium text-primary-600 hover:text-primary-500">
        &larr; OpenClaw 部署教程
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import DocCodeBlock from '../DocCodeBlock.vue'
import DocNote from '../DocNote.vue'
</script>
