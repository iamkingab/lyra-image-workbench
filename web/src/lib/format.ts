export function formatStatus(chinese: string, english: string, code: string) {
  return `${chinese} / ${english} / ${code}`
}

export function formatBytes(bytes = 0) {
  if (!bytes) return '0 B'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1024 / 1024).toFixed(2)} MB`
}
