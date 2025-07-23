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
잘 안되면 app/page.html 보면 로딩된 페이지에 대한 자료 있으니까 보고 뭐가 로딩안됐는지 확인하기.

## /scrape-twitter
```
curl "http://localhost:18081/scrape-twitter?url=https://x.com/naeng2_/status/1903488320367403357"
```

결과
```
{"text":"상시 흑백 그림커미션을 개장했습니다~\nhttps://kre.pe/V5LG\n자세한 사항 크레페 링크를 확인 부탁드립니다.","images":["https://pbs.twimg.com/media/GmqMNf0bYAAUnTD?format=png\u0026name=small","https://pbs.twimg.com/media/GmqMh1maAAAdIXL?format=png\u0026name=360x360"],"username":"@naeng2_","user_nickname":"냉이","user_profile_img":"https://pbs.twimg.com/profile_images/1843649710072225792/PyeAorAY_normal.jpg","meta_tag":"냉이 on X: \"상시 흑백 그림커미션을 개장했습니다~\nhttps://t.co/Bcu5BZZLkH\n자세한 사항은 크레페 링크를 확인 부탁드립니다. https://t.co/iFdaKGuPnH\" / X","links":["https://kre.pe/V5LG"]}
```

## /meta
```
curl "http://localhost:8080/meta?url=https://www.naver.com"
```

```
{"description":"네이버 메인에서 다양한 정보와 유용한 컨텐츠를 만나 보세요","img":"https://s.pstatic.net/static/www/mobile/edit/2016/0705/mobile_212852414260.png","title":"네이버"}
```