NAME   := 412530390400.dkr.ecr.ap-northeast-2.amazonaws.com/cmsn-scrap-browser
PUBLIC_NAME := public.ecr.aws/v0z0y2u1/cmsn-scrap-browser
TAG    := $$(git log -1 --pretty=%h)
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest
DIFF   := $$(git status --porcelain)
 
# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

## https://stackoverflow.com/questions/3878624/how-do-i-programmatically-determine-if-there-are-uncommitted-changes
## 참고용으로 남겨뒀고 현재 사용안함
diff:
	@echo $( if [ -n "$$(git status --porcelain)" ]; then \
		echo "there are changes"; \
		false; \
	else \
		echo "tree is clean"; \
		true; \
	fi )

## git에 diff가 있는지 확인하고 있으면 exit 1 반환
## || 기호는 앞에서부터 순차적으로 실행하되, 명령 실행에 성공하면 뒤에 오는 명령을 실행하지 않는다. | 기호는 둘을 이어서 실행하는 것, &&는 무조건 둘 다 실행
## echo -e \033[47;31m 어쩌구 하는 부분은 출력 메시지 색상 지정
checkdiff:
	@(git diff --exit-code && git diff --cached --exit-code) || ( echo -e "\033[47;31m Git tree is dirty! Commit first. \033[0m" && exit 1 )

echo:
	@echo "run!"

show-commit: ## Show a commit by hash. ex) make show-commit HASH=1q2w3e
	@git show $(HASH)

show-commit-list: ## Show pretty commit list.
	@git log --pretty=format:"%h - %an, %ar : %s"

login: ## ECR login.
	@rm -f ~/.docker/config.json && aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 412530390400.dkr.ecr.ap-northeast-2.amazonaws.com

login-public: ## ECR login.
	@aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/v0z0y2u1

build: ## Build the container and tag with git commit's hash.
	@docker build -t ${IMG} .
	@docker tag ${IMG} ${LATEST}
 
push: ## Push docker image to ecr repo.
	@docker push ${IMG}
	@docker push ${LATEST}

build-public: ## Build the container and tag with git commit's hash.
	@docker build -t ${PUBLIC_NAME}:${TAG} .
	@docker tag ${PUBLIC_NAME}:${TAG} ${PUBLIC_NAME}:latest

push-public: ## Push docker image to ecr repo.
	@docker push ${PUBLIC_NAME}:${TAG}
	@docker push ${PUBLIC_NAME}:latest

all: checkdiff login-public build-public push-public ## All publish process.

all-local: checkdiff build push ## All publish process.


run:
	go run ./cmd/server

run-dp:
	SCRAPER_ENGINE=chromedp go run ./cmd/server

browser: ## Deploy the browser stack and show logs.
	@sh browser.sh

rm-browser: ## Remove the browser stack.
	@docker stack rm browser