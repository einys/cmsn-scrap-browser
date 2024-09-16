import os
import platform
import re
import traceback
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
import datetime

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# Firefox 옵션 설정 초기화
firefox_options = Options()

logger.info(" 🌟 System info: " + platform.system() + " " + platform.machine())

# print current time with timezone and icon
logger.info(" 🚀 Starting the Flask app...")
logger.info(" 🦊 Initializing Firefox WebDriver...")
current_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
timezone = datetime.datetime.now().astimezone().tzinfo
logger.info(f" ⏰ Current time: {current_time} Timezone: {timezone}")

WAITING_TIME_SEC = 4
WAITING_TIME_SEC_LINK = 1

# Set the binary location
if platform.system() == "Darwin":  # macOS
    logger.info(" >> macOS system")
    firefox_options.binary_location = "/Applications/Firefox.app/Contents/MacOS/firefox"
elif platform.system() == "Linux":  # Ubuntu or Linux-based Docker container
    logger.info(" >> Linux system")
    firefox_options.binary_location = os.getenv(
        "FIREFOX_BINARY_PATH", "/usr/bin/firefox")
else:
    logger.info(" >> Windows system")
    firefox_options.binary_location = "C:\\Program Files\\Mozilla Firefox\\firefox.exe"

# firefox_options.add_argument("--headless")  # 브라우저 창을 열지 않고 실행

# GeckoDriver 로드
try:
    logger.info(" 🦎 Initializing GeckoDriver service...")
    driver = webdriver.Firefox(options=firefox_options)
    logger.info(" ✅ GeckoDriver service initialized successfully.")
except WebDriverException as e:
    logger.error("Failed to initialize GeckoDriver service: %s", str(e))
    exit(1)


@app.route('/scrape-twitter', methods=['POST'])
def scrape_twitter():

    try:

        data = request.json
        url = data.get("url")
        if not url:
            return jsonify({"error": "URL is required"}), 400

        logger.warning(
            f"🔗 Scraping the URL(by POST request! depreciated): {url}")

        # log request ip address
        logger.info(
            f"🌐 Who is requesting? Requesting IP address: {request.remote_addr}")

        # 페이지 로드 및 데이터 추출 로직은 동일
        driver.get(url)

        # WebDriverWait(driver, 10).until(
        #     lambda d: d.execute_script("return window.__runPxScript !== undefined")
        # )

        # Text content extraction
        text_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div[1]/div/section/div/div/div[1]/div/div/article/div/div/div[3]/div[1]/div"
        try:
            text = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located((By.XPATH, text_xpath))
            ).text
        except TimeoutException:
            logger.error("Timeout loading text content")
            return jsonify({"error": "Timeout loading text content"}), 500

        # Image extraction
        image_xpath = "//img[contains(@src, 'https://pbs.twimg.com/media')]"
        try:
            image = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located((By.XPATH, image_xpath))
            ).get_attribute('src')
        except TimeoutException:
            logger.warning("Timeout loading image.")
            image = None

        # Username extraction
        username_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div[1]/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[2]/div/div/div[1]/div/div/div[2]/div/div/a/div/span"
        try:
            username = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located((By.XPATH, username_xpath))
            ).text
        except TimeoutException:
            logger.error("Timeout loading username")
            return jsonify({"error": "Timeout loading username"}), 500

        user_nickname_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[2]/div/div/div[1]/div/div/div[1]/div/a/div/div[1]/span/span"
        try:
            user_nickname = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located((By.XPATH, user_nickname_xpath))
            ).text
        except TimeoutException:
            logger.warning("Timeout loading user_nickname")
            user_nickname = None

        # User profile image extraction
        user_profile_img = None
        user_profile_img_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[1]/div[1]/div/div/div/div[2]/div/div[2]/div/a/div[3]/div/div[2]/div/img"
        try:
            user_profile_img = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located(
                    (By.XPATH, user_profile_img_xpath))
            ).get_attribute('src')
        except TimeoutException:
            logger.warning("Timeout loading user profile image.")

        # Meta tag extraction
        meta_tag = None
        try:
            meta_tag = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located(
                    (By.XPATH, '//meta[@property="og:title"]')
                )
            ).get_attribute('content')
        except TimeoutException:
            logger.warning("Timeout loading meta tag.")

        link_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[3]/div[2]/div/a"
        content_link = None
        try:
            content_link = WebDriverWait(driver, WAITING_TIME_SEC_LINK).until(
                EC.presence_of_element_located((By.XPATH, link_xpath))
            ).text

            # 링크만 추출하는 패턴
            pattern = r"(?=[a-zA-Z0-9/-]*\.[a-zA-Z0-9/-])[a-zA-Z0-9./-]+"

            # 패턴에 맞는 부분 추출
            match = re.search(pattern, content_link)
            content_link = match.group(0)
            content_link = "https://" + content_link
        except TimeoutException:
            logger.warning("Timeout loading link.")
        except Exception as e:
            print(f"General exception occurred: {e}")

        return jsonify({
            "text": text,
            "image": image,
            "username": username,
            "user_nickname": user_nickname,
            "user_profile_img": user_profile_img,
            "meta_tag": meta_tag,
            "link": content_link
        })

    except Exception as e:
        print("❌ Error occured ")
        print("My log app print exception:", e)
        # 에러 핸들러가 자동으로 호출되므로, 별도의 처리 없이도 됩니다.
        # Return the error details in the response
        return jsonify({
            "message": str(e)
        }), 500

    # finally:
    #     if driver:
    #         logger.info("Closing the WebDriver...")
    #         driver.quit()


@app.route('/scrape-twitter', methods=['GET'])
def scrape_twitter_get():
    try:
        url = request.args.get("url")
        if not url:
            return jsonify({"error": "URL is required"}), 400

        # save current time and log it
        init_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        logger.info(f"🚀 Starting the scraping process at {init_time}")

        # print url
        logger.info(f"🔗 Scraping the URL: {url}")

        # 페이지 로드 및 데이터 추출 로직은 동일
        driver.get(url)

        # WebDriverWait(driver, 10).until(
        #     lambda d: d.execute_script("return window.__runPxScript !== undefined")
        # )

        # Text content extraction
        text_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div[1]/div/section/div/div/div[1]/div/div/article/div/div/div[3]/div[1]/div"
        try:
            text = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located((By.XPATH, text_xpath))
            ).text
        except TimeoutException:
            logger.error("Timeout loading text content")
            return jsonify({"error": "Timeout loading text content"}), 500

        # calculate running time until now and log it as seconds
        running_time = datetime.datetime.now() - datetime.datetime.strptime(init_time,
                                                                            "%Y-%m-%d %H:%M:%S")
        logger.info(
            f"⏱️  Text has been scraped. Running time: {running_time.total_seconds()} sec")

        # Image extraction
        image_xpath = "//img[contains(@src, 'https://pbs.twimg.com/media')]"
        try:
            image = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located((By.XPATH, image_xpath))
            ).get_attribute('src')
        except TimeoutException:
            logger.warning("Timeout loading image.")
            image = None

        # calculate running time until now and log it as seconds
        running_time = datetime.datetime.now() - datetime.datetime.strptime(init_time,
                                                                            "%Y-%m-%d %H:%M:%S")
        logger.info(
            f"⏱️  Image has been scraped. Running time: {running_time.total_seconds()} sec")

        # Username extraction
        username_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div[1]/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[2]/div/div/div[1]/div/div/div[2]/div/div/a/div/span"
        try:
            username = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located((By.XPATH, username_xpath))
            ).text
        except TimeoutException:
            logger.error("Timeout loading username")
            return jsonify({"error": "Timeout loading username"}), 500

        # calculate running time until now and log it as seconds
        running_time = datetime.datetime.now() - datetime.datetime.strptime(init_time,
                                                                            "%Y-%m-%d %H:%M:%S")
        logger.info(
            f"⏱️  Username has been scraped. Running time: {running_time.total_seconds()} sec")

        user_nickname_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[2]/div/div/div[1]/div/div/div[1]/div/a/div/div[1]/span/span"
        try:
            user_nickname = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located((By.XPATH, user_nickname_xpath))
            ).text
        except TimeoutException:
            logger.warning("Timeout loading user_nickname")
            user_nickname = None

        # calculate running time until now and log it as seconds
        running_time = datetime.datetime.now() - datetime.datetime.strptime(init_time,
                                                                            "%Y-%m-%d %H:%M:%S")
        logger.info(
            f"⏱️  User nickname has been scraped. Running time: {running_time.total_seconds()} sec")

        # User profile image extraction
        user_profile_img = None
        user_profile_img_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[2]/div[1]/div[1]/div/div/div/div[2]/div/div[2]/div/a/div[3]/div/div[2]/div/img"
        try:
            user_profile_img = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located(
                    (By.XPATH, user_profile_img_xpath))
            ).get_attribute('src')
        except TimeoutException:
            logger.warning("Timeout loading user profile image.")

        # calculate running time until now and log it as seconds
        running_time = datetime.datetime.now() - datetime.datetime.strptime(init_time,
                                                                            "%Y-%m-%d %H:%M:%S")
        logger.info(
            f"⏱️  User profile image has been scraped. Running time: {running_time.total_seconds()} sec")

        # Meta tag extraction
        meta_tag = None
        try:
            meta_tag = WebDriverWait(driver, WAITING_TIME_SEC).until(
                EC.presence_of_element_located(
                    (By.XPATH, '//meta[@property="og:title"]')
                )
            ).get_attribute('content')
        except TimeoutException:
            logger.warning("Timeout loading meta tag.")

        # calculate running time until now and log it as seconds
        running_time = datetime.datetime.now() - datetime.datetime.strptime(init_time,
                                                                            "%Y-%m-%d %H:%M:%S")
        logger.info(
            f"⏱️  Meta tag has been scraped. Running time: {running_time.total_seconds()} sec")

        link_xpath = "/html/body/div[1]/div/div/div[2]/main/div/div/div/div/div/section/div/div/div[1]/div/div/article/div/div/div[3]/div[2]/div/a"
        content_link = None
        try:
            content_link = WebDriverWait(driver, WAITING_TIME_SEC_LINK).until(
                EC.presence_of_element_located((By.XPATH, link_xpath))
            ).text

            # 링크만 추출하는 패턴
            pattern = r"(?=[a-zA-Z0-9/-]*\.[a-zA-Z0-9/-])[a-zA-Z0-9./-]+"

            # 패턴에 맞는 부분 추출
            match = re.search(pattern, content_link)
            content_link = match.group(0)
            content_link = "https://" + content_link
        except TimeoutException:
            logger.error("Timeout loading link")
            logger.warning("Timeout loading link.")
        except Exception as e:
            print(f"General exception occurred: {e}")

        # calculate running time until now and log it as seconds
        running_time = datetime.datetime.now() - datetime.datetime.strptime(init_time,
                                                                            "%Y-%m-%d %H:%M:%S")
        logger.info(
            f"⏱️  Link has been scraped. Running time: {running_time.total_seconds()} sec")

        # log total running time
        logger.info(
            f"🏁 Total running time: {running_time.total_seconds()} sec, return now.")

        # return scraped data
        return jsonify({
            "text": text,
            "image": image,
            "username": username,
            "user_nickname": user_nickname,
            "user_profile_img": user_profile_img,
            "meta_tag": meta_tag,
            "link": content_link
        })

    except Exception as e:
        print("❌ Error occured ")
        print("My log app print exception:", e)
        # 에러 핸들러가 자동으로 호출되므로, 별도의 처리 없이도 됩니다.
        # Return the error details in the response
        return jsonify({
            "message": str(e)
        }), 500

    # finally:
    #     if driver:
    #         logger.info("Closing the WebDriver...")
    #         driver.quit()


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=18081)
