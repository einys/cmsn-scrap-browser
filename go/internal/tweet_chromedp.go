package internal

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// ScrapeTweetChromedp는 chromedp로 공개 트윗 페이지에서 기본 정보를 긁어온다.
// parent는 재사용 컨텍스트(반복 크롤링용)를 권장한다.
func ScrapeTweetChromedp(parent context.Context, tweetURL string) (*TweetData, error) {
	if tweetURL == "" {
		return nil, errors.New("empty url")
	}
	if !strings.HasPrefix(tweetURL, "http") {
		tweetURL = "https://" + tweetURL
	}

	// 각 호출마다 타임아웃을 별도로 건다.
	ctx, cancel := context.WithTimeout(parent, 25*time.Second)
	defer cancel()

	// 결과 변수
	var (
		title, currentURL            string
		username, nickname, pfp, txt string
		ogTitle                      string
		imagesJSON, linksJSON        string
	)

	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"

	tasks := chromedp.Tasks{
		network.Enable(),
		emulation.SetUserAgentOverride(ua),
		emulation.SetLocaleOverride(),

		chromedp.Navigate(tweetURL),

		// 핵심 노드가 붙을 때까지 대기
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.WaitVisible("article", chromedp.ByQuery),

		chromedp.Title(&title),
		chromedp.Location(&currentURL),

		// @username
		chromedp.EvaluateAsDevTools(`(function(){
			const atexts = Array.from(document.querySelectorAll('article a'))
				.map(a => (a.textContent||'').trim());
			return atexts.find(t => t.includes('@')) || '';
		})()`, &username),

		// 닉네임(표시명)
		chromedp.EvaluateAsDevTools(`(function(){
			const el = document.querySelector('article [data-testid="User-Name"] span span')
				|| document.querySelector('article div[dir="ltr"] span span');
			return el ? el.textContent : '';
		})()`, &nickname),

		// 프로필 이미지 (대체 선택자 포함)
		chromedp.EvaluateAsDevTools(`(function(){
			return (document.querySelector('article img[src*="profile_images"]')
				|| document.querySelector('article img[alt][src*="pbs.twimg.com"]'))?.src || '';
		})()`, &pfp),

		// 본문 텍스트
		chromedp.EvaluateAsDevTools(`(function(){
			const el = document.querySelector('article div[data-testid="tweetText"]');
			return el ? el.innerText : '';
		})()`, &txt),

		// og:title (있으면 메타로 보완)
		chromedp.AttributeValue(`meta[property="og:title"]`, "content", &ogTitle, nil),

		// 이미지들
		chromedp.EvaluateAsDevTools(`JSON.stringify(
			Array.from(document.querySelectorAll('article img[src*="pbs.twimg.com/media"]'))
				.map(i => i.src)
		)`, &imagesJSON),

		// 링크들: 절대/상대/href/text 안의 URL 모두 수집(Set으로 중복 제거)
		chromedp.EvaluateAsDevTools(`(function(){
			const out = new Set();
			for (const a of document.querySelectorAll('article a[href]')) {
				const href = a.getAttribute('href') || '';
				if (href.startsWith('http')) out.add(href);
				else if (href.startsWith('/')) out.add('https://x.com' + href);

				const txt = (a.textContent||'').trim();
				const m = txt.match(/[a-zA-Z]+:\/\/[^\s]+/);
				if (m) out.add(m[0]);
			}
			return JSON.stringify(Array.from(out));
		})()`, &linksJSON),
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		// 디버깅용 스크린샷 남기기 (선택)
		var png []byte
		if ssErr := chromedp.Run(ctx, chromedp.CaptureScreenshot(&png)); ssErr == nil && len(png) > 0 {
			_ = os.WriteFile("tweet_error.png", png, 0644)
		}
		return nil, err
	}

	// JSON -> slice
	var images, links []string
	_ = json.Unmarshal([]byte(imagesJSON), &images)
	_ = json.Unmarshal([]byte(linksJSON), &links)

	// 메타 보정
	metaTitle := strings.TrimSpace(ogTitle)
	if metaTitle == "" {
		metaTitle = title
	}

	return &TweetData{
		Text:           strings.ReplaceAll(txt, "\n", " "),
		Images:         images,
		Username:       username,
		UserNickname:   nickname,
		UserProfileImg: pfp,
		MetaTag:        metaTitle,
		Links:          links,
	}, nil
}
