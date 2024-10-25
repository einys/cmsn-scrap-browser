curl -X POST "http://localhost:18081/scrape-twitter" \
-H "Content-Type: application/json" \
-d '{"url": "https://x.com/Rook_commission/status/1831722898014257457"}'

curl -X GET "http://localhost:18081/meta?url=https://cafe.naver.com/usemapcomputer/438708?tc=shared_link"