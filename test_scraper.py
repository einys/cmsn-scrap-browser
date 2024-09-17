import unittest
from app import app  # Import your Flask app
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

from selenium.webdriver.chrome.options import Options


class TestMetaScraper(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        # Use ChromeOptions instead of DesiredCapabilities
        chrome_options = Options()
        chrome_options.add_argument('--headless')  # Run in headless mode
        chrome_options.add_argument('--no-sandbox')
        chrome_options.add_argument('--disable-dev-shm-usage')

        # Use the remote Selenium WebDriver with ChromeOptions
        cls.driver = webdriver.Remote(
            command_executor='http://localhost:4444/wd/hub',
            options=chrome_options
        )
        cls.driver.implicitly_wait(10)
        cls.app = app.test_client()
        cls.app.testing = True

    @classmethod
    def tearDownClass(cls):
        # Quit the WebDriver after all tests
        cls.driver.quit()

    def test_meta_scraping(self):
        # Replace with a URL you want to test
        test_url = "https://daisyui.com/"

        # Start the Flask test client
        response = self.app.get(f'/meta?url={test_url}')
        self.assertEqual(response.status_code, 200)

        # Verify the response contains the meta data
        json_data = response.get_json()
        self.assertIn('title', json_data)
        self.assertIn('image', json_data)
        self.assertIn('description', json_data)

        # Further checks for data correctness
        self.assertIsInstance(json_data['title'], str)
        self.assertIsInstance(json_data['image'], str)
        self.assertIsInstance(json_data['description'], str)
        print(f"Meta data: {json_data}")


if __name__ == '__main__':
    unittest.main()
