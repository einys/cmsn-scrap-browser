package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// context 생성
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 실행 타임아웃 설정
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var title string
	url := "https://example.com"

	// 크롬 동작 시퀀스
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Title(&title), // 페이지 타이틀 가져오기
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Page title:", title)
}
