package internal

import (
	"errors"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

func FindTextByXPath(wd selenium.WebDriver, xpath string) string {
	elem, err := wd.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		log.Printf("❌ Failed to find element: %v", err)
		return ""
	}
	text, _ := elem.Text()
	return text
}

func FindAttrByXPath(wd selenium.WebDriver, xpath, attr string) string {
	elem, err := wd.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		log.Printf("❌ Failed to find element: %v", err)
		return ""
	}
	val, _ := elem.GetAttribute(attr)
	return val
}

func WaitForPageLoad(wd selenium.WebDriver, timeoutSeconds int) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	end := time.Now().Add(timeout)
	for {
		state, err := wd.ExecuteScript("return document.readyState", nil)
		if err == nil && state == "complete" {
			return nil
		}
		if time.Now().After(end) {
			return errors.New("timeout waiting for page load")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
