package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var (
	// chromedriver 위치
	chromeDriverPath = "/usr/bin/chromedriver"
	port             = 9515
)

func init() {
	// macOS 환경(테스트 환경)인 경우 path를 /opt/homebrew/bin/chromedriver 로 설정
	if runtime.GOOS == "darwin" {
		log.Println("🍏 macOS detected. Setting chromedriver path to /opt/homebrew/bin/chromedriver")
		chromeDriverPath = "/opt/homebrew/bin/chromedriver"
	} else {
		log.Println("🐧 Linux detected. Using default chromedriver path.")
	}
}

type TweetData struct {
	Text           string   `json:"text"`
	Images         []string `json:"images"`
	Username       string   `json:"username"`
	UserNickname   string   `json:"user_nickname"`
	UserProfileImg string   `json:"user_profile_img"`
	MetaTag        string   `json:"meta_tag"`
	Links          []string `json:"links"` // ✅ 여러 링크
}

func main() {

	http.HandleFunc("/scrape", ScrapeHandler)
	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	tweetURL := r.URL.Query().Get("url")
	if tweetURL == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	// Chrome options
	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--headless=new",
		"--disable-gpu",
		"--no-sandbox",
		"--window-size=1280,1024",
		"--disable-dev-shm-usage", // 추가된 플래그
		"--lang=ko-KR,ko",
		"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	}})

	// Start chromedriver
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start ChromeDriver: %v", err), http.StatusInternalServerError)
		return
	}
	defer service.Stop()

	// Connect to WebDriver
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to WebDriver: %v", err), http.StatusInternalServerError)
		return
	}
	defer wd.Quit()

	// Scrape tweet
	tweetData, err := ScrapeTweet(wd, tweetURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to scrape tweet: %v", err), http.StatusInternalServerError)
		return
	}

	// Return as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweetData)
}

func findTextByXPath(wd selenium.WebDriver, xpath string) string {

	elem, err := wd.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		log.Printf("❌ Error: Failed to find element: %v", err)
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
		log.Printf("❌ Error: Failed to find element: %v", err)
		return ""
	}
	val, err := elem.GetAttribute(attr)
	if err != nil {
		return ""
	}
	return val
}

func ScrapeTweet(wd selenium.WebDriver, url string) (*TweetData, error) {
	log.Printf("📥 크롤링 시작: %s", url)

	err := wd.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to load URL: %w", err)
	}

	// source, err := wd.PageSource()
	// if err != nil {
	// 	log.Printf("❌ Failed to get page source: %v\n", err)
	// } else {
	// 	log.Println("✅ Page loaded. First 2000 chars:")
	// 	log.Println(source[:2000])
	// }

	// 페이지 타이틀 출력
	title, _ := wd.Title()
	log.Printf("📄 Title: %s", title)

	// 페이지 로딩 대기
	for i := 0; i < 20; i++ {
		_, err := wd.FindElement(selenium.ByCSSSelector, "article")
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// 웹페이지 로딩 상태 확인
	_, err = wd.FindElement(selenium.ByCSSSelector, "article")
	if err != nil {
		log.Println("❌ Error: <article> 태그를 못 찾았어요. 아마 트윗이 안 보이거나 리디렉션 된 것 같아요. 페이지 소스를 저장할게요.")

		// 페이지 소스 가져오기 및 파일 저장
		source, err := wd.PageSource()
		if err != nil {
			log.Printf("페이지 소스 가져오기 실패: %v", err)
		} else {
			err = os.WriteFile("page.html", []byte(source), 0644)
			if err != nil {
				log.Printf("파일 저장 실패: %v", err)
			} else {
				log.Println("페이지 소스가 'page.html' 파일로 저장되었습니다. docker 환경인 경우 ./app/page.html 파일을 확인하세요.")
			}
		}
		return nil, fmt.Errorf("failed to find <article> element: %w", err)
	}
	currentURL, _ := wd.CurrentURL()
	log.Printf("🌐 현재 URL: %s\n", currentURL)

	// === Username ===
	log.Printf("🔍 Finding username...")
	username := findTextByXPath(wd, `//article//a[starts-with(@href, "/") and contains(., "@")]`)
	log.Printf("👤 Username: %s", username)
	// === Nickname ===
	log.Printf("🔍 Finding nickname...")
	nickname := findTextByXPath(wd, `//article//div[@dir="ltr"]//span/span`)
	log.Printf("👤 Nickname: %s", nickname)
	// === Profile Image ===
	log.Printf("🔍 Finding profile image...")
	profileImg := findAttrByXPath(wd, `/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[1]/div[1]/div/div/div/div[2]/div/div[2]/div/a/div[3]/div/div[2]/div/img`, "src")
	log.Printf("👤 Profile Image: %s", profileImg)
	// === Meta Tag ===
	log.Printf("🔍 Finding meta tag...")
	metaTag := findAttrByXPath(wd, `//meta[@property='og:title']`, "content")
	log.Printf("🏷 Meta Tag: %s", strings.ReplaceAll(metaTag, "\n", " "))
	// === Tweet Text ===
	log.Printf("🔍 Finding tweet text...")
	text := findTextByXPath(wd, `/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[3]/div[1]/div/div`)
	log.Printf("📝 Tweet Text: %s", strings.ReplaceAll(text, "\n", " "))

	// === Images (all) ===
	log.Printf("🔍 Finding images...")
	var images []string
	imgElements, _ := wd.FindElements(selenium.ByXPATH, `//img[contains(@src, 'https://pbs.twimg.com/media')]`)
	for _, img := range imgElements {
		src, _ := img.GetAttribute("src")
		images = append(images, src)
	}

	// === All links in tweet ===
	log.Printf("🔍 Finding links...")
	linkElems, err := wd.FindElements(selenium.ByXPATH, `//article//a`)
	var links []string

	if err == nil {

		re := regexp.MustCompile(`[a-zA-Z0-9/-]*\.[a-zA-Z0-9/-]+[a-zA-Z0-9./-]*`) // 링크 추출 정규식. 이걸 안 할 경우 결과  Links:[ Xylo @xxylolo #커미션 #rt http://kre.pe/nKn1 https://open.kakao.com/o/sr5J0Vmh  오후 4:48 · 2025년 3월 21일]
		// 💡 팁: href 를 쓰는 게 더 정확하긴 해
		// 일부 트윗은 <a> 태그에 텍스트가 없고, href 속성에만 URL이 있는 경우도 있어서:

		// href, err := el.GetAttribute("href")
		// if err == nil && strings.Contains(href, "t.co") {
		//     links = append(links, href)
		// }
		// 필요하면 text + href 조합으로도 만들 수 있어.

		for _, el := range linkElems {
			linkText, _ := el.Text()
			match := re.FindString(linkText)
			if match != "" {
				link := "https:" + match
				// 중복 제거 또는 필터링 필요 시 여기에 처리
				links = append(links, link)
			}
		}
	}

	log.Printf("🖼 Images: %v", images)
	log.Printf("🔗 Links: %v", links)
	log.Printf("✅ 크롤링 완료: %s", url)

	return &TweetData{
		Text:           text,
		Images:         images,
		Username:       username,
		UserNickname:   nickname,
		UserProfileImg: profileImg,
		MetaTag:        metaTag,
		Links:          links,
	}, nil
}
