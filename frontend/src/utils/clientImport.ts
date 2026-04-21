import type { ApiKey } from '@/types'

export type SupportedCcSwitchApp = 'claude' | 'gemini' | 'codex'

export interface CcSwitchImportTarget {
  app: SupportedCcSwitchApp
  endpoint: string
}

export interface CherryStudioImportConfig {
  providerType: 'OpenAI' | 'Anthropic' | 'Gemini'
  type: 'openai' | 'anthropic' | 'gemini'
  providerId: string
  providerName: string
  baseUrl: string
  docsEntry: string
}

export type ModelFamily = 'anthropic' | 'openai' | 'gemini' | 'unknown'

export function detectModelFamily(model: string): ModelFamily {
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

export function getCherryStudioProvider(model: string, baseUrl: string): CherryStudioImportConfig {
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

export function getCherryStudioImportUrl(apiKey: string, model: string, baseUrl: string) {
  const provider = getCherryStudioProvider(model, baseUrl)
  const payload = JSON.stringify({
    id: provider.providerId,
    name: provider.providerName,
    type: provider.type,
    apiKey,
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

export function resolveCcSwitchImportTarget(
  apiKeyRow: Pick<ApiKey, 'group'>,
  app: SupportedCcSwitchApp,
  baseUrl: string
): CcSwitchImportTarget | null {
  const platform = apiKeyRow.group?.platform
  if (!platform) {
    return null
  }

  if (platform === 'antigravity') {
    if (app === 'codex') {
      return null
    }
    return {
      app,
      endpoint: `${baseUrl}/antigravity`,
    }
  }

  if (platform === 'multi') {
    return {
      app,
      endpoint: baseUrl,
    }
  }

  if (platform === 'anthropic' && app === 'claude') {
    return { app, endpoint: baseUrl }
  }

  if (platform === 'gemini' && app === 'gemini') {
    return { app, endpoint: baseUrl }
  }

  if (platform === 'openai' && app === 'codex') {
    return { app, endpoint: baseUrl }
  }

  return null
}

export function buildCcSwitchImportUrl(
  apiKeyRow: Pick<ApiKey, 'key' | 'group'>,
  app: SupportedCcSwitchApp,
  baseUrl: string
) {
  const target = resolveCcSwitchImportTarget(apiKeyRow, app, baseUrl)
  if (!target) {
    return null
  }

  const usageScript = `({
    request: {
      url: "{{baseUrl}}/usage",
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
    app: target.app,
    name: 'ccoder.me',
    homepage: baseUrl,
    endpoint: target.endpoint,
    apiKey: apiKeyRow.key,
    configFormat: 'json',
    usageEnabled: 'true',
    usageScript: btoa(usageScript),
    usageAutoInterval: '30'
  })

  return `ccswitch://v1/import?${params.toString()}`
}
