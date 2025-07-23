# ───────────────────────────────────────────────────────────────
# Stage 1: Build Go binary
# ───────────────────────────────────────────────────────────────
FROM golang:1.23-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 정적 링크 빌드 (CGO_DISABLED)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -ldflags="-s -w" -o server .



# ───────────────────────────────────────────────────────────────
# Stage 2: Runtime image
# ───────────────────────────────────────────────────────────────
FROM debian:bullseye-slim

# 비루트 사용자 생성
RUN groupadd -r appuser && useradd -r -g appuser appuser

# 시스템 의존성 설치 (chromium, chromedriver 제거)
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    ca-certificates \
    fonts-liberation \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcups2 \
    libdbus-1-3 \
    libgdk-pixbuf2.0-0 \
    libnspr4 \
    libnss3 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    libgbm-dev \
    libasound2 \
    libxshmfence1 \
    libxss1 \
    libxtst6 \
    wget \
    unzip \
    libu2f-udev \
    libvulkan1 \
    && rm -rf /var/lib/apt/lists/*

# ARM64 환경에서 Chromium + Driver 설치
# 버전 120: Debian bullseye‑slim의 기본 APT 리포지터리에 들어 있는 chromium / chromium‑driver 패키지는 버전 120 계열로 고정돼 있습니다.
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    chromium \
    chromium-driver \
    && rm -rf /var/lib/apt/lists/* \
    # 경로 통일
    && ln -sf /usr/bin/chromium /usr/bin/chrome


# Chrome 커맨드 이름을 통일 (이미 ln -s로 처리)
# RUN ln -sf /usr/bin/chromium /usr/bin/chrome

# 환경 변수: headless Chrome 기본 경로
ENV CHROME_BIN=/usr/bin/chrome \
    # 일반적으로 CI 등에서 headless 모드로 실행할 때 필요
    DISPLAY=:99

WORKDIR /app
# 빌드된 Go 바이너리 복사 및 권한 설정
COPY --from=builder /app/server /usr/local/bin/server
RUN chown appuser:appuser /usr/local/bin/server

USER appuser

EXPOSE 18081

# 손쉬운 신호 처리(init)용 tini 사용을 원한다면 아래 추가
# RUN apt-get update && apt-get install -y --no-install-recommends tini && rm -rf /var/lib/apt/lists/*
# ENTRYPOINT ["/usr/bin/tini", "--"]

CMD ["/usr/local/bin/server"]
