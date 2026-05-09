import { requestJson } from './client'
import type { AdminConfig } from '../types'

export async function getAdminConfig() {
  const data = await requestJson<{ ok: boolean; config: AdminConfig }>('/api/admin/config', {}, '')
  return data.config
}

export async function saveAdminConfig(newApiBaseUrl: string, timeoutSec: number) {
  const data = await requestJson<{ ok: boolean; config: AdminConfig }>('/api/admin/config', {
    method: 'POST',
    body: JSON.stringify({ newApiBaseUrl, timeoutSec }),
  }, '')
  return data.config
}
