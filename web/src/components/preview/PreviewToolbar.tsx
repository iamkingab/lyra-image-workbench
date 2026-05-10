export type PreviewAction = () => void | string | Promise<void | string>

type Props = {
  onCopyImage?: PreviewAction
  onCopyUrl?: PreviewAction
  onDownload?: PreviewAction
  onUseAsReference?: PreviewAction
  onNotice: (value: string) => void
}

const TEXT = {
  actionFailed: '\u64cd\u4f5c\u5931\u8d25',
  downloadStarted: '\u4e0b\u8f7d\u5df2\u89e6\u53d1',
  imageCopied: '\u56fe\u7247\u5df2\u590d\u5236',
  linkCopied: '\u94fe\u63a5\u5df2\u590d\u5236',
  referenceAdded: '\u5df2\u52a0\u5165\u53c2\u8003\u56fe',
  download: '\u4e0b\u8f7d',
  copyImage: '\u590d\u5236\u56fe\u7247',
  copyLink: '\u590d\u5236\u94fe\u63a5',
  useAsReference: '\u4f5c\u4e3a\u53c2\u8003\u56fe',
}

export function PreviewToolbar({ onDownload, onCopyImage, onCopyUrl, onUseAsReference, onNotice }: Props) {
  async function runAction(action: PreviewAction | undefined, fallback: string) {
    if (!action) return
    try {
      const result = await action()
      onNotice(typeof result === 'string' && result ? result : fallback)
    } catch (err) {
      onNotice(err instanceof Error ? err.message : TEXT.actionFailed)
    }
  }

  return (
    <div className="image-preview-toolbar">
      {onDownload ? <button type="button" onClick={() => void runAction(onDownload, TEXT.downloadStarted)}>{TEXT.download}</button> : null}
      {onCopyImage ? <button type="button" onClick={() => void runAction(onCopyImage, TEXT.imageCopied)}>{TEXT.copyImage}</button> : null}
      {onCopyUrl ? <button type="button" onClick={() => void runAction(onCopyUrl, TEXT.linkCopied)}>{TEXT.copyLink}</button> : null}
      {onUseAsReference ? <button type="button" onClick={() => void runAction(onUseAsReference, TEXT.referenceAdded)}>{TEXT.useAsReference}</button> : null}
    </div>
  )
}
