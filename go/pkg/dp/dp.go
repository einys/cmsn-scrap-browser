package dp

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// Metadata 구조체
type Metadata struct {
	Title       string
	Description string
	Image       string
	URL         string
}

// ScrapeMetadata 크롤링 함수
func ScrapeMetadata(targetURL string) (*Metadata, error) {

	log.Println("🔍 크롤링 대상 URL:", targetURL)

	// 크롬 context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 타임아웃
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var title, desc, image string

	// chromedp 실행
	tasks := chromedp.Tasks{
		chromedp.Navigate(targetURL),

		// 일반 <title>
		chromedp.Title(&title),

		// <meta name="description">
		chromedp.AttributeValue(`meta[name="description"]`, "content", &desc, nil),

		// <meta property="og:image">
		chromedp.AttributeValue(`meta[property="og:image"]`, "content", &image, nil),
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return nil, err
	}

	// 정리
	meta := &Metadata{
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(desc),
		Image:       strings.TrimSpace(image),
		URL:         targetURL,
	}

	return meta, nil
}
