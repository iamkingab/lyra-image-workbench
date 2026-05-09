import type { Task } from '../types'
import { formatBytes } from '../lib/format'

export function ResultGrid({ task }: { task?: Task }) {
  return (
    <section className="panel results-panel">
      <h2>生成结果</h2>
      {!task ? <p className="muted">选择或创建一个任务后查看结果</p> : null}
      {task ? <p className="muted">{task.stageText} / {task.stage} / {task.stageCode} · {task.progress}%</p> : null}
      <div className="result-grid">
        {task?.results.map((result) => (
          <article className="result-card" key={result.index}>
            {result.ok && result.imageUrl ? <img src={result.imageUrl} alt={`result-${result.index + 1}`} /> : <div className="result-error">{result.error || result.statusText}</div>}
            <footer>{result.statusText} / {result.status} / {result.statusCode} · {formatBytes(result.bytes)}</footer>
          </article>
        ))}
      </div>
    </section>
  )
}
