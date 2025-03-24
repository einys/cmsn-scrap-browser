from flask import Flask, request, jsonify
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By

app = Flask(__name__)

@app.route('/scrape-twitter', methods=['POST'])
def scrape_twitter():
    try:
        # Request에서 URL 가져오기
        data = request.json
        url = data.get("url")
        if not url:
            return jsonify({"error": "URL is required"}), 400

        # 크롬 옵션 설정
        chrome_options = Options()
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--headless")  # 브라우저 창을 열지 않고 실행
        chrome_options.add_argument("--disable-gpu")
        chrome_options.add_argument("--disable-dev-shm-usage")
        chrome_options.add_argument("--remote-debugging-port=9222")  # this
        chrome_options.add_argument(
            "user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
        )

        # ChromeDriver 설정
        service = Service(ChromeDriverManager().install())
        driver = webdriver.Chrome(service=service, options=chrome_options)

        # 트위터 페이지 로드
        driver.get(url)

        # 필요한 스크립트 파일들이 로드되고 실행될 때까지 기다림
        WebDriverWait(driver, 10).until(
            lambda d: d.execute_script("return window.__runPxScript !== undefined")
        )

        # XPath로 요소가 로딩될 때까지 최대 10초 대기
        text_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div[1]/div/section/div/div/div[1]/div/div/article/div/div/div[3]/div[1]/div"   # 실제 트위터 텍스트가 있는 XPath
        text = WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.XPATH, text_xpath))
        )

        # XPath로 src가 "https://pbs.twimg.com/media"로 시작하는 이미지 요소를 찾기
        image_xpath = "//img[contains(@src, 'https://pbs.twimg.com/media')]"

        # 이미지가 로드될 때까지 최대 10초 대기
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

        # 메타 태그의 content 속성 추출
        meta_content = meta_tag.get_attribute("content")

        # 텍스트 요소에서 텍스트 추출
        text_content = text.text

        # 이미지 요소에서 src 속성 값 추출
        image_url = image.get_attribute('src')

        # 유저네임 요소에서 텍스트 추출
        username_content = username.text

        # 드라이버 종료
        driver.quit()

        # JSON 응답 생성
        return jsonify({
            "meta_content": meta_content,
            "text_content": text_content,
            "image_url": image_url,
            "username_content": username_content
        })

    except Exception as e:
        # 에러 발생 시 JSON 응답 생성
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=18081)
