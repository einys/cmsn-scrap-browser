package internal

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

type TweetData struct {
	Text           string   `json:"text"`
	Images         []string `json:"images"`
	Username       string   `json:"username"`
	UserNickname   string   `json:"user_nickname"`
	UserProfileImg string   `json:"user_profile_img"`
	MetaTag        string   `json:"meta_tag"`
	Links          []string `json:"links"`
}

func ScrapeTweet(wd selenium.WebDriver, url string) (*TweetData, error) {
	log.Printf("📥 Scraping tweet: %s", url)
	if err := wd.Get(url); err != nil {
		return nil, fmt.Errorf("failed to load URL: %w", err)
	}

	// 대기
	for i := 0; i < 20; i++ {
		if _, err := wd.FindElement(selenium.ByCSSSelector, "article"); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if _, err := wd.FindElement(selenium.ByCSSSelector, "article"); err != nil {
		src, _ := wd.PageSource()
		_ = os.WriteFile("page.html", []byte(src), 0644)
		return nil, fmt.Errorf("failed to find <article>: %w", err)
	}

	username := FindTextByXPath(wd, `//article//a[starts-with(@href, "/") and contains(., "@")]`)
	nickname := FindTextByXPath(wd, `//article//div[@dir="ltr"]//span/span`)
	profileImg := FindAttrByXPath(wd, `//article//img[contains(@src, 'profile_images')]`, "src")
	metaTag := FindAttrByXPath(wd, `//meta[@property='og:title']`, "content")
	text := FindTextByXPath(wd, `//article//div[@data-testid="tweetText"]`)

	// 이미지
	var images []string
	imgElements, _ := wd.FindElements(selenium.ByXPATH, `//img[contains(@src, 'https://pbs.twimg.com/media')]`)
	for _, img := range imgElements {
		src, _ := img.GetAttribute("src")
		images = append(images, src)
	}

	// 링크
	var links []string
	if linkElems, err := wd.FindElements(selenium.ByXPATH, `//article//a`); err == nil {
		re := regexp.MustCompile(`[a-zA-Z0-9/-]*\.[a-zA-Z0-9/-]+[a-zA-Z0-9./-]*`)
		for _, el := range linkElems {
			linkText, _ := el.Text()
			if match := re.FindString(linkText); match != "" {
				links = append(links, "https:"+match)
			}
		}
	}

	return &TweetData{
		Text:           strings.ReplaceAll(text, "\n", " "),
		Images:         images,
		Username:       username,
		UserNickname:   nickname,
		UserProfileImg: profileImg,
		MetaTag:        strings.ReplaceAll(metaTag, "\n", " "),
		Links:          links,
	}, nil
}
