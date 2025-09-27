package internal

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tebeka/selenium"

	"github.com/einys/cmsn-scraper/lib"
)

// MetaData : 메타데이터 결과 구조체
type MetaData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"img"`
	URL         string `json:"url"`
}

// ScrapeMeta : 일반 페이지의 메타데이터 스크래핑
func ScrapeMeta(wd selenium.WebDriver, pageURL string) (*MetaData, error) {
	log.Printf("📥 Scraping meta: %s", pageURL)
	startTime := time.Now()
	meta := &MetaData{URL: pageURL}

	// 페이지 로딩
	if err := wd.Get(pageURL); err != nil {
		return nil, err
	}
	if err := WaitForPageLoad(wd, 10); err != nil {
		return nil, fmt.Errorf("failed to wait for page load: %v", err)
	}
	log.Printf("✅ Page loaded in %v", time.Since(startTime))

	// === Title ===
	if strings.Contains(pageURL, ".notion.") {
		log.Printf("🔍 Handling Notion title...")
		wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			script := `return document.querySelector(".notion-page-content")?.innerText;`
			text, err := wd.ExecuteScript(script, nil)
			return text != nil && text.(string) != "", err
		}, 10*time.Second)

		titleJS, err := wd.ExecuteScript("return document.title;", nil)
		if err == nil {
			meta.Title = titleJS.(string)
		}
	} else {
		elem, err := wd.FindElement(selenium.ByXPATH, `//meta[@property="og:title"]`)
		if err != nil {
			elem, err = wd.FindElement(selenium.ByXPATH, `//head/title`)
			if err == nil {
				meta.Title, _ = elem.Text()
			}
		} else {
			meta.Title, _ = elem.GetAttribute("content")
		}
	}
	log.Printf("🏷 Title: %s", meta.Title)

	// === Image ===
	imgElem, err := wd.FindElement(selenium.ByXPATH, `//meta[@property="og:image"]`)
	if err != nil {
		imgElem, err = wd.FindElement(selenium.ByXPATH, `//meta[@name="image"]`)
	}
	if imgElem != nil {
		meta.Image, _ = imgElem.GetAttribute("content")
	}
	log.Printf("🖼 Image: %s", meta.Image)

	// === Description ===
	if strings.Contains(pageURL, ".notion.") {
		log.Printf("🔍 Handling Notion description...")
		wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			script := `return document.querySelector(".notion-page-content")?.innerText;`
			text, err := wd.ExecuteScript(script, nil)
			return text != nil && text.(string) != "", err
		}, 10*time.Second)

		descJS, err := wd.ExecuteScript(`return document.querySelector(".notion-page-content")?.innerText.slice(0, 200);`, nil)
		if err == nil {
			clean := lib.CleanText(descJS.(string))
			if len(clean) > 200 {
				clean = clean[:200] + "..."
			}
			meta.Description = clean
		}
	} else {
		descElem, err := wd.FindElement(selenium.ByXPATH, `//meta[@property="og:description"]`)
		if err != nil {
			descElem, err = wd.FindElement(selenium.ByCSSSelector, `meta[name="description"]`)
			if err == nil {
				meta.Description, _ = descElem.GetAttribute("content")
			}
		} else {
			meta.Description, _ = descElem.GetAttribute("content")
		}
	}
	log.Printf("📝 Description: %s", meta.Description)

	log.Printf("✅ Done scraping meta: %s (%v)", pageURL, time.Since(startTime))
	return meta, nil
}
