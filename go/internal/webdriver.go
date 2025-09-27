package internal

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var (
	chromeDriverPath = "/usr/bin/chromedriver"
	chromiumPath     = "/usr/bin/chromium"
	port             = 9515
	myOS             = runtime.GOOS
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if myOS == "darwin" {
		chromeDriverPath = "/opt/homebrew/bin/chromedriver"
		chromiumPath = "/opt/homebrew/bin/chromium"
		log.Println("🍏 macOS detected. Adjusting chromedriver path:", chromeDriverPath)
	}
}

func InitWebDriver() (selenium.WebDriver, func(), error) {
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeArgs := []string{
		"--window-size=1280,1024",
		"--disable-dev-shm-usage",
		"--lang=ko-KR,ko",
		"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	}
	if myOS != "darwin" {
		log.Println("🐧 Linux: running in headless mode.")
		chromeArgs = append([]string{"--headless", "--disable-gpu", "--no-sandbox"}, chromeArgs...)
	}
	chromeCaps := chrome.Capabilities{
		Path: chromiumPath,
		Args: chromeArgs,
	}
	caps.AddChrome(chromeCaps)

	logFile, _ := os.Create("/tmp/chromedriver.log")
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, selenium.Output(logFile))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start ChromeDriver: %v", err)
	}
	quitFunc := func() { service.Stop() }

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		quitFunc()
		return nil, nil, fmt.Errorf("failed to connect to WebDriver: %v", err)
	}
	return wd, quitFunc, nil
}
