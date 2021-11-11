.PHONY: frontend server postgres clickhouse-amd64 clickhouse-arm64 clickhouse push
release: frontend server postgres clickhouse push

server:
	docker build -t ghcr.io/inovex/attractify/server -f server/Dockerfile .

postgres:
	cp server/schema/postgres.sql docker/postgres/1_schema.sql
	cp server/testdata/fixtures/postgres.sql docker/postgres/2_fixtures.sql
	docker build -t ghcr.io/inovex/attractify/postgres -f docker/postgres/Dockerfile .

clickhouse-amd64:
	cp server/schema/clickhouse.sql docker/clickhouse/schema.sql
	docker build -t ghcr.io/inovex/attractify/clickhouse:amd64 -f docker/clickhouse/Dockerfile.amd64 .

clickhouse-arm64:
	cp server/schema/clickhouse.sql docker/clickhouse/schema.sql
	docker build -t ghcr.io/inovex/attractify/clickhouse:arm64 -f docker/clickhouse/Dockerfile.arm64 .

clickhouse: clickhouse-arm64 clickhouse-amd64

frontend:
	cd frontend; ./build.sh

push:
	docker push ghcr.io/inovex/attractify/clickhouse:amd64
	docker push ghcr.io/inovex/attractify/clickhouse:arm64

	docker manifest rm ghcr.io/inovex/attractify/clickhouse
	docker manifest create ghcr.io/inovex/attractify/clickhouse ghcr.io/inovex/attractify/clickhouse:amd64 ghcr.io/inovex/attractify/clickhouse:arm64
	docker manifest push ghcr.io/inovex/attractify/clickhouse

	docker push ghcr.io/inovex/attractify/postgres
	docker push ghcr.io/inovex/attractify/server

login-test:
	docker login ghcr.io/inovex/attractify/

up:
	docker-compose up

down:
	docker-compose down