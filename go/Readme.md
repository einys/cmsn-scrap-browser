# headless chrome twitter scraper

```
"--headless=new",
"--disable-gpu",
"--no-sandbox",
"--window-size=1280,1024",
"--lang=ko-KR,ko",
"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
```
이 옵션을 써야 됨. 아니면 트위터가 차단함. 또는, headless 주석처리하고 렌더링을 하면 됨.