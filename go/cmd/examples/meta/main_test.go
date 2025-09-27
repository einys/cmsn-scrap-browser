package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScrapeTwitterAPI(t *testing.T) {
	req := httptest.NewRequest("GET", "/scrape?url=https://x.com/naeng2_/status/1903488320367403357", nil)
	w := httptest.NewRecorder()

	TweetScrapeHandler(w, req)

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

func TestScrapeMetaApi(t *testing.T) {
	req := httptest.NewRequest("GET", "/meta?url=https://grand-mistake-cc4.notion.site/1bd539ffbe2b80849b14c504a7509543", nil)
	w := httptest.NewRecorder()

	MetaHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	t.Logf("📄 응답 코드: %d", resp.StatusCode)
	t.Logf("📦 응답 본문: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("❌ 응답 실패: %d", resp.StatusCode)
	}

	if !strings.Contains(string(body), `"title":`) {
		t.Errorf("❌ 예상하는 필드 없음 (title)")
	}
}
