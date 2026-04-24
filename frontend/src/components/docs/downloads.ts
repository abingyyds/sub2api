export type DownloadLink = {
  label: string
  href: string
  recommended?: boolean
}

export type DownloadGroup = {
  title: string
  links: DownloadLink[]
}

export type DownloadTool = {
  title: string
  version: string
  descZh: string
  descEn: string
  installGuide?: string
  providerDocs?: string
  officialRepo?: string
  releases?: string
  groups: DownloadGroup[]
}

const localDownload = (id: string) => `/downloads/${id}`

export const nodejsLtsVersion = 'v24.15.0'

export const nodejsDownloads: DownloadGroup[] = [
  {
    title: 'Windows',
    links: [
      {
        label: `Windows x64 MSI (${nodejsLtsVersion})`,
        href: `https://nodejs.org/dist/${nodejsLtsVersion}/node-${nodejsLtsVersion}-x64.msi`,
        recommended: true,
      },
      {
        label: `Windows ARM64 ZIP (${nodejsLtsVersion})`,
        href: `https://nodejs.org/dist/${nodejsLtsVersion}/node-${nodejsLtsVersion}-win-arm64.zip`,
      },
    ],
  },
  {
    title: 'macOS',
    links: [
      {
        label: `macOS Intel PKG (${nodejsLtsVersion})`,
        href: `https://nodejs.org/dist/${nodejsLtsVersion}/node-${nodejsLtsVersion}.pkg`,
        recommended: true,
      },
      {
        label: `macOS Apple Silicon TAR (${nodejsLtsVersion})`,
        href: `https://nodejs.org/dist/${nodejsLtsVersion}/node-${nodejsLtsVersion}-darwin-arm64.tar.gz`,
      },
    ],
  },
  {
    title: 'Linux',
    links: [
      {
        label: `Linux x64 TAR (${nodejsLtsVersion})`,
        href: `https://nodejs.org/dist/${nodejsLtsVersion}/node-${nodejsLtsVersion}-linux-x64.tar.xz`,
        recommended: true,
      },
      {
        label: `Linux ARM64 TAR (${nodejsLtsVersion})`,
        href: `https://nodejs.org/dist/${nodejsLtsVersion}/node-${nodejsLtsVersion}-linux-arm64.tar.xz`,
      },
    ],
  },
]

export const tutorialDownloadTools: DownloadTool[] = [
  {
    title: 'CC Switch',
    version: 'v3.13.0',
    descZh: '统一管理 Claude Code、Gemini CLI、Codex 等客户端，支持一键导入 Provider。',
    descEn: 'Manage Claude Code, Gemini CLI, Codex and more in one place, with one-click provider import.',
    officialRepo: 'https://github.com/farion1231/cc-switch',
    releases: 'https://github.com/farion1231/cc-switch/releases/latest',
    groups: [
      {
        title: 'Windows',
        links: [
          {
            label: 'Windows x64 MSI',
            href: localDownload('cc-switch-windows-x64-msi'),
            recommended: true,
          },
          {
            label: 'Windows x64 Portable ZIP',
            href: localDownload('cc-switch-windows-x64-portable-zip'),
          },
        ],
      },
      {
        title: 'macOS',
        links: [
          {
            label: 'macOS Universal DMG',
            href: localDownload('cc-switch-macos-universal-dmg'),
            recommended: true,
          },
        ],
      },
      {
        title: 'Linux',
        links: [
          {
            label: 'Linux x64 DEB',
            href: localDownload('cc-switch-linux-x64-deb'),
            recommended: true,
          },
          {
            label: 'Linux ARM64 DEB',
            href: localDownload('cc-switch-linux-arm64-deb'),
          },
          {
            label: 'Linux x64 AppImage',
            href: localDownload('cc-switch-linux-x64-appimage'),
          },
          {
            label: 'Linux ARM64 AppImage',
            href: localDownload('cc-switch-linux-arm64-appimage'),
          },
        ],
      },
    ],
  },
  {
    title: 'Codex',
    version: '0.122.0',
    descZh: 'OpenAI 官方 Codex 工具，提供 Windows / macOS / Linux 直链安装包。',
    descEn: 'OpenAI official Codex tool with direct Windows, macOS, and Linux package links.',
    installGuide: 'https://github.com/openai/codex#installation',
    releases: 'https://github.com/openai/codex/releases/latest',
    groups: [
      {
        title: 'Windows',
        links: [
          {
            label: 'Windows x64 EXE',
            href: localDownload('codex-windows-x64-exe'),
            recommended: true,
          },
          {
            label: 'Windows ARM64 EXE',
            href: localDownload('codex-windows-arm64-exe'),
          },
        ],
      },
      {
        title: 'macOS',
        links: [
          {
            label: 'macOS Apple Silicon DMG',
            href: localDownload('codex-macos-arm64-dmg'),
            recommended: true,
          },
          {
            label: 'macOS Intel DMG',
            href: localDownload('codex-macos-intel-dmg'),
          },
        ],
      },
      {
        title: 'Linux',
        links: [
          {
            label: 'Linux x64 TAR.GZ',
            href: localDownload('codex-linux-x64-tar-gz'),
            recommended: true,
          },
          {
            label: 'Linux ARM64 TAR.GZ',
            href: localDownload('codex-linux-arm64-tar-gz'),
          },
        ],
      },
    ],
  },
  {
    title: 'Cherry Studio',
    version: 'v1.9.2',
    descZh: 'Cherry Studio 支持接入 OpenAI / Anthropic / Gemini 服务商，可配合本站生成参数使用。',
    descEn: 'Cherry Studio can connect to OpenAI, Anthropic, and Gemini providers and works well with generated settings from this site.',
    providerDocs: 'https://docs.cherry-ai.com/advanced-basic/providers/custom-providers',
    releases: 'https://github.com/CherryHQ/cherry-studio/releases/latest',
    groups: [
      {
        title: 'Windows',
        links: [
          {
            label: 'Windows x64 Setup EXE',
            href: localDownload('cherry-studio-windows-x64-setup-exe'),
            recommended: true,
          },
          {
            label: 'Windows x64 Portable EXE',
            href: localDownload('cherry-studio-windows-x64-portable-exe'),
          },
          {
            label: 'Windows ARM64 Setup EXE',
            href: localDownload('cherry-studio-windows-arm64-setup-exe'),
          },
        ],
      },
      {
        title: 'macOS',
        links: [
          {
            label: 'macOS Intel DMG',
            href: localDownload('cherry-studio-macos-intel-dmg'),
            recommended: true,
          },
          {
            label: 'macOS Apple Silicon DMG',
            href: localDownload('cherry-studio-macos-arm64-dmg'),
          },
        ],
      },
      {
        title: 'Linux',
        links: [
          {
            label: 'Linux x64 AppImage',
            href: localDownload('cherry-studio-linux-x64-appimage'),
            recommended: true,
          },
          {
            label: 'Linux x64 DEB',
            href: localDownload('cherry-studio-linux-x64-deb'),
          },
          {
            label: 'Linux ARM64 AppImage',
            href: localDownload('cherry-studio-linux-arm64-appimage'),
          },
          {
            label: 'Linux ARM64 DEB',
            href: localDownload('cherry-studio-linux-arm64-deb'),
          },
        ],
      },
    ],
  },
]

export const codexRecommendedDownloads = tutorialDownloadTools.find((tool) => tool.title === 'Codex')?.groups ?? []
