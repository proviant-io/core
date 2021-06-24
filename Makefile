ifndef TAG
TAG := "dev"
endif

.PHONY: docker/compile
docker/compile:
	CGO_ENABLED=1 go build -ldflags="-X 'main.SqliteLocation=/app/db/pantry.db'" -o app ./cmd

.PHONY: test/2e2/docker-build
test/2e2/docker-build:
	docker build -t proviant:e2e -f Dockerfile .

.PHONY: docker/build
docker/build:
	docker build -t proviant:$(TAG) -f Dockerfile .

.PHOMY: docker/publish
docker/publish:
	docker push brushknight/proviant:$(TAG)

.PHONY: docker/run
docker/run: docker/build
	docker run --rm -t --name "proviant" -p8080:80 proviant:latest

.PHONY: docker/run/test
docker/run/test: docker/build
	docker run -d --name "proviant-test" -p8081:80 proviant:latest

.PHONY: docker/stop/test
docker/stop/test:
	docker rm -f "proviant-test"

.PHONY: test/e2e
test/e2e: test/2e2/docker-build
	go test -v ./test/e2e/
