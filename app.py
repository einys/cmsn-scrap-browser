import os
import platform
from flask import Flask, request, jsonify
from selenium import webdriver
from selenium.webdriver.firefox.service import Service
from selenium.webdriver.firefox.options import Options
from webdriver_manager.firefox import GeckoDriverManager
from webdriver_manager.core.os_manager import OperationSystemManager
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By

app = Flask(__name__)


# Firefox 옵션 설정
firefox_options = Options()

# Set the binary location
if platform.system() == "Darwin":  # macOS
    firefox_options.binary_location = "/Applications/Firefox.app/Contents/MacOS/firefox"
elif platform.system() == "Linux":  # Ubuntu or Linux-based Docker container
    firefox_options.binary_location = os.getenv("FIREFOX_BINARY_PATH", "/usr/bin/firefox")
    
firefox_options.add_argument("--headless")  # 브라우저 창을 열지 않고 실행

os_manager = OperationSystemManager("linux_aarch64")
# GeckoDriver 설정
service = Service(GeckoDriverManager(version="v0.35.0", os_system_manager=os_manager).install())
driver = webdriver.Firefox(service=service, options=firefox_options)

@app.route('/scrape-twitter', methods=['POST'])
def scrape_twitter():
    try:
        data = request.json
        url = data.get("url")
        if not url:
            return jsonify({"error": "URL is required"}), 400
        

        # 페이지 로드 및 데이터 추출 로직은 동일
        driver.get(url)

        WebDriverWait(driver, 10).until(
            lambda d: d.execute_script("return window.__runPxScript !== undefined")
        )

        text_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div[1]/div/section/div/div/div[1]/div/div/article/div/div/div[3]/div[1]/div"
        text = WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.XPATH, text_xpath))
        )

        image_xpath = "//img[contains(@src, 'https://pbs.twimg.com/media')]"
        image = WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.XPATH, image_xpath))
        )

        username_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div[1]/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[2]/div/div/div[1]/div/div/div[2]/div/div/a/div/span"
        username = WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.XPATH, username_xpath))
        )

        meta_tag = WebDriverWait(driver, 10).until(
            EC.presence_of_element_located(
                (By.XPATH, '//meta[@property="og:title"]')
            )
        )

        meta_content = meta_tag.get_attribute("content")
        text_content = text.text
        image_url = image.get_attribute('src')
        username_content = username.text

        driver.quit()

        return jsonify({
            "meta_content": meta_content,
            "text_content": text_content,
            "image_url": image_url,
            "username_content": username_content
        })

    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=18081)
