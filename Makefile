.PHONY: up-db
up-db:
#	@mkdir ./mariadb
	@docker-compose -f docker-compose.dev.yml up -d