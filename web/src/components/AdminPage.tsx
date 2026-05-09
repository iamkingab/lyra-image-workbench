import { type FormEvent, useEffect, useState } from 'react'
import { getAdminConfig, saveAdminConfig } from '../api/admin'
import type { AdminConfig } from '../types'

export function AdminPage() {
  const [config, setConfig] = useState<AdminConfig | null>(null)
  const [url, setUrl] = useState('')
  const [timeout, setTimeoutSec] = useState(600)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  useEffect(() => {
    void getAdminConfig().then((cfg) => { setConfig(cfg); setUrl(cfg.newApiBaseUrl); setTimeoutSec(cfg.timeoutSec) }).catch((err) => setError(String(err)))
  }, [])
  async function submit(event: FormEvent) {
    event.preventDefault()
    setError('')
    try {
      const cfg = await saveAdminConfig(url, timeout)
      setConfig(cfg)
      setMessage('管理配置已保存')
    } catch (err) {
      setError(err instanceof Error ? err.message : '保存失败')
    }
  }
  return (
    <main className="center-shell">
      <form className="panel admin-panel" onSubmit={submit}>
        <p className="eyebrow">Admin</p>
        <h1>后台管理</h1>
        <label>NewAPI 请求 URL<input value={url} onChange={(e) => setUrl(e.target.value)} placeholder="http://127.0.0.1:3000/v1" /></label>
        <label>超时时间（秒）<input type="number" min={config?.limits.minTimeoutSec || 60} max={config?.limits.maxTimeoutSec || 3600} value={timeout} onChange={(e) => setTimeoutSec(Number(e.target.value))} /></label>
        <div className="status-line">模型：{config?.model || 'gpt-image-2'} {config?.modelLocked ? '（首版固定）' : ''}</div>
        <button className="primary" type="submit">保存管理配置</button>
        <a href="/">返回工作台</a>
        {message ? <div className="ok">{message}</div> : null}
        {error ? <div className="error">{error}</div> : null}
      </form>
    </main>
  )
}
