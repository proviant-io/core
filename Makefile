ifndef TAG
TAG := dev
endif

ifndef UI_VERSION
UI_VERSION := 0.0.19
endif

# Compile
.PHONY: compile
compile:
	CGO_ENABLED=1 go build -o app ./cmd

# Testing
.PHONY: test/2e2/docker-build
test/2e2/docker-build:
	docker build --target app \
		--build-arg CONFIG_VERSION_ARG="./config/api-sqlite.yml" \
		--build-arg UI_VERSION_ARG=$(UI_VERSION) \
		-t brushknight/proviant:e2e \
		-f Dockerfile .

.PHONY: test/e2e
test/e2e: test/2e2/docker-build
	docker rm -f proviant-e2e
	go test -v ./test/e2e/

.PHONY: test/unit
test/unit:
	go test -v ./internal/...

# Download assets
.PHONY: download/ui
download/ui:
	curl -L https://github.com/brushknight/proviant-ui/releases/download/$(UI_VERSION)/ui-release-$(UI_VERSION).tar.gz -o /tmp/ui-release.tar.gz
	tar -xvf /tmp/ui-release.tar.gz -C ./public/

# Docker
.PHONY: docker/build
docker/build:
	docker build --target app \
		--build-arg CONFIG_VERSION_ARG="./config/web-sqlite.yml" \
		--build-arg UI_VERSION_ARG=$(UI_VERSION) \
		-t brushknight/proviant:$(TAG) \
		-t brushknight/proviant:latest \
		-f Dockerfile .

.PHOMY: docker/publish
docker/publish:
	docker push brushknight/proviant:$(TAG)
	docker push brushknight/proviant:latest

.PHONY: docker/run
docker/run: docker/build docker/prepare-folders
	docker rm -f proviant
	docker run --rm -t \
		--name "proviant" \
		-v $(PWD)/runtime/sqlite:/app/db/ \
		-v $(PWD)/runtime/user_content:/app/user_content/ \
		-v $(PWD)/config/web-sqlite.yml:/app/default-config.yml \
		-p8080:80 \
		brushknight/proviant:$(TAG)

.PHONY: docker/pull-latest
docker/pull-latest:
	docker pull brushknight/proviant:latest

.PHONY: docker/run-latest
docker/run-latest: docker/pull-latest docker/prepare-folders
	mkdir -p $(PWD)/db
	docker run --rm -t \
		--name "proviant" \
		-v $(PWD)/runtime/sqlite:/app/db/ \
		-v $(PWD)/runtime/user_content:/app/user_content/ \
		-v $(PWD)/config/simple.yml:/app/config.yml \
		-p8080:80 \
		brushknight/proviant:latest

# docker compose
.PHONY: docker/compose
docker/compose: docker/build docker/prepare-folders
	docker-compose up -d --force-recreate

.PHONY: docker/prepare-folders
docker/prepare-folders:
	mkdir -p $(PWD)/runtime/mysql
	mkdir -p $(PWD)/runtime/user_content

