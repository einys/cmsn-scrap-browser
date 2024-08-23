import os
import platform
from flask import Flask, request, jsonify
from selenium import webdriver
from selenium.webdriver.firefox.service import Service
from selenium.webdriver.firefox.options import Options
from webdriver_manager.firefox import GeckoDriverManager
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from selenium.common.exceptions import TimeoutException, WebDriverException

import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# Firefox 옵션 설정 초기화
firefox_options = Options()

logger.info(" 🌟 System info: " + platform.system() + " " + platform.machine()) 

# Set the binary location
if platform.system() == "Darwin":  # macOS
    logger.info(" >> macOS system")
    firefox_options.binary_location = "/Applications/Firefox.app/Contents/MacOS/firefox"
elif platform.system() == "Linux":  # Ubuntu or Linux-based Docker container
    logger.info(" >> Linux system")
    firefox_options.binary_location = os.getenv("FIREFOX_BINARY_PATH", "/usr/bin/firefox")
else :
    logger.info(" >> Windows system")
    firefox_options.binary_location = "C:\\Program Files\\Mozilla Firefox\\firefox.exe"
    
firefox_options.add_argument("--headless")  # 브라우저 창을 열지 않고 실행

@app.route('/scrape-twitter', methods=['POST'])
def scrape_twitter():
    driver = None
    try:
        # GeckoDriver 로드
        try:
            logger.info(" 🦎 Initializing GeckoDriver service...")
            driver = webdriver.Firefox(options=firefox_options)
        except WebDriverException as e:
            logger.error("Failed to initialize GeckoDriver service: %s", str(e))
            return jsonify({"error": "Failed to initialize GeckoDriver service"}), 500
            
        data = request.json
        url = data.get("url")
        if not url:
            return jsonify({"error": "URL is required"}), 400
        
        # 페이지 로드 및 데이터 추출 로직은 동일
        driver.get(url)

        WebDriverWait(driver, 10).until(
            lambda d: d.execute_script("return window.__runPxScript !== undefined")
        )

        # Text content extraction
        text_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div[1]/div/section/div/div/div[1]/div/div/article/div/div/div[3]/div[1]/div"
        try:
            text = WebDriverWait(driver, 10).until(
                EC.presence_of_element_located((By.XPATH, text_xpath))
            ).text
        except TimeoutException:
            logger.error("Timeout loading text content")
            return jsonify({"error": "Timeout loading text content"}), 500
        
        # Image extraction
        image_xpath = "//img[contains(@src, 'https://pbs.twimg.com/media')]"
        try:
            image = WebDriverWait(driver, 10).until(
                EC.presence_of_element_located((By.XPATH, image_xpath))
            ).get_attribute('src')
        except TimeoutException:
            logger.warning("Timeout loading image.")
            image = None

        # Username extraction
        username_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div[1]/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[2]/div/div/div[1]/div/div/div[2]/div/div/a/div/span"
        try:
            username = WebDriverWait(driver, 10).until(
                EC.presence_of_element_located((By.XPATH, username_xpath))
            ).text
        except TimeoutException:
            logger.error("Timeout loading username")
            return jsonify({"error": "Timeout loading username"}), 500

        # Meta tag extraction
        try:
            meta_tag = WebDriverWait(driver, 10).until(
                EC.presence_of_element_located(
                    (By.XPATH, '//meta[@property="og:title"]')
                )
            ).get_attribute('content')
        except TimeoutException:
            logger.warning("Timeout loading meta tag.")
            meta_tag = None

        return jsonify({
            "text": text,
            "image": image,
            "username": username,
            "meta_tag": meta_tag
        })

    except Exception as e:
        logger.error("Internal server error: %s", str(e))
        return jsonify({"error": "Internal server error", "message": str(e)}), 500

    finally:
        if driver:
            logger.info("Closing the WebDriver...")
            driver.quit()

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=18081)
