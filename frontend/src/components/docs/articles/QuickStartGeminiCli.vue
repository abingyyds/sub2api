<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Gemini CLI 配置教程</h1>
      <p class="mt-2 text-gray-500 dark:text-dark-400">安装并配置 Gemini CLI 工具，连接 Ai Go Code 平台使用 Gemini 模型。</p>
    </div>

    <DocNote type="info">
      Gemini CLI 是 Google 官方提供的终端 AI 助手工具。本教程将引导你完成安装和配置。
    </DocNote>

    <DocNote type="warning">
      开始前请先完成
      <router-link :to="{ path: '/tutorial', query: { doc: 'nodejs' } }" class="underline font-medium">Node.js 环境安装教程</router-link>。
    </DocNote>

    <!-- Step 1 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第一步：安装 Gemini CLI</h2>

      <DocTabGroup :tabs="installTabs">
        <template #windows>
          <DocCodeBlock code="npm install -g @google/gemini-cli" language="bash" />
        </template>
        <template #macos-npm>
          <DocCodeBlock code="npm install -g @google/gemini-cli" language="bash" />
        </template>
        <template #macos-brew>
          <DocCodeBlock code="brew install gemini-cli" language="bash" />
        </template>
        <template #linux>
          <DocCodeBlock code="npm install -g @google/gemini-cli" language="bash" />
        </template>
      </DocTabGroup>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-4">验证安装</h3>
      <DocCodeBlock code="gemini --version" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300">如果显示版本号，说明安装成功了！</p>
    </section>

    <!-- Step 2 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第二步：配置环境变量</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-4">设置环境变量让 Gemini CLI 连接到 Ai Go Code 平台。需要设置三个变量：</p>
      <ul class="text-sm text-gray-600 dark:text-dark-300 mb-4 space-y-1">
        <li><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">GOOGLE_GEMINI_BASE_URL</code> — 服务地址</li>
        <li><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">GEMINI_API_KEY</code> — 你的 API Key</li>
        <li><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">GEMINI_MODEL</code> — 使用的模型</li>
      </ul>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2">临时设置（当前会话）</h3>
      <DocTabGroup :tabs="['Windows', 'macOS', 'Linux']">
        <template #Windows>
          <DocCodeBlock :code="tempWindows" filename="PowerShell" />
        </template>
        <template #macOS>
          <DocCodeBlock :code="tempUnix" filename="Terminal" />
        </template>
        <template #Linux>
          <DocCodeBlock :code="tempUnix" filename="Terminal" />
        </template>
      </DocTabGroup>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-4">永久设置</h3>
      <DocTabGroup :tabs="permanentTabs">
        <template #windows>
          <DocCodeBlock :code="permWindows" filename="PowerShell" />
        </template>
        <template #macos>
          <DocCodeBlock :code="permZsh" filename="~/.zshrc" />
        </template>
        <template #linux>
          <DocCodeBlock :code="permBash" filename="~/.bashrc" />
        </template>
      </DocTabGroup>

      <DocNote type="info">
        请将 <code class="bg-blue-100 dark:bg-blue-800/30 px-1.5 py-0.5 rounded text-xs font-mono">YOUR_API_KEY</code> 替换为在
        <router-link to="/keys" class="underline font-medium">控制台</router-link> 页面中复制的任意 API Key。如需更换模型，只需修改 GEMINI_MODEL 即可。
      </DocNote>
    </section>

    <!-- Step 3 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第三步：开始使用</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">配置完成后，在项目目录中运行：</p>
      <DocCodeBlock code="gemini" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-4">Gemini CLI 会分析当前目录的代码，并提供智能建议。</p>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2">测试命令</h3>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">非交互式查询：</p>
      <DocCodeBlock :code='`gemini -p "你好"`' language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300">
        更多使用说明请参考
        <a href="https://ai.google.dev/gemini-api/docs" target="_blank" class="text-primary-600 hover:underline">Google AI 官方文档</a>。
      </p>
    </section>

    <!-- Prev / Next -->
    <div class="flex justify-between border-t border-gray-200 dark:border-dark-700 pt-6">
      <router-link :to="{ path: '/tutorial', query: { doc: 'claude-code' } }" class="text-sm font-medium text-primary-600 hover:text-primary-500">
        &larr; Claude Code 配置教程
      </router-link>
      <router-link :to="{ path: '/tutorial', query: { doc: 'codex' } }" class="text-sm font-medium text-primary-600 hover:text-primary-500">
        Codex (OpenAI) 配置教程 &rarr;
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import DocCodeBlock from '../DocCodeBlock.vue'
import DocTabGroup from '../DocTabGroup.vue'
import DocNote from '../DocNote.vue'

const installTabs = [
  { id: 'windows', label: 'Windows' },
  { id: 'macos-npm', label: 'macOS (npm)' },
  { id: 'macos-brew', label: 'macOS (Homebrew)' },
  { id: 'linux', label: 'Linux' },
]

const permanentTabs = [
  { id: 'windows', label: 'Windows' },
  { id: 'macos', label: 'macOS (zsh)' },
  { id: 'linux', label: 'Linux (bash)' },
]

const tempWindows = `$env:GOOGLE_GEMINI_BASE_URL = "https://ccoder.me"
$env:GEMINI_API_KEY = "YOUR_API_KEY"
$env:GEMINI_MODEL = "gemini-3-pro-preview"`

const tempUnix = `export GOOGLE_GEMINI_BASE_URL="https://ccoder.me"
export GEMINI_API_KEY="YOUR_API_KEY"
export GEMINI_MODEL="gemini-3-pro-preview"`

const permWindows = `[System.Environment]::SetEnvironmentVariable("GOOGLE_GEMINI_BASE_URL", "https://ccoder.me", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("GEMINI_API_KEY", "YOUR_API_KEY", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("GEMINI_MODEL", "gemini-3-pro-preview", [System.EnvironmentVariableTarget]::User)`

const permZsh = `echo 'export GOOGLE_GEMINI_BASE_URL="https://ccoder.me"' >> ~/.zshrc
echo 'export GEMINI_API_KEY="YOUR_API_KEY"' >> ~/.zshrc
echo 'export GEMINI_MODEL="gemini-3-pro-preview"' >> ~/.zshrc
source ~/.zshrc`

const permBash = `echo 'export GOOGLE_GEMINI_BASE_URL="https://ccoder.me"' >> ~/.bashrc
echo 'export GEMINI_API_KEY="YOUR_API_KEY"' >> ~/.bashrc
echo 'export GEMINI_MODEL="gemini-3-pro-preview"' >> ~/.bashrc
source ~/.bashrc`
</script>
