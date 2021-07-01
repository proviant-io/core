ifndef TAG
TAG := dev
endif

ifndef UI_VERSION
UI_VERSION := 0.0.8
endif

.PHONY: docker/compile
docker/compile:
	CGO_ENABLED=1 go build -ldflags="-X 'main.SqliteLocation=/app/db/proviant.db'" -o app ./cmd

.PHONY: test/2e2/docker-build
test/2e2/docker-build:
	docker build -t brushknight/proviant:e2e -f Dockerfile .

.PHONY: docker/build
docker/build:
	docker build --build-arg UI_VERSION_ARG=$(UI_VERSION) -t brushknight/proviant:$(TAG) -t brushknight/proviant:latest -f Dockerfile .

.PHOMY: docker/publish
docker/publish:
	docker push brushknight/proviant:$(TAG)
	docker push brushknight/proviant:latest

.PHONY: docker/run
docker/run: docker/build
	mkdir -p $(PWD)/db
	docker run --rm -t --name "proviant" -v $(PWD)/db:/app/db/ -p8080:80 brushknight/proviant:$(TAG)

.PHONY: test/e2e
test/e2e: test/2e2/docker-build
	go test -v ./test/e2e/

.PHONY: download/ui
download/ui:
	curl -L https://github.com/brushknight/proviant-ui/releases/download/$(UI_VERSION)/ui-release-$(UI_VERSION).tar.gz -o /tmp/ui-release.tar.gz
	tar -xvf /tmp/ui-release.tar.gz -C ./public/