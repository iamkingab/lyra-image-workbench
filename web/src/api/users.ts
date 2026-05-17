import { clearLocalKeyScope, requestJson, saveLocalKeyScope } from './client'
import type { UserSession } from '../types'

type UserSessionResponse = { ok: boolean; session: UserSession }
export type TwoFactorSetup = { secret: string; otpauthUrl: string }

export async function registerUser(username: string, password: string, legacySpacePassword = '') {
  const data = await requestJson<UserSessionResponse>('/api/users/register', {
    method: 'POST',
    body: JSON.stringify({ username, password, legacySpacePassword }),
  })
  saveLocalKeyScope(data.session.user.username)
  return data.session
}

export async function loginUser(username: string, password: string, twoFactorCode = '') {
  const data = await requestJson<UserSessionResponse>('/api/users/session', {
    method: 'POST',
    body: JSON.stringify({ username, password, twoFactorCode }),
  })
  saveLocalKeyScope(data.session.user.username)
  return data.session
}

export async function getCurrentUser() {
  const data = await requestJson<UserSessionResponse>('/api/users/session')
  saveLocalKeyScope(data.session.user.username)
  return data.session
}

export async function logoutUser() {
  await requestJson<{ ok: boolean }>('/api/users/session', { method: 'DELETE' })
  clearLocalKeyScope()
}

export async function setupTwoFactor() {
  const data = await requestJson<{ ok: boolean; setup: TwoFactorSetup }>('/api/users/2fa/setup', { method: 'POST' })
  return data.setup
}

export async function enableTwoFactor(code: string) {
  const data = await requestJson<UserSessionResponse>('/api/users/2fa/enable', {
    method: 'POST',
    body: JSON.stringify({ code }),
  })
  return data.session
}

export async function disableTwoFactor(code: string) {
  const data = await requestJson<UserSessionResponse>('/api/users/2fa/disable', {
    method: 'POST',
    body: JSON.stringify({ code }),
  })
  return data.session
}
