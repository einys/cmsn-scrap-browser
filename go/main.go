package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	seleniumPath     = ""                               // selenium-server-standalone jar 경로 (안 써도 됨: chromedriver만 쓸 경우)
	chromeDriverPath = "/opt/homebrew/bin/chromedriver" // chromedriver 위치
	port             = 9515
)

type TweetData struct {
	Text           string   `json:"text"`
	Images         []string `json:"images"`
	Username       string   `json:"username"`
	UserNickname   string   `json:"user_nickname"`
	UserProfileImg string   `json:"user_profile_img"`
	MetaTag        string   `json:"meta_tag"`
	Link           string   `json:"link"`
}

func main() {
	// Chrome 옵션 설정
	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		// "--headless", // 브라우저 UI 없이 실행. 테스트 시 주석 해제
		"--disable-gpu",
		"--no-sandbox",
	}})

	// WebDriver 실행 (selenium server 없이 chromedriver로 바로 실행)
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	if err != nil {
		log.Fatalf("Error starting the ChromeDriver server: %v", err)
	}
	defer service.Stop()

	// WebDriver 연결
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Fatalf("Error connecting to WebDriver: %v", err)
	}
	defer wd.Quit()

	tweet, err := ScrapeTweet(wd, "https://x.com/xxylolo/status/1903111391026012368")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", tweet)
}

func findTextByXPath(wd selenium.WebDriver, xpath string) string {
	elem, err := wd.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		return ""
	}
	text, err := elem.Text()
	if err != nil {
		return ""
	}
	return text
}

func findAttrByXPath(wd selenium.WebDriver, xpath, attr string) string {
	elem, err := wd.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		return ""
	}
	val, err := elem.GetAttribute(attr)
	if err != nil {
		return ""
	}
	return val
}

func ScrapeTweet(wd selenium.WebDriver, url string) (*TweetData, error) {
	err := wd.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to load URL: %w", err)
	}

	time.Sleep(7 * time.Second) // JS 로딩 대기

	// === Username ===
	username := findTextByXPath(wd, `/html/body/div[1]/div/div/div[2]/main/div/div/div/div[1]/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[2]/div/div/div[1]/div/div/div[2]/div/div/a/div/span`)

	// === Nickname ===
	nickname := findTextByXPath(wd, `/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[2]/div/div/div[1]/div/div/div[1]/div/a/div/div[1]/span/span`)

	// === Profile Image ===
	profileImg := findAttrByXPath(wd, `/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[1]/div[1]/div/div/div/div[2]/div/div[2]/div/a/div[3]/div/div[2]/div/img`, "src")

	// === Meta Tag ===
	metaTag := findAttrByXPath(wd, `//meta[@property='og:title']`, "content")

	// === Tweet Text ===
	text := findTextByXPath(wd, `/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[3]/div[1]/div/div`)

	// === Images (all) ===
	var images []string
	imgElements, _ := wd.FindElements(selenium.ByXPATH, `//img[contains(@src, 'https://pbs.twimg.com/media')]`)
	for _, img := range imgElements {
		src, _ := img.GetAttribute("src")
		images = append(images, src)
	}

	// === External Link ===
	contentLink := ""
	linkXPath := `/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div[1]/div/article/div/div/div[3]/div[1]/div/div/a[1]`

	linkElem, err := wd.FindElement(selenium.ByXPATH, linkXPath)
	if err == nil {
		linkText, _ := linkElem.Text()
		re := regexp.MustCompile(`[a-zA-Z0-9./-]+\.[a-zA-Z0-9/-]+`)
		match := re.FindString(linkText)
		if match != "" {
			contentLink = "https://" + match
		}
	}

	return &TweetData{
		Text:           text,
		Images:         images,
		Username:       username,
		UserNickname:   nickname,
		UserProfileImg: profileImg,
		MetaTag:        metaTag,
		Link:           contentLink,
	}, nil
}
