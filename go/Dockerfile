# Stage 1: Go 애플리케이션 빌드
FROM golang:1.22-bullseye AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o server .

# Stage 2: 런타임 환경
FROM debian:bullseye-slim
WORKDIR /app

# 필수 패키지 설치
RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    gnupg \
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
    wget

# Chromium 설치 (ARM64용)
RUN apt-get update && apt-get install -y chromium

# Symlink 생성 (필요시)
RUN ln -s /usr/bin/chromium /usr/bin/chrome

# ChromeDriver 설치
RUN apt-get update && apt-get install -y chromium-driver


# Go 바이너리 복사
COPY --from=builder /app/server /usr/local/bin/server

EXPOSE 8080

CMD ["/usr/local/bin/server"]
