import { useEffect, useMemo, useState } from 'react'
import { createPortal } from 'react-dom'
import { formatBytes } from '../lib/format'

type ImageDimensions = { width: number; height: number }
type PreviewAction = () => void | string | Promise<void | string>
type ViewportSize = { width: number; height: number }

type Props = {
  src: string
  title: string
  requestedSize?: string
  ratio?: string
  bytes?: number
  onCopyImage?: PreviewAction
  onCopyUrl?: PreviewAction
  onDownload?: PreviewAction
  onUseAsReference?: PreviewAction
  onClose: () => void
}

export function ImagePreviewModal({ src, title, bytes, onCopyImage, onCopyUrl, onDownload, onUseAsReference, onClose }: Props) {
  const [dimensions, setDimensions] = useState<ImageDimensions>()
  const [byteSize, setByteSize] = useState(bytes || 0)
  const [notice, setNotice] = useState('')
  const [viewport, setViewport] = useState<ViewportSize>(() => readViewport())

  useEffect(() => {
    setDimensions(undefined)
  }, [src])

  useEffect(() => {
    if (bytes) {
      setByteSize(bytes)
      return
    }
    let cancelled = false
    setByteSize(0)
    void fetch(src)
      .then((response) => response.blob())
      .then((blob) => {
        if (!cancelled) setByteSize(blob.size)
      })
      .catch(() => {
        if (!cancelled) setByteSize(0)
      })
    return () => {
      cancelled = true
    }
  }, [src, bytes])

  useEffect(() => {
    const previousOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose()
    }
    const onResize = () => setViewport(readViewport())
    window.addEventListener('keydown', onKeyDown)
    window.addEventListener('resize', onResize)
    window.visualViewport?.addEventListener('resize', onResize)
    return () => {
      document.body.style.overflow = previousOverflow
      window.removeEventListener('keydown', onKeyDown)
      window.removeEventListener('resize', onResize)
      window.visualViewport?.removeEventListener('resize', onResize)
    }
  }, [onClose])

  const actualRatio = useMemo(() => formatActualRatio(dimensions), [dimensions])
  const previewImageStyle = useMemo(() => getPreviewImageStyle(dimensions, viewport), [dimensions, viewport])

  async function runAction(action: PreviewAction | undefined, fallback: string) {
    if (!action) return
    try {
      const result = await action()
      setNotice(typeof result === 'string' && result ? result : fallback)
    } catch (err) {
      setNotice(err instanceof Error ? err.message : '操作失败')
    }
    window.setTimeout(() => setNotice(''), 1800)
  }

  return createPortal(
    <div className="preview-mask" onMouseDown={(event) => event.target === event.currentTarget && onClose()}>
      <div className="preview-dialog" role="dialog" aria-modal="true" aria-label={title}>
        <button type="button" className="preview-close" onClick={onClose} aria-label="关闭预览">×</button>
        <div className="preview-info">
          <span>{formatDimensions(dimensions)}</span>
          <span>{actualRatio}</span>
          <span>{byteSize ? formatBytes(byteSize) : '读取大小中'}</span>
        </div>
        <div className="preview-stage">
          <img
            src={src}
            alt={title}
            style={previewImageStyle}
            onLoad={(event) => setDimensions({
              width: event.currentTarget.naturalWidth,
              height: event.currentTarget.naturalHeight,
            })}
          />
        </div>
        {notice ? <div className="preview-notice" role="status">{notice}</div> : null}
        <div className="preview-actions">
          {onDownload ? <button type="button" onClick={() => void runAction(onDownload, '下载已触发')}>下载</button> : null}
          {onCopyImage ? <button type="button" onClick={() => void runAction(onCopyImage, '图片已复制')}>复制图片</button> : null}
          {onCopyUrl ? <button type="button" onClick={() => void runAction(onCopyUrl, '链接已复制')}>复制链接</button> : null}
          {onUseAsReference ? <button type="button" onClick={() => void runAction(onUseAsReference, '已加入参考图')}>作为参考图</button> : null}
        </div>
      </div>
    </div>,
    document.body,
  )
}

function readViewport(): ViewportSize {
  return {
    width: Math.round(window.visualViewport?.width || window.innerWidth || 1024),
    height: Math.round(window.visualViewport?.height || window.innerHeight || 768),
  }
}

function getPreviewImageStyle(dimensions: ImageDimensions | undefined, viewport: ViewportSize) {
  if (!dimensions?.width || !dimensions.height) return undefined
  const horizontalPadding = viewport.width <= 620 ? 24 : 56
  const reservedHeight = viewport.width <= 620 ? 142 : 132
  const maxWidth = Math.max(220, viewport.width - horizontalPadding)
  const maxHeight = Math.max(220, viewport.height - reservedHeight)
  const scale = Math.min(maxWidth / dimensions.width, maxHeight / dimensions.height)
  return {
    width: `${Math.max(1, Math.floor(dimensions.width * scale))}px`,
    height: `${Math.max(1, Math.floor(dimensions.height * scale))}px`,
  }
}

function formatDimensions(dimensions?: ImageDimensions) {
  if (!dimensions) return '读取尺寸中'
  return `${dimensions.width}×${dimensions.height}`
}

function formatActualRatio(dimensions?: ImageDimensions) {
  if (!dimensions) return '读取比例中'
  const divisor = gcd(dimensions.width, dimensions.height)
  return `${dimensions.width / divisor}:${dimensions.height / divisor}`
}

function gcd(a: number, b: number): number {
  while (b) {
    const t = b
    b = a % b
    a = t
  }
  return Math.max(1, a)
}
