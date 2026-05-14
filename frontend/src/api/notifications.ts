import { apiClient } from './client'

export interface WechatNotificationStatus {
  enabled: boolean
  configured: boolean
  bound: boolean
  openid_masked?: string
  bound_at?: string
}

export async function getWechatNotificationStatus(): Promise<WechatNotificationStatus> {
  const { data } = await apiClient.get<WechatNotificationStatus>('/user/wechat-official/status')
  return data
}

export async function getWechatBindURL(returnTo = '/profile'): Promise<string> {
  const { data } = await apiClient.get<{ url: string }>('/user/wechat-official/bind-url', {
    params: { return_to: returnTo }
  })
  return data.url
}

export async function unbindWechatNotification(): Promise<void> {
  await apiClient.post('/user/wechat-official/unbind')
}

export default {
  getWechatNotificationStatus,
  getWechatBindURL,
  unbindWechatNotification
}
