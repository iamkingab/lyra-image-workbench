import { type FormEvent, useState } from 'react'
import { openSpace } from '../api/spaces'
import type { SpaceSession } from '../types'

export function SpaceLogin({ onSession }: { onSession: (session: SpaceSession) => void }) {
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  async function submit(event: FormEvent) {
    event.preventDefault()
    setError('')
    try {
      onSession(await openSpace(password))
    } catch (err) {
      setError(err instanceof Error ? err.message : '进入空间失败')
    }
  }
  return (
    <main className="center-shell">
      <form className="panel login-panel" onSubmit={submit}>
        <p className="eyebrow">Personal Space</p>
        <h1>进入个人空间</h1>
        <p className="muted">输入相同空间密码会进入同一个本机个人空间。浏览器只保存不可逆空间令牌。</p>
        <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="至少 10 位复杂空间密码" autoFocus />
        <button className="primary" type="submit">进入工作台</button>
        {error ? <div className="error">{error}</div> : null}
      </form>
    </main>
  )
}
