FROM python:3.9-slim


# 작업 디렉토리 설정
WORKDIR /app


# Firefox 및 GeckoDriver 설치
RUN apt-get update && apt-get install -y --no-install-recommends \
    firefox-esr \
    wget \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# GeckoDriver 다운로드 및 설치
RUN wget https://github.com/mozilla/geckodriver/releases/download/v0.35.0/geckodriver-v0.35.0-linux-aarch64.tar.gz\
    && tar -xzf geckodriver-v0.35.0-linux-aarch64.tar.gz \
    && mv geckodriver /usr/local/bin/ \
    && rm geckodriver-v0.35.0-linux-aarch64.tar.gz



# requirements.txt 파일을 컨테이너에 복사하고, 필요한 패키지 설치
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt


# 애플리케이션 코드 복사
COPY app.py .


# 포트 설정
EXPOSE 18081

# Gunicorn을 사용해 Flask 애플리케이션 실행
CMD ["gunicorn", "-w", "4", "-b", "0.0.0.0:18081", "app:app"]
