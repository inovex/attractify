.PHONY: frontend server postgres clickhouse-amd64 clickhouse-arm64 clickhouse push
release: frontend server postgres clickhouse push

server:
	docker build -t registry.inovex.de:4567/attractify/platform/demo/server -f server/Dockerfile .

postgres:
	cp server/schema/postgres.sql docker/postgres/1_schema.sql
	cp server/testdata/fixtures/postgres.sql docker/postgres/2_fixtures.sql
	docker build -t registry.inovex.de:4567/attractify/platform/demo/postgres -f docker/postgres/Dockerfile .

clickhouse-amd64:
	cp server/schema/clickhouse.sql docker/clickhouse/schema.sql
	docker build -t registry.inovex.de:4567/attractify/platform/demo/clickhouse:amd64 -f docker/clickhouse/Dockerfile.amd64 .

clickhouse-arm64:
	cp server/schema/clickhouse.sql docker/clickhouse/schema.sql
	docker build -t registry.inovex.de:4567/attractify/platform/demo/clickhouse:arm64 -f docker/clickhouse/Dockerfile.arm64 .

clickhouse: clickhouse-arm64 clickhouse-amd64

frontend:
	cd frontend; ./build.sh

push:
	docker push registry.inovex.de:4567/attractify/platform/demo/clickhouse:amd64
	docker push registry.inovex.de:4567/attractify/platform/demo/clickhouse:arm64

	docker manifest rm registry.inovex.de:4567/attractify/platform/demo/clickhouse
	docker manifest create registry.inovex.de:4567/attractify/platform/demo/clickhouse registry.inovex.de:4567/attractify/platform/demo/clickhouse:amd64 registry.inovex.de:4567/attractify/platform/demo/clickhouse:arm64
	docker manifest push registry.inovex.de:4567/attractify/platform/demo/clickhouse

	docker push registry.inovex.de:4567/attractify/platform/demo/postgres
	docker push registry.inovex.de:4567/attractify/platform/demo/server

login-test:
	docker login registry.inovex.de:4567/attractify/

up:
	docker-compose up

down:
	docker-compose down