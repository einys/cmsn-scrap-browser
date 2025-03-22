package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScrapeAPI(t *testing.T) {
	req := httptest.NewRequest("GET", "/scrape?url=https://x.com/xxylolo/status/1903111391026012368", nil)
	w := httptest.NewRecorder()

	ScrapeHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	t.Logf("📄 응답 코드: %d", resp.StatusCode)
	t.Logf("📦 응답 본문: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("❌ 응답 실패: %d", resp.StatusCode)
	}

	if !strings.Contains(string(body), `"username":`) {
		t.Errorf("❌ 예상하는 필드 없음 (username)")
	}
}
