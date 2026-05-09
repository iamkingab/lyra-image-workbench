import { type FormEvent, useCallback, useEffect, useMemo, useState } from 'react'
import { cancelTask, createTask, listTasks, retryTask } from '../api/tasks'
import { clearSpaceToken, getSpaceToken } from '../api/client'
import { getCurrentSpace, leaveSpace } from '../api/spaces'
import { deleteReferenceUpload, listReferenceUploads, uploadReferenceImages } from '../api/uploads'
import type { CreateTaskRequest, Mode, ReferenceUpload, SpaceSession, Task, TaskEvent } from '../types'
import { SpaceLogin } from './SpaceLogin'
import { SettingsPanel } from './SettingsPanel'
import { UploadPanel } from './UploadPanel'
import { TaskQueue } from './TaskQueue'
import { ResultGrid } from './ResultGrid'
import { useTaskEvents } from '../hooks/useTaskEvents'

export function WorkbenchPage() {
  const [session, setSession] = useState<SpaceSession | null>(null)
  const [spaceReady, setSpaceReady] = useState(false)
  const [keyReady, setKeyReady] = useState(false)
  const [tasks, setTasks] = useState<Task[]>([])
  const [activeId, setActiveId] = useState<string | null>(null)
  const [uploads, setUploads] = useState<ReferenceUpload[]>([])
  const [mode, setMode] = useState<Mode>('text-to-image')
  const [prompt, setPrompt] = useState('')
  const [ratio, setRatio] = useState('1:1')
  const [resolution, setResolution] = useState('standard')
  const [count, setCount] = useState(1)
  const [concurrency, setConcurrency] = useState(1)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  const activeTask = useMemo(() => tasks.find((task) => task.id === activeId), [tasks, activeId])
  const upsertTask = useCallback((task: Task) => {
    setTasks((prev) => {
      const index = prev.findIndex((item) => item.id === task.id)
      if (index < 0) return [task, ...prev]
      const next = [...prev]
      next[index] = task
      return next
    })
  }, [])

  const handleTaskEvent = useCallback((event: TaskEvent) => {
    if (event.event !== 'heartbeat') setMessage(`${event.chinese} / ${event.english} / ${event.code}`)
  }, [])

  useTaskEvents(activeId, upsertTask, handleTaskEvent)

  useEffect(() => {
    const token = getSpaceToken()
    if (!token) return
    void getCurrentSpace(token).then((next) => { setSession(next); setSpaceReady(true) }).catch(() => { clearSpaceToken(); setSpaceReady(false) })
  }, [])

  useEffect(() => {
    if (!spaceReady) return
    void refreshTasks()
    void refreshUploads()
  }, [spaceReady])

  async function refreshTasks() {
    const items = await listTasks()
    setTasks(items)
    if (!activeId && items[0]) setActiveId(items[0].id)
  }

  async function refreshUploads() {
    setUploads(await listReferenceUploads())
  }

  async function submit(event: FormEvent) {
    event.preventDefault()
    setError('')
    if (!keyReady) { setError('请先保存当前空间的 NewAPI Key'); return }
    const payload: CreateTaskRequest = { mode, prompt, ratio, resolution, count, concurrency, uploadIds: mode === 'image-to-image' ? uploads.map((item) => item.id) : [] }
    try {
      const job = await createTask(payload)
      upsertTask(job)
      setActiveId(job.id)
      setMessage('任务已提交，后端会继续执行')
    } catch (err) {
      setError(err instanceof Error ? err.message : '提交失败')
    }
  }

  async function handleUpload(files: File[]) {
    if (!files.length) return
    setUploads([...(await uploadReferenceImages(files)), ...(await listReferenceUploads())])
    await refreshUploads()
  }

  async function handleDeleteUpload(id: string) {
    await deleteReferenceUpload(id)
    await refreshUploads()
  }

  async function logout() {
    await leaveSpace()
    setSession(null)
    setSpaceReady(false)
  }

  if (!session) return <SpaceLogin onSession={(next) => { setSession(next); setSpaceReady(true) }} />

  return (
    <div className="app-shell">
      <header className="topbar">
        <div className="brand">
          <div className="brand-mark">AI</div>
          <div>
            <h1>本机生图工作台</h1>
            <p>{session.space.displayName} · {session.tokenPreview}</p>
          </div>
        </div>
        <nav className="top-actions"><a className="ghost-link" href="/admin">Admin</a><button onClick={logout}>退出空间</button></nav>
      </header>
      <main className="workspace">
        <aside className="left-col">
          <SettingsPanel onReady={setKeyReady} />
          <UploadPanel uploads={uploads} onUpload={handleUpload} onDelete={handleDeleteUpload} />
          <section className="panel">
            <h2>创建任务</h2>
            <form onSubmit={submit} className="task-form">
              <select value={mode} onChange={(e) => setMode(e.target.value as Mode)}><option value="text-to-image">文生图</option><option value="image-to-image">图生图</option></select>
              <textarea value={prompt} onChange={(e) => setPrompt(e.target.value)} placeholder="输入提示词" rows={5} />
              <div className="grid-2"><select value={ratio} onChange={(e) => setRatio(e.target.value)}><option>1:1</option><option>16:9</option><option>9:16</option><option>3:4</option><option>4:3</option><option>auto</option></select><select value={resolution} onChange={(e) => setResolution(e.target.value)}><option value="standard">标准</option><option value="2k">2K</option><option value="4k">4K</option><option value="auto">自动</option></select></div>
              <div className="grid-2"><input type="number" min={1} max={12} value={count} onChange={(e) => setCount(Number(e.target.value))} /><input type="number" min={1} max={4} value={concurrency} onChange={(e) => setConcurrency(Number(e.target.value))} /></div>
              <button className="primary" type="submit">提交任务</button>
            </form>
            {message ? <div className="ok">{message}</div> : null}
            {error ? <div className="error">{error}</div> : null}
          </section>
        </aside>
        <ResultGrid task={activeTask} />
        <TaskQueue tasks={tasks} activeId={activeId || undefined} onSelect={setActiveId} onRetry={(id) => void retryTask(id).then(upsertTask)} onCancel={(id) => void cancelTask(id).then(upsertTask)} />
      </main>
    </div>
  )
}
