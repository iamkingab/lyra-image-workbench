import { requestJson } from './client'
import type { ReferenceUpload } from '../types'

export async function listReferenceUploads() {
  const data = await requestJson<{ ok: boolean; uploads: ReferenceUpload[] }>('/api/uploads/reference')
  return data.uploads
}

export async function uploadReferenceImages(files: File[]) {
  const form = new FormData()
  for (const file of files) form.append('image[]', file, file.name)
  const data = await requestJson<{ ok: boolean; uploads: ReferenceUpload[] }>('/api/uploads/reference', {
    method: 'POST',
    body: form,
  })
  return data.uploads
}

export async function deleteReferenceUpload(id: string) {
  await requestJson<{ ok: boolean }>(`/api/uploads/reference/${encodeURIComponent(id)}`, { method: 'DELETE' })
}
