.PHONY: up-db
up-db:
#	@mkdir ./mariadb
	@docker-compose -f docker-compose.dev.yml up -d

.PHONY: docker/compile
docker/compile:
	CGO_ENABLED=1 go build -ldflags="-X 'main.SqliteLocation=/app/db/pantry.db'" -o app ./cmd


.PHONY: docker/build
docker/build:
	docker build -t pantry:latest -f Dockerfile .


.PHONY: docker/run
docker/run: docker/build
	docker run --rm -t --name "pantry" -p8080:80 pantry:latest
