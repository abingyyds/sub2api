 文档中心 - Ai Go Code
AlGoCode
ccoder.me Docs

个人中心
Home
快速开始
Node.js 环境安装教程
Claude Code 配置教程
Gemini CLI 配置教程
Codex (OpenAI) 配置教程
进阶玩法
OpenClaw 部署教程
OpenCode 配置教程
文章心得
Claude Code 创始人 Boris Cherny 的 13 个使用技巧完整解析
Claude Code 效率提升：用 Slash Commands 把高频操作一键化
Claude Code Hooks 完整指南：让自动化渗透到每个环节
Claude Code Subagents 完整指南：打造你的专属 AI 助手团队
Claude Code 验证策略：让 AI 对自己的工作负责
Build with Ai Go Code
面向开发者的文档入口。从快速接入到完整的 API 参考，这里有你需要的一切。

快速开始
Node.js 环境安装教程

在 Windows、macOS、Linux 安装 Node.js LTS，并验证 node/npm 命令可用。

阅读文档
Claude Code 配置教程

安装并配置 Claude Code CLI 工具，连接 Ai Go Code 平台。

阅读文档
Gemini CLI 配置教程

安装并配置 Gemini CLI 工具，连接 Ai Go Code 平台使用 Gemini 模型。

阅读文档
Codex (OpenAI) 配置教程

安装并配置 Codex CLI 工具，通过配置文件连接 Ai Go Code 平台。

阅读文档
进阶玩法
OpenClaw 部署教程

从零开始部署 OpenClaw Telegram 机器人，并接入 ccoder.me 平台。

阅读文档
OpenCode 配置教程

通过配置文件将 OpenCode 接入 Ai Go Code 平台，快速完成 Claude、GPT 和 Gemini 模型配置。

阅读文档
文章心得
Claude Code 创始人 Boris Cherny 的 13 个使用技巧完整解析

深入解析 Claude Code 创始人 Boris Cherny 分享的 13 个使用技巧，涵盖并行工作、Plan mode、Hooks、Subagents 等高效工作流。

阅读文档
Claude Code 效率提升：用 Slash Commands 把高频操作一键化

详细介绍 Claude Code 自定义 Slash Commands 的创建方法、使用技巧和团队共享实践。

阅读文档
Claude Code Hooks 完整指南：让自动化渗透到每个环节

详解 Claude Code Hooks 机制，涵盖任务通知、等待提醒、代码自动格式化、自动测试等实用场景。

阅读文档
Claude Code Subagents 完整指南：打造你的专属 AI 助手团队

详解 Claude Code Subagents 的创建、配置和使用，包括代码简化、验证代理等实用案例。

阅读文档
Claude Code 验证策略：让 AI 对自己的工作负责

介绍 Claude Code 的多种验证方式，包括直接指令、Subagent 验证、Stop Hook 自动验证和浏览器自动化测试。

阅读文档
 Node.js 环境安装教程 - 文档中心 - Ai Go Code
AlGoCode
ccoder.me Docs

个人中心
Home
快速开始
Node.js 环境安装教程
Claude Code 配置教程
Gemini CLI 配置教程
Codex (OpenAI) 配置教程
进阶玩法
文章心得
Docs
快速开始
快速开始
Node.js 环境安装教程
在 Windows、macOS、Linux 安装 Node.js LTS，并验证 node/npm 命令可用。

后续的 Claude Code、Gemini CLI、Codex 都依赖 Node.js 运行环境。先完成本教程，再继续各自工具配置。

Windows#
方法一：官方下载（推荐）

前往 Node.js 官网 下载 LTS 版本，双击安装包按提示安装即可。

方法二：使用 Chocolatey

Copy
choco install nodejs-lts
方法三：使用 Scoop

Copy
scoop install nodejs-lts
macOS#
方法一：使用 Homebrew（推荐）

Copy
brew install node
方法二：官方下载

前往 Node.js 官网 下载 macOS 安装包。

Linux#
Ubuntu / Debian：

Copy
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs
CentOS / RHEL：

Copy
curl -fsSL https://rpm.nodesource.com/setup_lts.x | sudo bash -
sudo yum install -y nodejs
验证安装#
安装完成后，在终端执行：

Copy
node --version
npm --version
如果能输出版本号，说明 Node.js 与 npm 已可用。

下一步#
Claude Code 配置教程
Gemini CLI 配置教程
Codex (OpenAI) 配置教程
Next
Claude Code 配置教程
On this page

Windows
macOS
Linux
验证安装
下一步  Claude Code 配置教程 - 文档中心 - Ai Go Code
AlGoCode
ccoder.me Docs

个人中心
Home
快速开始
Node.js 环境安装教程
Claude Code 配置教程
Gemini CLI 配置教程
Codex (OpenAI) 配置教程
进阶玩法
文章心得
Docs
快速开始
快速开始
Claude Code 配置教程
安装并配置 Claude Code CLI 工具，连接 Ai Go Code 平台。

Claude Code 是 Anthropic 官方 CLI 工具，强大的 AI 编程助手。本教程将引导你完成安装和配置。

开始前请先完成 Node.js 环境安装教程。

第一步：安装 Claude Code#
打开终端（Windows 使用 PowerShell 或 CMD），运行以下命令：

Copy
npm install -g @anthropic-ai/claude-code
这个命令会从 npm 官方仓库下载并安装最新版本的 Claude Code。

如果在 macOS / Linux 上遇到权限问题，可以使用 sudo：

Copy
sudo npm install -g @anthropic-ai/claude-code
验证安装#
安装完成后，输入以下命令检查是否安装成功：

Copy
claude --version
如果显示版本号，恭喜你！Claude Code 已经成功安装了。

第二步：设置环境变量#
为了让 Claude Code 连接到 Ai Go Code 平台，需要设置两个环境变量：

ANTHROPIC_BASE_URL — 服务地址
ANTHROPIC_AUTH_TOKEN — 你的 API Key
临时设置（当前会话）#
Windows
macOS
Linux
Copy
$env:ANTHROPIC_BASE_URL = "https://ccoder.me"
$env:ANTHROPIC_AUTH_TOKEN = "YOUR_API_KEY"
永久设置#
Windows
macOS (zsh)
Linux (bash)
Copy
[System.Environment]::SetEnvironmentVariable("ANTHROPIC_BASE_URL", "https://ccoder.me", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("ANTHROPIC_AUTH_TOKEN", "YOUR_API_KEY", [System.EnvironmentVariableTarget]::User)
请将 YOUR_API_KEY 替换为在 控制台 页面中复制的任意 API Key。

第三步：开始使用#
现在你可以开始使用 Claude Code 了！

Copy
claude
启动后，Claude Code 会分析当前目录的代码，并提供智能编程辅助。

更多使用说明请参考 Anthropic 官方文档。

Previous
Node.js 环境安装教程
Next
Gemini CLI 配置教程
On this page

第一步：安装 Claude Code
验证安装
第二步：设置环境变量
临时设置（当前会话）
永久设置 
第三步：开始使用  Gemini CLI 配置教程 - 文档中心 - Ai Go Code
AlGoCode
ccoder.me Docs

个人中心
Home
快速开始
Node.js 环境安装教程
Claude Code 配置教程
Gemini CLI 配置教程
Codex (OpenAI) 配置教程
进阶玩法
文章心得
Docs
快速开始
快速开始
Gemini CLI 配置教程
安装并配置 Gemini CLI 工具，连接 Ai Go Code 平台使用 Gemini 模型。

Gemini CLI 是 Google 官方提供的终端 AI 助手工具。本教程将引导你完成安装和配置。

开始前请先完成 Node.js 环境安装教程。

第一步：安装 Gemini CLI#
Windows
macOS (npm)
macOS (Homebrew)
Linux
Copy
npm install -g @google/gemini-cli
验证安装#
Copy
gemini --version
如果显示版本号，说明安装成功了！

第二步：配置环境变量#
设置环境变量让 Gemini CLI 连接到 Ai Go Code 平台。需要设置三个变量：

GOOGLE_GEMINI_BASE_URL — 服务地址
GEMINI_API_KEY — 你的 API Key
GEMINI_MODEL — 使用的模型
临时设置（当前会话）#
Windows
macOS
Linux
Copy
$env:GOOGLE_GEMINI_BASE_URL = "https://ccoder.me"
$env:GEMINI_API_KEY = "YOUR_API_KEY"
$env:GEMINI_MODEL = "gemini-3-pro-preview"
永久设置#
Windows
macOS (zsh)
Linux (bash)
Copy
[System.Environment]::SetEnvironmentVariable("GOOGLE_GEMINI_BASE_URL", "https://ccoder.me", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("GEMINI_API_KEY", "YOUR_API_KEY", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("GEMINI_MODEL", "gemini-3-pro-preview", [System.EnvironmentVariableTarget]::User)
请将 YOUR_API_KEY 替换为在 控制台 页面中复制的任意 API Key。 如需更换模型，只需修改 GEMINI_MODEL 即可。

第三步：开始使用#
配置完成后，在项目目录中运行：

Copy
gemini
Gemini CLI 会分析当前目录的代码，并提供智能建议。

测试命令#
非交互式查询：

Copy
gemini -p "你好"
更多使用说明请参考 Google AI 官方文档。

Previous
Claude Code 配置教程
Next
Codex (OpenAI) 配置教程
On this page

第一步：安装 Gemini CLI
验证安装
第二步：配置环境变量
临时设置（当前会话）
永久设置
第三步：开始使用
测试命令  Codex (OpenAI) 配置教程 - 文档中心 - Ai Go Code
AlGoCode
ccoder.me Docs

个人中心
Home
快速开始
Node.js 环境安装教程
Claude Code 配置教程
Gemini CLI 配置教程
Codex (OpenAI) 配置教程
进阶玩法
文章心得
Docs
快速开始
快速开始
Codex (OpenAI) 配置教程
安装并配置 Codex CLI 工具，通过配置文件连接 Ai Go Code 平台。

Codex 是支持 OpenAI API 兼容接口的 CLI 工具，使用配置文件的方式进行设置。本教程将引导你完成安装和配置。

开始前请先完成 Node.js 环境安装教程。

第一步：安装 Codex#
Windows
macOS
Linux
Copy
npm install -g @openai/codex@latest
以上命令会从 npm 官方仓库下载并安装最新版本的 Codex 工具。

验证安装#
Copy
codex --version
如果显示版本号，恭喜你！Codex 已经成功安装了。

第二步：创建配置文件#
Codex 使用配置文件进行连接设置，需要创建 config.toml 和 auth.json 两个文件。

创建配置目录#
Windows
macOS
Linux
Copy
# 删除旧目录并创建新目录
if (Test-Path "$env:USERPROFILE\.codex") { Remove-Item -Recurse -Force "$env:USERPROFILE\.codex" }
mkdir "$env:USERPROFILE\.codex"
创建配置文件#
Windows
macOS
Linux
Copy
# 创建 config.toml
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
"@ | Out-File -FilePath "$env:USERPROFILE\.codex\config.toml" -Encoding utf8

# 创建 auth.json
@"
{
  "OPENAI_API_KEY": "YOUR_API_KEY"
}
"@ | Out-File -FilePath "$env:USERPROFILE\.codex\auth.json" -Encoding utf8
请将 YOUR_API_KEY 替换为在 控制台 页面中复制的任意 API Key。

第三步：开始使用#
现在你可以开始使用 Codex 了！

Copy
codex
更多使用说明请参考 OpenAI 官方文档。

Previous
Gemini CLI 配置教程
Next
OpenClaw 部署教程
On this page

第一步：安装 Codex
验证安装
第二步：创建配置文件
创建配置目录
创建配置文件
第三步：开始使用  OpenClaw 部署教程 - 文档中心 - Ai Go Code
AlGoCode
ccoder.me Docs

个人中心
Home
快速开始
进阶玩法
OpenClaw 部署教程
OpenCode 配置教程
文章心得
Docs
进阶玩法
进阶玩法
OpenClaw 部署教程
从零开始部署 OpenClaw Telegram 机器人，并接入 ccoder.me 平台。

OpenClaw 从 0 到 1 部署教程（接入 ccoder.me）#
还是有不少人在问龙虾怎么部署、ccoder.me 的服务怎么接入。我直接开了一台全新的服务器，从零开始重新走了一遍完整流程，每一步都截了图，容易踩坑的地方也都标出来了。核心 6 步，跟着做一把过。

准备工作#
一台全新的 Ubuntu/Debian 服务器（推荐 Ubuntu 22.04+）。
服务器已安装 Git：
Copy
sudo apt update && sudo apt install -y git
安装 Git

服务器已安装 Node.js v22+（如未安装，可一键搞定）：
Copy
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo bash -
sudo apt-get install -y nodejs
安装 Node.js

一个 ccoder.me 账号，拿到 API Key（sk-xxx 格式）。
一个 Telegram 账号。
第一步：安装 OpenClaw#
在服务器终端执行一行命令，全局安装 OpenClaw：

Copy
sudo npm install -g openclaw
看到安装成功的提示后，核心程序就位了。

安装 OpenClaw

第二步：申请 Telegram 机器人#
在 Telegram 里搜索 @BotFather（认准蓝色官方认证标记），按以下步骤操作：

BotFather

发送 /newbot。
新建 Bot

按提示给机器人起一个昵称（支持中文，后续可改）。
设置一个唯一的用户名（必须以 bot 结尾，例如 my_ai_bot）。
设置用户名

创建成功后，BotFather 会返回一串 Bot Token（类似 123456789:ABCdefGHIJ...），复制保存好。
获取 Token

安全提醒：Bot Token 相当于机器人的最高控制权密码，不要泄露到公开场合。如果泄露了，可以在 BotFather 里发送 /revoke 重新生成。

第三步：运行向导，完成基础配置#
在终端输入以下命令启动交互式向导：

Copy
openclaw onboard
向导会一步步问你问题，按以下策略选择：

I understand this is powerful and inherently risky. Continue? → 选 Yes。
确认继续

Onboarding mode → 选 QuickStart（新服务器直接用默认配置，最省事）。
QuickStart

Model/auth provider → 翻到最底部，选 Skip for now (configure later)。这是关键一步，因为我们要接入的 ccoder.me 不在官方列表里，所以先跳过，后面手动配。
跳过 Provider

Filter models by provider → 保持 All providers，直接回车。
All Providers

Default model → 选 Keep current（保持默认即可，后面我们会在配置文件里手动指定 ccoder.me 的模型）。
Keep Current

Channels → 选择 Telegram，填入你刚才拿到的 Bot Token。向导会自动帮你完成管理员配对。
选择 Telegram

填入 Token

Configure skills now? → 选 No（Skills 后续随时可以加，先跑起来再说）。
跳过 Skills

Show Homebrew install command? → 选 No（如果上一步选了 Skills → No，这一步可能会被跳过）。
Enable hooks? → 建议全选（boot-md、bootstrap-extra-files、command-logger、session-memory 都挺实用，尤其 session-memory 能在切换会话时自动保存上下文）。
启用 Hooks

Install Gateway service → 选 Yes（如果出现这个选项，表示开启开机自启）。
向导在背后帮你完成了以下工作：

生成配置文件：默认写入 ~/.openclaw/openclaw.json，包含你刚才选择的所有设置。
开启 systemd lingering：确保你退出 SSH 后服务不会被杀掉，Bot 持续在线。
安装 systemd 服务：注册开机自启服务（openclaw-gateway.service），服务器重启后 Bot 自动拉起。
验证 Telegram 连接：自动检测 Bot Token 是否有效，看到 Telegram: ok 就说明连接成功。
验证成功

How do you want to hatch your bot? → 选 Do this later（模型还没配，现在 hatch 会失败，等配好 ccoder.me 再说）。
Hatch 是 OpenClaw 的首次人格初始化流程，给 Bot 设定名字、性格、身份。配好模型重启后直接在 Telegram 聊就行，Bot 没有 hatch 也能正常工作，只是没有自定义人格（会用默认的）。

向导跑完后会提示 Onboarding complete。到此基础环境部署完成。

部署完成

第四步：注入 ccoder.me 模型配置#
向导帮你生成了一个完整的配置文件，但里面没有模型（因为我们刚才选了 Skip）。现在要把 ccoder.me 的模型塞进去。

打开配置文件：

Copy
# 如果有 nano：
nano ~/.openclaw/openclaw.json

# 如果没有 nano，用 vi：
vi ~/.openclaw/openclaw.json

# 或者先装 nano：
sudo apt install -y nano
在文件的最外层大括号内，找到合适的位置（比如紧跟在 meta 块后面），插入以下两段配置片段。注意替换成你真实的 ccoder.me API Key：

jsonc
Copy
"models": {
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
}
踩坑提醒：如果配置文件里已经有 agents 字段（向导可能会生成一个默认的），你需要把它合并到一起，而不是直接粘贴第二份。JSON 里同一层级不能有两个同名 key，否则后面的会覆盖前面的，导致模型配置丢失报错。

比如下面图片是我初始的配置文件，可以看到没有 models 模块，所以可以把完整的 models 模块粘贴进去。

但已经有 agents 模块了，所以我要做的是把 model 字段贴到现有的 agents 里面，而不是再贴一整个 agents 模块进去。

初始配置文件

添加 models 配置

指定默认模型

如果还想加 Codex 模型#
如果你还想加 Codex 模型，可以在 models.providers 下面继续增加一个 provider：

jsonc
Copy
"ccoder.me-gpt": {
  "baseUrl": "https://ccoder.me",
  "apiKey": "填入你的API key",
  "api": "openai-responses",
  "models": [
    {
      "id": "gpt-5.3-codex",
      "name": "GPT-5.3 Codex (ccoder.me)",
      "reasoning": false,
      "input": ["text", "image"],
      "cost": {
        "input": 0,
        "output": 0,
        "cacheRead": 0,
        "cacheWrite": 0
      },
      "contextWindow": 128000,
      "maxTokens": 16384
    }
  ]
}
补充说明#
配置完 providers 后，还需要注意以下两点：

1. 设置默认模型（可选）
如果想把 Codex 模型作为默认模型，可以修改 agents 配置中的 primary 字段：

jsonc
Copy
"agents": {
  "defaults": {
    "model": {
      "primary": "ccoder.me-gpt/gpt-5.3-codex"
    }
  }
}
2. 配置模型白名单（如果需要）
如果你的配置文件中有 models.allowlist 字段，需要把新添加的模型加入白名单才能使用：

jsonc
Copy
"models": {
  "allowlist": [
    "ccoder.me/claude-opus-4-6",
    "ccoder.me-gpt/gpt-5.3-codex"
  ],
  "providers": {
    // 这里保留你现有的 providers 配置
  }
}
提示：如果配置文件中没有 allowlist 字段，说明没有启用白名单机制，所有已配置的模型都可以直接使用。由于上面展示的是配置片段，粘贴到现有 JSON 时记得检查前后逗号是否正确。

保存退出（Ctrl+O 回车保存，Ctrl+X 退出）。

关于 api 字段的说明：

anthropic-messages：用于 Claude 系列模型（Claude Opus 4.6、Claude Sonnet 4 等）。
openai-responses：用于 GPT 系列模型（GPT-4o、GPT-3.5 等），以及大部分兼容 OpenAI 格式的开源模型。
根据你在 ccoder.me 上实际使用的模型类型，选择对应的 api 格式即可。

第五步：重启服务，开始对话#
配置文件修改后，需要重启 OpenClaw 才能生效：

Copy
openclaw gateway restart
openclaw gateway status
只要状态是 active (running)，就说明服务已经正常运行，可以继续下一步配对。

服务运行中

第六步：Telegram 配对#
重启成功后，打开 Telegram，给你的 Bot 发一条任意消息（哈喽、你好）。Bot 会回复一段配对信息：

OpenClaw: access not configured. Your Telegram user id: xxxxxxxx Pairing code: XXXXXXXX

这是正常的，因为你还没有完成身份绑定。复制最后一行命令，回到服务器终端执行：

Copy
openclaw pairing approve telegram <你的配对码>
配对信息

看到 Approved 提示后，回到 Telegram 再发一条消息，Bot 就能正常回复了。部署完成。

部署成功

总结#
整个流程 6 步：安装程序 → 申请机器人 → 跑向导 → 改配置 → 重启 → 配对。向导帮你搞定环境和基础配置，模型接入部分则通过配置文件完成，各司其职。

跑起来之后，就可以直接和机器人对话。后续如果想加更多模型或调整接入配置，也可以继续在现有配置上扩展。

常见问题#
No API key found for provider anthropic

检查配置文件里的 models 和 agents 板块是否正确，primary 是否指向 ccoder.me/claude-opus-4-6。

端口被占用（EADDRINUSE）

说明有其他进程占了默认端口，用以下命令排查：

Copy
ss -tlnp | grep 18789
openclaw: command not found

确认 npm 全局安装成功，运行 npm list -g openclaw 检查。如果安装了但命令找不到，可能是 npm 全局 bin 目录不在 PATH 里，运行 npm config get prefix 查看路径并添加到 PATH。

飞书配置#
如果你想接入飞书而不是 Telegram，可以参考这篇单独教程：OpenClaw 飞书配置教程。

Previous
Codex (OpenAI) 配置教程
Next
OpenCode 配置教程
On this page

OpenClaw 从 0 到 1 部署教程（接入 ccoder.me）
准备工作
第一步：安装 OpenClaw
第二步：申请 Telegram 机器人
第三步：运行向导，完成基础配置
第四步：注入 ccoder.me 模型配置
如果还想加 Codex 模型
补充说明
第五步：重启服务，开始对话
第六步：Telegram 配对
总结
常见问题
飞书配置  OpenCode 配置教程 - 文档中心 - Ai Go Code
AlGoCode
ccoder.me Docs

个人中心
Home
快速开始
进阶玩法
OpenClaw 部署教程
OpenCode 配置教程
文章心得
Docs
进阶玩法
进阶玩法
OpenCode 配置教程
通过配置文件将 OpenCode 接入 Ai Go Code 平台，快速完成 Claude、GPT 和 Gemini 模型配置。

OpenCode 支持接入第三方 AI 服务。如果你想在 OpenCode 里使用 Ai Go Code 平台提供的 Claude、GPT 或 Gemini 模型，可以直接通过配置文件完成接入，整个过程非常快。

这篇文档会带你完成 API Key 获取、配置文件编写和最终验证，适合希望一次性配好的用户。

重要说明#
不同模型族使用的 API 端点后缀不同，请先确认：

Anthropic（Claude 系列）：https://ccoder.me/v1
OpenAI（GPT 系列）：https://ccoder.me/v1
Google（Gemini 系列）：https://ccoder.me/v1beta
必须使用原生 SDK：Anthropic 使用 @ai-sdk/anthropic，OpenAI 使用 OpenCode 内置 OpenAI provider，Gemini 使用 @ai-sdk/google。

准备工作#
开始之前，请先准备好以下内容：

一个 Ai Go Code 账号。
已创建好的 API Key。
已安装并可正常启动的 OpenCode。
获取 API Key#
如果你还没有 API Key，可以按下面步骤创建（⚠️注意选择支持第三方的分组）：

访问 Ai Go Code 并注册账号。
登录后进入控制台。
在“API 密钥”页面创建新的 API Key。
复制并妥善保存生成的 Key（格式通常为 sk-xxxxxx...）。
通过配置文件接入#
这是最推荐的方式。你可以手动修改文件，也可以把下面的配置示例发给能帮你改文件的 AI 工具（如 Codex、Claude Code、Cursor 等），让它直接替你写入。

Claude Code 自动配置示例

第一步：找到配置文件#
OpenCode 主要会用到下面两个文件：

模型配置：~/.config/opencode/opencode.jsonc
认证配置：~/.local/share/opencode/auth.json
第二步：配置模型提供商#
编辑 ~/.config/opencode/opencode.jsonc，将 provider 配置指向 Ai Go Code：

Copy
{
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
}
第三步：确认关键配置#
上面这份配置里，最容易出错的是下面几项：

必须使用原生 SDK，不能随便替换成其他兼容层。
Claude 和 GPT 系列使用 /v1。
Gemini 系列使用 /v1beta。
模型名称必须与 Ai Go Code 支持的模型 ID 完全一致。
apiKey 需要替换成你自己的真实 API Key。
第四步：处理 auth.json#
如果你已经把 API Key 写进了 opencode.jsonc 的 options.apiKey，那么 auth.json 可以保持空对象：

Copy
{}
OpenCode 支持两种 API Key 配置方式：一是写在 opencode.jsonc 的 options.apiKey 中，二是写在 auth.json 中。两者任选其一即可，不需要重复配置。为了维护更简单，推荐直接写在 opencode.jsonc 中。

验证配置是否成功#
完成配置后，重新打开 OpenCode，确认是否能看到你刚刚配置的模型。

建议按下面顺序检查：

模型列表中是否出现 Claude、GPT 或 Gemini 模型。
选择任意一个模型发送测试消息。
确认请求能正常返回结果，没有鉴权或 endpoint 报错。
OpenCode 模型列表

OpenCode 测试成功

常见问题#
为什么 Gemini 不能用 /v1？#
因为 Gemini 需要使用 Google 对应的接口路径，所以在 Ai Go Code 平台接入 Gemini 时也要使用 https://ccoder.me/v1beta。

为什么我已经配置了 auth.json，还要写 apiKey？#
不需要同时写。auth.json 和 opencode.jsonc 二选一即可，只保留一种方式最稳妥。

看不到模型怎么办？#
优先检查以下几项：

配置文件路径是否正确。
JSON / JSONC 格式是否有语法错误。
模型 ID 是否与平台支持的名称完全一致。
baseURL 是否填成了对应模型族需要的地址。
API Key 是否有效。
Previous
OpenClaw 部署教程
Next
Claude Code 创始人 Boris Cherny 的 13 个使用技巧完整解析
On this page

重要说明
准备工作
获取 API Key
通过配置文件接入
第一步：找到配置文件
第二步：配置模型提供商
第三步：确认关键配置
第四步：处理 auth.json
验证配置是否成功
常见问题
为什么 Gemini 不能用 /v1？
为什么我已经配置了 auth.json，还要写 apiKey？
看不到模型怎么办？ 
