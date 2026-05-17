import { requestJson } from './client'
import type { UserConfig } from '../types'
import { mergeLocalApiKeys, saveLocalApiKeys } from '../lib/localApiKeys'

export async function getUserConfig() {
  const data = await requestJson<{ ok: boolean; config: UserConfig }>('/api/config')
  return mergeLocalApiKeys(data.config)
}

export async function saveApiKey(apiKey: string) {
  return saveUserConfig({ apiKey })
}

export type SaveUserConfigPayload = {
  apiKey?: string
  bananaApiKey?: string
  saveApiKeyToCloud?: boolean
  saveBananaKeyToCloud?: boolean
  clearCloudApiKey?: boolean
  clearCloudBananaApiKey?: boolean
  defaultCount?: number
  defaultConcurrency?: number
  autoUploadPixhost?: boolean
}

export async function saveUserConfig(payload: SaveUserConfigPayload) {
  const { apiKey, bananaApiKey, saveApiKeyToCloud, saveBananaKeyToCloud, ...rest } = payload
  saveLocalApiKeys({ apiKey, bananaApiKey })
  const serverPayload = {
    ...rest,
    ...(saveApiKeyToCloud ? { apiKey, saveApiKeyToCloud } : {}),
    ...(saveBananaKeyToCloud ? { bananaApiKey, saveBananaKeyToCloud } : {}),
  }
  const data = await requestJson<{ ok: boolean; config: UserConfig }>('/api/config', {
    method: 'POST',
    body: JSON.stringify(serverPayload),
  })
  return mergeLocalApiKeys(data.config)
}
