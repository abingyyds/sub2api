<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Node.js 环境安装教程</h1>
      <p class="mt-2 text-gray-500 dark:text-dark-400">在 Windows、macOS、Linux 安装 Node.js LTS，并验证 node/npm 命令可用。</p>
    </div>

    <DocNote type="info">
      后续的 Claude Code、Gemini CLI、Codex 都依赖 Node.js 运行环境。先完成本教程，再继续各自工具配置。
    </DocNote>

    <!-- Windows -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">Windows</h2>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2">方法一：官方下载（推荐）</h3>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">
        前往 <a href="https://nodejs.org" target="_blank" class="text-primary-600 hover:underline">Node.js 官网</a> 下载 LTS 版本，双击安装包按提示安装即可。
      </p>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-4">方法二：使用 Chocolatey</h3>
      <DocCodeBlock code="choco install nodejs-lts" language="bash" />

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-4">方法三：使用 Scoop</h3>
      <DocCodeBlock code="scoop install nodejs-lts" language="bash" />
    </section>

    <!-- macOS -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">macOS</h2>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2">方法一：使用 Homebrew（推荐）</h3>
      <DocCodeBlock code="brew install node" language="bash" />

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-4">方法二：官方下载</h3>
      <p class="text-sm text-gray-600 dark:text-dark-300">
        前往 <a href="https://nodejs.org" target="_blank" class="text-primary-600 hover:underline">Node.js 官网</a> 下载 macOS 安装包。
      </p>
    </section>

    <!-- Linux -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">Linux</h2>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2">Ubuntu / Debian：</h3>
      <DocCodeBlock :code="`curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -\nsudo apt-get install -y nodejs`" language="bash" />

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-4">CentOS / RHEL：</h3>
      <DocCodeBlock :code="`curl -fsSL https://rpm.nodesource.com/setup_lts.x | sudo bash -\nsudo yum install -y nodejs`" language="bash" />
    </section>

    <!-- Verify -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">验证安装</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">安装完成后，在终端执行：</p>
      <DocCodeBlock :code="`node --version\nnpm --version`" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300">如果能输出版本号，说明 Node.js 与 npm 已可用。</p>
    </section>

    <!-- Next Steps -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">下一步</h2>
      <ul class="space-y-2 text-sm">
        <li>
          <router-link :to="{ path: '/tutorial', query: { doc: 'claude-code' } }" class="text-primary-600 hover:underline">
            Claude Code 配置教程 &rarr;
          </router-link>
        </li>
        <li>
          <router-link :to="{ path: '/tutorial', query: { doc: 'gemini-cli' } }" class="text-primary-600 hover:underline">
            Gemini CLI 配置教程 &rarr;
          </router-link>
        </li>
        <li>
          <router-link :to="{ path: '/tutorial', query: { doc: 'codex' } }" class="text-primary-600 hover:underline">
            Codex (OpenAI) 配置教程 &rarr;
          </router-link>
        </li>
      </ul>
    </section>

    <!-- Prev / Next Navigation -->
    <div class="flex justify-end border-t border-gray-200 dark:border-dark-700 pt-6">
      <router-link :to="{ path: '/tutorial', query: { doc: 'claude-code' } }" class="text-sm font-medium text-primary-600 hover:text-primary-500">
        Claude Code 配置教程 &rarr;
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import DocCodeBlock from '../DocCodeBlock.vue'
import DocNote from '../DocNote.vue'
</script>
