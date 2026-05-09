import type { ReferenceUpload } from '../types'
import { formatBytes } from '../lib/format'

export function UploadPanel({ uploads, onUpload, onDelete }: { uploads: ReferenceUpload[]; onUpload: (files: File[]) => void; onDelete: (id: string) => void }) {
  return (
    <section className="panel">
      <h2>图生图参考图</h2>
      <input type="file" accept="image/png,image/jpeg,image/webp" multiple onChange={(e) => onUpload(Array.from(e.target.files || []))} />
      <div className="upload-list">
        {uploads.map((item) => (
          <div className="upload-item" key={item.id}>
            <span>{item.originalName}</span>
            <small>{item.mime} · {formatBytes(item.size)}</small>
            <button type="button" onClick={() => onDelete(item.id)}>删除</button>
          </div>
        ))}
      </div>
    </section>
  )
}
