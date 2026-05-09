import type { Task } from '../types'

export function TaskQueue({ tasks, activeId, onSelect, onRetry, onCancel }: { tasks: Task[]; activeId?: string; onSelect: (id: string) => void; onRetry: (id: string) => void; onCancel: (id: string) => void }) {
  return (
    <section className="panel task-panel">
      <h2>任务队列</h2>
      {!tasks.length ? <p className="muted">暂无任务</p> : null}
      {tasks.map((task) => (
        <article className={`task-card ${task.id === activeId ? 'active' : ''}`} key={task.id} onClick={() => onSelect(task.id)}>
          <strong>{task.statusText} / {task.status} / {task.statusCode}</strong>
          <p>{task.stageText} / {task.stage} / {task.stageCode}</p>
          <progress value={task.progress} max={100} />
          <small>{task.mode === 'image-to-image' ? '图生图' : '文生图'} · {task.size === '自动' ? '自动尺寸' : task.size} · {task.results.filter((r) => r.ok).length}/{task.count}</small>
          <div className="task-actions">
            <button type="button" onClick={(e) => { e.stopPropagation(); onRetry(task.id) }}>重试</button>
            <button type="button" onClick={(e) => { e.stopPropagation(); onCancel(task.id) }}>取消</button>
          </div>
        </article>
      ))}
    </section>
  )
}
