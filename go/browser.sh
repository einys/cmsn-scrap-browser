docker stack deploy -c docker-compose.browser.yml browser --with-registry-auth

# 서비스 이름 자동 추출
SERVICE_NAME=$(docker stack services browser --format '{{.Name}}' | head -n 1)

# 서비스 로그 실시간 모니터링
docker service logs -f $SERVICE_NAME --tail 10