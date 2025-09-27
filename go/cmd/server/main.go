package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/einys/cmsn-scraper/internal"
)

var (
	ENGINE = "chromedp" // 기본값
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 환경변수로 엔진 설정
	ENGINE = os.Getenv("SCRAPER_ENGINE")
	log.Println("🛠️  Using SCRAPER_ENGINE:", ENGINE)

	// 서버 시작
	http.HandleFunc("/scrape-twitter", tweetHandler)
	http.HandleFunc("/meta", metaHandler)
	log.Println("🚀 Server running on http://localhost:18081")
	log.Fatal(http.ListenAndServe(":18081", nil))

}

func tweetHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing 'url'", http.StatusBadRequest)
		return
	}

	log.Println("🐦 트윗 스크래핑 요청 URL:", url)

	if ENGINE == "chromedp" {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		data, err := internal.ScrapeTweetChromedp(ctx, url)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	// 기본: selenium
	wd, quit, err := internal.InitWebDriver()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer quit()
	defer wd.Quit()

	data, err := internal.ScrapeTweet(wd, url)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func normalizeURL(u string) string {
	u = strings.TrimSpace(u)
	if u == "" {
		return u
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "https://" + u
	}
	return u
}

func metaHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing 'url'", http.StatusBadRequest)
		return
	}

	url = normalizeURL(url)
	if url == "" {
		http.Error(w, "Invalid 'url'", http.StatusBadRequest)
		return
	}

	log.Println("🌐 메타데이터 스크래핑 요청 URL:", url)

	if ENGINE == "chromedp" {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		data, err := internal.ScrapeMetaChromedp(ctx, url)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	// 기본: selenium
	wd, quit, err := internal.InitWebDriver()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer quit()
	defer wd.Quit()

	data, err := internal.ScrapeMeta(wd, url)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(data)
}
