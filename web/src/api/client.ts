export const SPACE_TOKEN_KEY = 'image-workbench:space-token:v1'

export function getSpaceToken() {
  return localStorage.getItem(SPACE_TOKEN_KEY) || ''
}

export function saveSpaceToken(token: string) {
  localStorage.setItem(SPACE_TOKEN_KEY, token)
}

export function clearSpaceToken() {
  localStorage.removeItem(SPACE_TOKEN_KEY)
}

export async function requestJson<T>(path: string, options: RequestInit = {}, token = getSpaceToken()): Promise<T> {
  const headers = new Headers(options.headers)
  if (token) headers.set('X-Space-Token', token)
  if (options.body && !(options.body instanceof FormData) && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }
  const response = await fetch(path, { ...options, headers, cache: 'no-store' })
  const data = await response.json().catch(() => null) as { message?: string; chinese?: string } | null
  if (!response.ok) throw new Error(data?.chinese || data?.message || `HTTP ${response.status}`)
  return data as T
}
