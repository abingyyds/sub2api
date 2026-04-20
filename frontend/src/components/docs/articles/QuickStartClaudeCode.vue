<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Claude Code 配置教程</h1>
      <p class="mt-2 text-gray-500 dark:text-dark-400">安装并配置 Claude Code CLI 工具，连接 Ai Go Code 平台。</p>
    </div>

    <DocNote type="info">
      Claude Code 是 Anthropic 官方 CLI 工具。官方当前推荐优先使用原生安装脚本，npm 安装仅建议在受限环境下作为备用方案。
    </DocNote>

    <DocNote type="warning">
      如果你的网络环境限制 <code class="bg-amber-100 dark:bg-amber-800/30 px-1.5 py-0.5 rounded text-xs font-mono">curl</code> 或 <code class="bg-amber-100 dark:bg-amber-800/30 px-1.5 py-0.5 rounded text-xs font-mono">PowerShell</code> 在线脚本执行，再改用 npm 备用安装；只有备用方案才需要
      <router-link :to="{ path: '/tutorial', query: { doc: 'nodejs' } }" class="underline font-medium">Node.js 环境</router-link>。
    </DocNote>

    <!-- Step 1 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第一步：安装 Claude Code</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">Anthropic 官方当前推荐优先使用原生安装脚本。根据你的系统选择对应命令即可：</p>

      <DocTabGroup :tabs="installTabs">
        <template #native>
          <DocCodeBlock code="curl -fsSL https://claude.ai/install.sh | bash" language="bash" />
          <p class="text-sm text-gray-600 dark:text-dark-300 mt-2">适用于 macOS、Linux 与 WSL。安装完成后建议重新打开终端一次，确保 PATH 已刷新。</p>
        </template>
        <template #windows>
          <DocCodeBlock code="irm https://claude.ai/install.ps1 | iex" filename="PowerShell" />
          <p class="text-sm text-gray-600 dark:text-dark-300 mt-2">适用于 Windows PowerShell。若遇到执行策略限制，请先以管理员身份打开 PowerShell 后重试。</p>
        </template>
        <template #npm>
          <DocCodeBlock code="npm install -g @anthropic-ai/claude-code" language="bash" />
          <p class="text-sm text-gray-600 dark:text-dark-300 mt-2">这是备用方案，仅建议在无法使用官方原生安装脚本时使用。使用前请先准备好 Node.js 环境。</p>
        </template>
      </DocTabGroup>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-4">验证安装</h3>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">安装完成后，输入以下命令检查是否安装成功：</p>
      <DocCodeBlock code="claude --version" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300">如果显示版本号，恭喜你！Claude Code 已经成功安装了。</p>
    </section>

    <!-- Step 2 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第二步：设置环境变量</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-4">为了让 Claude Code 连接到 Ai Go Code 平台，需要设置两个环境变量：</p>
      <ul class="text-sm text-gray-600 dark:text-dark-300 mb-4 space-y-1">
        <li><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">ANTHROPIC_BASE_URL</code> — 服务地址</li>
        <li><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">ANTHROPIC_AUTH_TOKEN</code> — 你的 API Key</li>
      </ul>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2">临时设置（当前会话）</h3>
      <DocTabGroup :tabs="['Windows', 'macOS', 'Linux']">
        <template #Windows>
          <DocCodeBlock :code="envTempWindows" filename="PowerShell" />
        </template>
        <template #macOS>
          <DocCodeBlock :code="envTempUnix" filename="Terminal" />
        </template>
        <template #Linux>
          <DocCodeBlock :code="envTempUnix" filename="Terminal" />
        </template>
      </DocTabGroup>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-4">永久设置</h3>
      <DocTabGroup :tabs="permanentTabs">
        <template #windows>
          <DocCodeBlock :code="envPermWindows" filename="PowerShell" />
        </template>
        <template #macos>
          <DocCodeBlock :code="envPermZsh" filename="~/.zshrc" />
        </template>
        <template #linux>
          <DocCodeBlock :code="envPermBash" filename="~/.bashrc" />
        </template>
      </DocTabGroup>

      <DocNote type="info">
        请将 <code class="bg-blue-100 dark:bg-blue-800/30 px-1.5 py-0.5 rounded text-xs font-mono">YOUR_API_KEY</code> 替换为在
        <router-link to="/keys" class="underline font-medium">控制台</router-link> 页面中复制的任意 API Key。
      </DocNote>
    </section>

    <!-- Step 3 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第三步：开始使用</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">现在你可以开始使用 Claude Code 了！</p>
      <DocCodeBlock code="claude" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300">
        启动后，Claude Code 会分析当前目录的代码，并提供智能编程辅助。更多使用说明请参考
        <a href="https://code.claude.com/docs/en/setup#native-install-recommended" target="_blank" rel="noopener noreferrer" class="text-primary-600 hover:underline">Claude Code 官方安装文档</a>。
      </p>
    </section>

    <!-- Prev / Next -->
    <div class="flex justify-between border-t border-gray-200 dark:border-dark-700 pt-6">
      <router-link :to="{ path: '/tutorial', query: { doc: 'nodejs' } }" class="text-sm font-medium text-primary-600 hover:text-primary-500">
        &larr; Node.js 环境安装教程
      </router-link>
      <router-link :to="{ path: '/tutorial', query: { doc: 'gemini-cli' } }" class="text-sm font-medium text-primary-600 hover:text-primary-500">
        Gemini CLI 配置教程 &rarr;
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import DocCodeBlock from '../DocCodeBlock.vue'
import DocTabGroup from '../DocTabGroup.vue'
import DocNote from '../DocNote.vue'

const installTabs = [
  { id: 'native', label: 'macOS / Linux / WSL' },
  { id: 'windows', label: 'Windows PowerShell' },
  { id: 'npm', label: 'npm 备用方案' },
]

const permanentTabs = [
  { id: 'windows', label: 'Windows' },
  { id: 'macos', label: 'macOS (zsh)' },
  { id: 'linux', label: 'Linux (bash)' },
]

const envTempWindows = `$env:ANTHROPIC_BASE_URL = "https://ccoder.me"
$env:ANTHROPIC_AUTH_TOKEN = "YOUR_API_KEY"`

const envTempUnix = `export ANTHROPIC_BASE_URL="https://ccoder.me"
export ANTHROPIC_AUTH_TOKEN="YOUR_API_KEY"`

const envPermWindows = `[System.Environment]::SetEnvironmentVariable("ANTHROPIC_BASE_URL", "https://ccoder.me", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("ANTHROPIC_AUTH_TOKEN", "YOUR_API_KEY", [System.EnvironmentVariableTarget]::User)`

const envPermZsh = `echo 'export ANTHROPIC_BASE_URL="https://ccoder.me"' >> ~/.zshrc
echo 'export ANTHROPIC_AUTH_TOKEN="YOUR_API_KEY"' >> ~/.zshrc
source ~/.zshrc`

const envPermBash = `echo 'export ANTHROPIC_BASE_URL="https://ccoder.me"' >> ~/.bashrc
echo 'export ANTHROPIC_AUTH_TOKEN="YOUR_API_KEY"' >> ~/.bashrc
source ~/.bashrc`
</script>
