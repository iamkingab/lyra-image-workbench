import type { TaskResult } from '../types'

export function errorReasonLabel(result: TaskResult) {
  const reason = normalizeErrorReason(result)
  return [reason.chinese, reason.code, reason.english].filter(Boolean).join(' / ')
}

export function normalizeErrorReason(result: TaskResult) {
  if (result.errorText || result.errorCode || result.errorEnglish) {
    return {
      chinese: result.errorText || '任务执行失败',
      code: result.errorCode || result.statusCode,
      english: result.errorEnglish || result.status,
    }
  }
  const raw = (result.error || '').trim().toLowerCase()
  if (raw.includes('unexpected eof')) {
    return { chinese: '上游响应提前结束', code: 'E_UPSTREAM_EOF', english: 'upstream_response_truncated' }
  }
  if (raw === 'eof' || raw.includes('empty response')) {
    return { chinese: '上游返回空响应', code: 'E_UPSTREAM_EMPTY', english: 'upstream_empty_response' }
  }
  if (raw.includes('context deadline exceeded') || raw.includes('timeout')) {
    return { chinese: '上游请求超时', code: 'E_UPSTREAM_TIMEOUT', english: 'upstream_timeout' }
  }
  const http = raw.match(/http\s+(\d{3})/)
  if (http) {
    return { chinese: `上游接口返回 HTTP ${http[1]}`, code: `E_UPSTREAM_HTTP_${http[1]}`, english: `upstream_http_${http[1]}` }
  }
  if (raw.includes('invalid character') || raw.includes('cannot unmarshal') || raw.includes('bad json')) {
    return { chinese: '上游返回不是有效 JSON', code: 'E_UPSTREAM_BAD_JSON', english: 'upstream_bad_json' }
  }
  if (raw.includes('没有返回可用图片') || raw.includes('no usable image')) {
    return { chinese: '上游没有返回可用图片', code: 'E_UPSTREAM_NO_IMAGE', english: 'upstream_no_image' }
  }
  return { chinese: '任务执行失败', code: result.statusCode || 'J500', english: result.status || 'failed' }
}
