<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">OpenClaw 部署教程</h1>
      <p class="mt-2 text-gray-500 dark:text-dark-400">从零开始部署 OpenClaw Telegram 机器人，并接入 ccoder.me 平台。</p>
    </div>

    <DocNote type="info">
      核心 6 步，跟着做一把过。本教程从全新服务器开始，每一步都详细说明，容易踩坑的地方也都标出来了。
    </DocNote>

    <!-- 准备工作 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">准备工作</h2>
      <ul class="text-sm text-gray-600 dark:text-dark-300 space-y-2">
        <li>一台全新的 Ubuntu/Debian 服务器（推荐 Ubuntu 22.04+）</li>
        <li>服务器已安装 Git：</li>
      </ul>
      <DocCodeBlock code="sudo apt update && sudo apt install -y git" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">服务器已安装 Node.js v22+（如未安装，可一键搞定）：</p>
      <DocCodeBlock :code="`curl -fsSL https://deb.nodesource.com/setup_22.x | sudo bash -\nsudo apt-get install -y nodejs`" language="bash" />
      <ul class="text-sm text-gray-600 dark:text-dark-300 space-y-2 mt-2">
        <li>一个 ccoder.me 账号，拿到 API Key（<code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">sk-xxx</code> 格式）</li>
        <li>一个 Telegram 账号</li>
      </ul>
    </section>

    <!-- Step 1 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第一步：安装 OpenClaw</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">在服务器终端执行一行命令，全局安装 OpenClaw：</p>
      <DocCodeBlock code="sudo npm install -g openclaw" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300">看到安装成功的提示后，核心程序就位了。</p>
    </section>

    <!-- Step 2 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第二步：申请 Telegram 机器人</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">在 Telegram 里搜索 <strong>@BotFather</strong>（认准蓝色官方认证标记），按以下步骤操作：</p>
      <ol class="text-sm text-gray-600 dark:text-dark-300 space-y-2 list-decimal list-inside">
        <li>发送 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">/newbot</code></li>
        <li>按提示给机器人起一个昵称（支持中文，后续可改）</li>
        <li>设置一个唯一的用户名（必须以 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">bot</code> 结尾，例如 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">my_ai_bot</code>）</li>
        <li>创建成功后，BotFather 会返回一串 Bot Token（类似 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">123456789:ABCdefGHIJ...</code>），复制保存好</li>
      </ol>
      <DocNote type="warning">
        <strong>安全提醒：</strong>Bot Token 相当于机器人的最高控制权密码，不要泄露到公开场合。如果泄露了，可以在 BotFather 里发送 <code class="bg-yellow-100 dark:bg-yellow-800/30 px-1.5 py-0.5 rounded text-xs font-mono">/revoke</code> 重新生成。
      </DocNote>
    </section>

    <!-- Step 3 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第三步：运行向导，完成基础配置</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">在终端输入以下命令启动交互式向导：</p>
      <DocCodeBlock code="openclaw onboard" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-4">向导会一步步问你问题，按以下策略选择：</p>
      <ol class="text-sm text-gray-600 dark:text-dark-300 space-y-3 list-decimal list-inside">
        <li><strong>I understand this is powerful and inherently risky. Continue?</strong> → 选 <strong>Yes</strong></li>
        <li><strong>Onboarding mode</strong> → 选 <strong>QuickStart</strong>（新服务器直接用默认配置，最省事）</li>
        <li><strong>Model/auth provider</strong> → 翻到最底部，选 <strong>Skip for now (configure later)</strong>。这是关键一步，因为我们要接入的 ccoder.me 不在官方列表里，所以先跳过，后面手动配</li>
        <li><strong>Filter models by provider</strong> → 保持 All providers，直接回车</li>
        <li><strong>Default model</strong> → 选 <strong>Keep current</strong>（保持默认即可，后面我们会在配置文件里手动指定 ccoder.me 的模型）</li>
        <li><strong>Channels</strong> → 选择 <strong>Telegram</strong>，填入你刚才拿到的 Bot Token。向导会自动帮你完成管理员配对</li>
        <li><strong>Configure skills now?</strong> → 选 <strong>No</strong>（Skills 后续随时可以加，先跑起来再说）</li>
        <li><strong>Enable hooks?</strong> → 建议全选（boot-md、bootstrap-extra-files、command-logger、session-memory 都挺实用）</li>
        <li><strong>Install Gateway service</strong> → 选 <strong>Yes</strong>（开启开机自启）</li>
        <li><strong>How do you want to hatch your bot?</strong> → 选 <strong>Do this later</strong>（模型还没配，现在 hatch 会失败）</li>
      </ol>

      <DocNote type="tip">
        向导在背后帮你完成了以下工作：生成配置文件（默认写入 <code class="bg-green-100 dark:bg-green-800/30 px-1.5 py-0.5 rounded text-xs font-mono">~/.openclaw/openclaw.json</code>）、开启 systemd lingering、安装 systemd 服务、验证 Telegram 连接。
      </DocNote>
    </section>

    <!-- Step 4 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第四步：注入 ccoder.me 模型配置</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">向导帮你生成了一个完整的配置文件，但里面没有模型（因为我们刚才选了 Skip）。现在要把 ccoder.me 的模型塞进去。</p>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">打开配置文件：</p>
      <DocCodeBlock code="nano ~/.openclaw/openclaw.json" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">在文件的最外层大括号内，找到合适的位置，插入以下两段配置片段。注意替换成你真实的 ccoder.me API Key：</p>
      <DocCodeBlock :code='`"models": {
  "providers": {
    "ccoder.me": {
      "baseUrl": "https://ccoder.me",
      "apiKey": "sk-这里填你的ccoder.me的APIKey",
      "api": "anthropic-messages",
      "models": [
        {
          "id": "claude-opus-4-6",
          "name": "Claude Opus 4.6 (ccoder.me)",
          "reasoning": true,
          "input": ["text", "image"],
          "contextWindow": 200000,
          "maxTokens": 16384
        }
      ]
    }
  }
},
"agents": {
  "defaults": {
    "model": {
      "primary": "ccoder.me/claude-opus-4-6"
    }
  }
}`' filename="~/.openclaw/openclaw.json" />

      <DocNote type="warning">
        <strong>踩坑提醒：</strong>如果配置文件里已经有 <code class="bg-yellow-100 dark:bg-yellow-800/30 px-1.5 py-0.5 rounded text-xs font-mono">agents</code> 字段（向导可能会生成一个默认的），你需要把它合并到一起，而不是直接粘贴。JSON 里同一层级不能有两个同名 key，否则后面的会覆盖前面的。
      </DocNote>

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-6">如果还想加 Codex 模型</h3>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">可以在 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">models.providers</code> 下面继续增加一个 provider：</p>
      <DocCodeBlock :code='`"ccoder.me-gpt": {
  "baseUrl": "https://ccoder.me",
  "apiKey": "填入你的API key",
  "api": "openai-responses",
  "models": [
    {
      "id": "gpt-5.3-codex",
      "name": "GPT-5.3 Codex (ccoder.me)",
      "reasoning": false,
      "input": ["text", "image"],
      "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
      "contextWindow": 128000,
      "maxTokens": 16384
    }
  ]
}`' filename="~/.openclaw/openclaw.json" />

      <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-2 mt-6">关于 api 字段的说明</h3>
      <ul class="text-sm text-gray-600 dark:text-dark-300 space-y-2">
        <li><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">anthropic-messages</code>：用于 Claude 系列模型（Claude Opus 4.6、Claude Sonnet 4 等）</li>
        <li><code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">openai-responses</code>：用于 GPT 系列模型（GPT-4o、GPT-3.5 等），以及大部分兼容 OpenAI 格式的开源模型</li>
      </ul>

      <p class="text-sm text-gray-600 dark:text-dark-300 mt-4">保存退出（Ctrl+O 回车保存，Ctrl+X 退出）。</p>
    </section>

    <!-- Step 5 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第五步：重启服务，开始对话</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">配置文件修改后，需要重启 OpenClaw 才能生效：</p>
      <DocCodeBlock :code="`openclaw gateway restart\nopenclaw gateway status`" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300">只要状态是 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">active (running)</code>，就说明服务已经正常运行，可以继续下一步配对。</p>
    </section>

    <!-- Step 6 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">第六步：Telegram 配对</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">重启成功后，打开 Telegram，给你的 Bot 发一条任意消息。Bot 会回复一段配对信息：</p>
      <div class="bg-gray-100 dark:bg-dark-800 rounded-lg p-3 text-sm font-mono text-gray-700 dark:text-dark-300 my-2">
        OpenClaw: access not configured. Your Telegram user id: xxxxxxxx Pairing code: XXXXXXXX
      </div>
      <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">这是正常的，因为你还没有完成身份绑定。复制配对码，回到服务器终端执行：</p>
      <DocCodeBlock code="openclaw pairing approve telegram <你的配对码>" language="bash" />
      <p class="text-sm text-gray-600 dark:text-dark-300">看到 Approved 提示后，回到 Telegram 再发一条消息，Bot 就能正常回复了。部署完成！</p>
    </section>

    <!-- 总结 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">总结</h2>
      <p class="text-sm text-gray-600 dark:text-dark-300">整个流程 6 步：安装程序 → 申请机器人 → 跑向导 → 改配置 → 重启 → 配对。向导帮你搞定环境和基础配置，模型接入部分则通过配置文件完成，各司其职。</p>
    </section>

    <!-- 常见问题 -->
    <section>
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">常见问题</h2>

      <div class="space-y-4">
        <div>
          <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-1">No API key found for provider anthropic</h3>
          <p class="text-sm text-gray-600 dark:text-dark-300">检查配置文件里的 models 和 agents 板块是否正确，primary 是否指向 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">ccoder.me/claude-opus-4-6</code>。</p>
        </div>

        <div>
          <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-1">端口被占用（EADDRINUSE）</h3>
          <p class="text-sm text-gray-600 dark:text-dark-300 mb-2">说明有其他进程占了默认端口，用以下命令排查：</p>
          <DocCodeBlock code="ss -tlnp | grep 18789" language="bash" />
        </div>

        <div>
          <h3 class="text-base font-medium text-gray-800 dark:text-dark-200 mb-1">openclaw: command not found</h3>
          <p class="text-sm text-gray-600 dark:text-dark-300">确认 npm 全局安装成功，运行 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">npm list -g openclaw</code> 检查。如果安装了但命令找不到，可能是 npm 全局 bin 目录不在 PATH 里，运行 <code class="bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-xs font-mono">npm config get prefix</code> 查看路径并添加到 PATH。</p>
        </div>
      </div>
    </section>

    <!-- Prev / Next -->
    <div class="flex justify-between border-t border-gray-200 dark:border-dark-700 pt-6">
      <router-link :to="{ path: '/tutorial', query: { doc: 'codex' } }" class="text-sm font-medium text-primary-600 hover:text-primary-500">
        &larr; Codex (OpenAI) 配置教程
      </router-link>
      <router-link :to="{ path: '/tutorial', query: { doc: 'opencode' } }" class="text-sm font-medium text-primary-600 hover:text-primary-500">
        OpenCode 配置教程 &rarr;
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import DocCodeBlock from '../DocCodeBlock.vue'
import DocNote from '../DocNote.vue'
</script>
