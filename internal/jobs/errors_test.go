package jobs

import "testing"

func TestErrorMetaMapsUnexpectedEOF(t *testing.T) {
	meta := ErrorMeta("unexpected EOF")
	if meta.Code != "E_UPSTREAM_EOF" || meta.English != "upstream_response_truncated" || meta.Chinese != "上游响应提前结束" {
		t.Fatalf("unexpected meta: %+v", meta)
	}
	result := NewResult(0, StatusFailed, "unexpected EOF")
	if result.ErrorCode != meta.Code || result.ErrorText != meta.Chinese || result.ErrorEnglish != meta.English {
		t.Fatalf("result error fields not populated: %+v", result)
	}
}

func TestErrorMetaMapsHTTPStatus(t *testing.T) {
	meta := ErrorMeta("上游请求失败：HTTP 524：timeout")
	if meta.Code != "E_UPSTREAM_HTTP_524" || meta.English != "upstream_http_524" || meta.Chinese != "上游接口返回 HTTP 524" {
		t.Fatalf("unexpected meta: %+v", meta)
	}
}
