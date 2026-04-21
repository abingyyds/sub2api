<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Codex (OpenAI) 配置教程</h1>
      <p class="mt-2 text-gray-500 dark:text-dark-400">安装并配置 Codex CLI 工具，通过配置文件连接 Ai Go Code 平台。</p>
    </div>

    <DocNote type="info">
      Codex 是支持 OpenAI API 兼容接口的 CLI 工具，使用配置文件的方式进行设置。本教程将引导你完成安装和配置。
    </DocNote>

    <DocNote type="warning">
      开始前请先完成
      <router-link :to="{ path: '/tutorial', query: { doc: 'nodejs' } }" class="underline font-medium">Node.js 环境安装教程</router-link>。
    </DocNote>

    <!-- Step 1 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第一步：安装 Codex</h2>

      <DocTabGroup :tabs="['Windows', 'macOS', 'Linux']">
        <template #Windows>
          <DocCodeBlock code="npm install -g @openai/codex@latest" language="bash" />
        </template>
        <template #macOS>
          <DocCodeBlock code="npm install -g @openai/codex@latest" language="bash" />
        </template>
        <template #Linux>
          <DocCodeBlock code="npm install -g @openai/codex@latest" language="bash" />
        </template>
      </DocTabGroup>

      <p class="text-sm text-gray-600 dark:text-dark-300 mt-2 mb-2">以上命令会从 npm 官方仓库下载并安装最新版本的 Codex 工具。</p>
      <div class="mb-2 space-y-3">
        <p class="text-sm text-gray-600 dark:text-dark-300">如果你更希望直接下载安装包，可以直接选择对应系统版本：</p>
        <div
          v-for="group in codexDownloads"
          :key="group.title"
          class="rounded-xl bg-gray-50 p-3 dark:bg-dark-800/70"
        >
          <p class="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">{{ group.title }}</p>
          <div class="mt-2 flex flex-wrap gap-2">
            <a
              v-for="link in group.links"
              :key="link.href"
              :href="link.href"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center rounded-lg border px-3 py-2 text-sm font-medium transition-colors"
              :class="link.recommended
                ? 'border-primary-200 bg-primary-50 text-primary-700 hover:border-primary-300 dark:border-primary-800 dark:bg-primary-900/20 dark:text-primary-300'
                : 'border-gray-200 text-gray-700 hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-400'"
            >
              {{ link.label }}
            </a>
          </div>
        </div>
      </div>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-4">验证安装</h3>
      <DocCodeBlock code="codex --version" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300">如果显示版本号，恭喜你！Codex 已经成功安装了。</p>
    </section>

    <!-- Step 2 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第二步：创建配置文件</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-4">Codex 使用配置文件进行连接设置，需要创建 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">config.toml</code> 和 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">auth.json</code> 两个文件。</p>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2">创建配置目录</h3>
      <DocTabGroup :tabs="['Windows', 'macOS', 'Linux']">
        <template #Windows>
          <DocCodeBlock :code="mkdirWindows" filename="PowerShell" />
        </template>
        <template #macOS>
          <DocCodeBlock :code="mkdirUnix" filename="Terminal" />
        </template>
        <template #Linux>
          <DocCodeBlock :code="mkdirUnix" filename="Terminal" />
        </template>
      </DocTabGroup>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-4">创建配置文件</h3>
      <DocTabGroup :tabs="['Windows', 'macOS', 'Linux']">
        <template #Windows>
          <DocCodeBlock :code="configWindows" filename="PowerShell" />
        </template>
        <template #macOS>
          <DocCodeBlock :code="configUnix" filename="Terminal" />
        </template>
        <template #Linux>
          <DocCodeBlock :code="configUnix" filename="Terminal" />
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
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">现在你可以开始使用 Codex 了！</p>
      <DocCodeBlock code="codex" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300">
        更多使用说明请参考
        <a href="https://github.com/openai/codex" target="_blank" class="text-primary-600 hover:underline">OpenAI 官方文档</a>。
      </p>
    </section>

    <!-- Prev / Next -->
    <div class="flex justify-between border-t border-gray-200 dark:border-dark-700 pt-6">
      <router-link :to="{ path: '/tutorial', query: { doc: 'gemini-cli' } }" class="text-sm font-medium text-primary-600 hover:text-primary-500">
        &larr; Gemini CLI 配置教程
      </router-link>
      <router-link :to="{ path: '/tutorial', query: { doc: 'openclaw' } }" class="text-sm font-medium text-primary-600 hover:text-primary-500">
        OpenClaw 部署教程 &rarr;
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import DocCodeBlock from '../DocCodeBlock.vue'
import DocTabGroup from '../DocTabGroup.vue'
import DocNote from '../DocNote.vue'
import { codexRecommendedDownloads } from '../downloads'

const codexDownloads = codexRecommendedDownloads

const mkdirWindows = `# 删除旧目录并创建新目录
if (Test-Path "$env:USERPROFILE\\.codex") { Remove-Item -Recurse -Force "$env:USERPROFILE\\.codex" }
mkdir "$env:USERPROFILE\\.codex"`

const mkdirUnix = `rm -rf ~/.codex
mkdir -p ~/.codex`

const configWindows = `# 创建 config.toml
@"
model_provider = "ccoder.me"
model = "gpt-5.3-codex"
model_reasoning_effort = "high"
disable_response_storage = true
preferred_auth_method = "apikey"

[model_providers.ccoder.me]
name = "ccoder.me"
base_url = "https://ccoder.me"
wire_api = "responses"
requires_openai_auth = true
"@ | Out-File -FilePath "$env:USERPROFILE\\.codex\\config.toml" -Encoding utf8

# 创建 auth.json
@"
{
  "OPENAI_API_KEY": "YOUR_API_KEY"
}
"@ | Out-File -FilePath "$env:USERPROFILE\\.codex\\auth.json" -Encoding utf8`

const configUnix = `# 创建 config.toml
cat > ~/.codex/config.toml << 'EOF'
model_provider = "ccoder.me"
model = "gpt-5.3-codex"
model_reasoning_effort = "high"
disable_response_storage = true
preferred_auth_method = "apikey"

[model_providers.ccoder.me]
name = "ccoder.me"
base_url = "https://ccoder.me"
wire_api = "responses"
requires_openai_auth = true
EOF

# 创建 auth.json
cat > ~/.codex/auth.json << 'EOF'
{
  "OPENAI_API_KEY": "YOUR_API_KEY"
}
EOF`
</script>
