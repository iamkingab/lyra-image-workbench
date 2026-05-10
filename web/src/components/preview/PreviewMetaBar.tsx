import type { ImageDimensions } from './PreviewImageStage'

type Props = {
  dimensions?: ImageDimensions
  byteSizeLabel: string
}

export function PreviewMetaBar({ dimensions, byteSizeLabel }: Props) {
  return (
    <div className="image-preview-meta" aria-label={'\u56fe\u7247\u4fe1\u606f'}>
      <span>{formatDimensions(dimensions)}</span>
      <span>{formatActualRatio(dimensions)}</span>
      <span>{byteSizeLabel}</span>
    </div>
  )
}

function formatDimensions(dimensions?: ImageDimensions) {
  if (!dimensions) return '\u8bfb\u53d6\u5c3a\u5bf8\u4e2d'
  return `${dimensions.width}\u00d7${dimensions.height}`
}

function formatActualRatio(dimensions?: ImageDimensions) {
  if (!dimensions) return '\u8bfb\u53d6\u6bd4\u4f8b\u4e2d'
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
