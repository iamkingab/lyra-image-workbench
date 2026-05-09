import { clearSpaceToken, requestJson, saveSpaceToken } from './client'
import type { SpaceSession } from '../types'

export async function openSpace(password: string) {
  const data = await requestJson<{ ok: boolean; session: SpaceSession }>('/api/spaces/session', {
    method: 'POST',
    body: JSON.stringify({ password }),
  }, '')
  saveSpaceToken(data.session.token)
  return data.session
}

export async function getCurrentSpace(token: string) {
  const data = await requestJson<{ ok: boolean; session: SpaceSession }>('/api/spaces/session', {}, token)
  return data.session
}

export async function leaveSpace() {
  clearSpaceToken()
}
