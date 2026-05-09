import type { Task } from '../types'

export function TaskTimeline({ tasks, activeId, onSelect, onRetry, onCancel }: { tasks: Task[]; activeId?: string; onSelect: (id: string) => void; onRetry: (id: string) => void; onCancel: (id: string) => void }) {
  return (
    <aside className="task-timeline">
      <header className="timeline-header">
        <p className="eyebrow">请求栏</p>
        <h2>请求</h2>
        <span>{tasks.length ? `${tasks.length} 个请求` : '暂无请求'}</span>
      </header>
      {!tasks.length ? <p className="muted">每次提交都会形成一个请求，并按时间显示在这里。</p> : null}
      <div className="task-stack">
        {tasks.map((task) => (
          <article className={`timeline-item ${task.id === activeId ? 'active' : ''}`} key={task.id} onClick={() => onSelect(task.id)}>
            <div className="timeline-status">
              <strong>{modeLabel(task)} · {sizeLabel(task)} · {task.results.filter((result) => result.ok).length}/{task.count}</strong>
              <span>{task.statusText} / {task.statusCode}</span>
              <span>{task.stageText} / {task.stageCode}</span>
            </div>
            <progress value={task.progress} max={100} />
            <small>{ratioLabel(task)} · {resolutionLabel(task)} · 并发 {task.concurrency}</small>
            <Thumbs task={task} />
            <div className="task-actions">
              {isFinal(task) ? <button type="button" onClick={(event) => { event.stopPropagation(); onRetry(task.id) }}>重试</button> : <button type="button" onClick={(event) => { event.stopPropagation(); onCancel(task.id) }}>取消</button>}
            </div>
          </article>
        ))}
      </div>
    </aside>
  )
}

function Thumbs({ task }: { task: Task }) {
  const images = task.results.filter((result) => result.ok && (result.remoteThumbUrl || result.imageUrl)).slice(0, 4)
  if (!images.length) return null
  const extra = task.results.filter((result) => result.ok && (result.remoteThumbUrl || result.imageUrl)).length - images.length
  return (
    <div className="timeline-thumbs">
      {images.map((result) => <img key={result.index} src={result.remoteThumbUrl || result.imageUrl} alt={`结果缩略图 ${result.index + 1}`} />)}
      {extra > 0 ? <span>+{extra}</span> : null}
    </div>
  )
}

function isFinal(task: Task) {
  return ['succeeded', 'partial_failed', 'failed', 'cancelled', 'interrupted'].includes(task.status)
}

function modeLabel(task: Task) {
  return task.mode === 'image-to-image' ? '图生图' : '文生图'
}

function sizeLabel(task: Task) {
  return task.size && task.size !== '自动' ? task.size : '自动尺寸'
}

function ratioLabel(task: Task) {
  return task.ratio === 'auto' ? '自动比例' : `比例 ${task.ratio}`
}

function resolutionLabel(task: Task) {
  const labels: Record<string, string> = { auto: '自动清晰度', standard: '标准', '2k': '2K', '4k': '4K' }
  return labels[task.resolution] || task.resolution
}
