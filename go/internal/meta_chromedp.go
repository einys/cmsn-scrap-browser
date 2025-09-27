package internal

import (
	"context"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// ScrapeMetaChromedp : chromedp 버전
func ScrapeMetaChromedp(parent context.Context, pageURL string) (*MetaData, error) {
	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	var title, desc, image string

	tasks := chromedp.Tasks{
		chromedp.Navigate(pageURL),
		chromedp.WaitReady("body", chromedp.ByQuery),

		chromedp.Title(&title),
		chromedp.AttributeValue(`meta[name="description"]`, "content", &desc, nil),
		chromedp.AttributeValue(`meta[property="og:image"]`, "content", &image, nil),
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return nil, err
	}

	return &MetaData{
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(desc),
		Image:       strings.TrimSpace(image),
		URL:         pageURL,
	}, nil
}
