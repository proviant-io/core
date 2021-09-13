#!make
include envfile
export $(shell sed 's/=.*//' envfile)

ifndef TAG
TAG := dev
endif

ifndef UI_VERSION
UI_VERSION := 0.0.38
endif

ifndef VOLUMES_PREFIX
VOLUMES_PREFIX=$(PWD)
endif

# Compile
.PHONY: compile
compile:
	CGO_ENABLED=1 go build -ldflags="-X 'main.Version=$(TAG)'" -o app ./cmd

# Testing
.PHONY: test/2e2/docker-build
test/2e2/docker-build:
	docker build --target app \
		--build-arg CONFIG_VERSION_ARG="./examples/config/api-sqlite.yml" \
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
	curl -L https://github.com/proviant-io/ui/releases/download/$(UI_VERSION)/ui-release-$(UI_VERSION)-ce.tar.gz -o /tmp/ui-release.tar.gz
	tar -xvf /tmp/ui-release.tar.gz -C ./public/

# Docker
.PHONY: docker/build
docker/build:
	docker build --target app \
		--build-arg CONFIG_VERSION_ARG="./examples/config/web-sqlite.yml" \
		--build-arg UI_VERSION_ARG=$(UI_VERSION) \
		--build-arg TAG_ARG=$(TAG) \
		-t brushknight/proviant-core:$(TAG) \
		-t brushknight/proviant-core:latest \
		-f Dockerfile .

.PHOMY: docker/publish
docker/publish:
	docker push brushknight/proviant-core:$(TAG)
	docker push brushknight/proviant-core:latest

.PHONY: docker/run
docker/run: docker/build docker/prepare-folders
	docker rm -f proviant
	docker run --rm -t \
		--name "proviant" \
		-v $(VOLUMES_PREFIX)/sqlite:/app/db/ \
		-v $(VOLUMES_PREFIX)/user_content:/app/user_content/ \
		-v $(PWD)/examples/config/web-sqlite.yml:/app/default-config.yml \
		-p8100:80 \
		brushknight/proviant-core:$(TAG)

.PHONY: docker/pull-latest
docker/pull-latest:
	docker pull brushknight/proviant-core:latest

.PHONY: docker/run-latest
docker/run-latest: docker/pull-latest docker/prepare-folders
	docker run --rm -t \
		--name "proviant" \
		-v $(VOLUMES_PREFIX)/sqlite:/app/db/ \
		-v $(VOLUMES_PREFIX)/user_content:/app/user_content/ \
		-v $(PWD)/examples/config/web-sqlite.yml:/app/default-config.yml \
		-p8080:80 \
		brushknight/proviant-core:latest

# docker compose
.PHONY: docker/compose
docker/compose: docker/build docker/prepare-folders
	docker-compose up -d --force-recreate

.PHONY: docker/prepare-folders
docker/prepare-folders:
	mkdir -p $(VOLUMES_PREFIX)/mysql
	mkdir -p $(VOLUMES_PREFIX)/user_content

