.PHONY: frontend server postgres clickhouse-amd64 clickhouse-arm64 clickhouse push
release: frontend server postgres clickhouse push

server:
	docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/inovex/attractify/attractify-server -f server/Dockerfile . --push

postgres:
	docker build -t ghcr.io/inovex/attractify/attractify-postgres -f docker/postgres/Dockerfile .

clickhouse-amd64:
	docker build -t ghcr.io/inovex/attractify/attractify-clickhouse:amd64 -f docker/clickhouse/Dockerfile.amd64 .

clickhouse-arm64:
	docker build -t ghcr.io/inovex/attractify/attractify-clickhouse:arm64 -f docker/clickhouse/Dockerfile.arm64 .

clickhouse: clickhouse-arm64 clickhouse-amd64

frontend:
	cd frontend; ./build.sh

push:
	docker push ghcr.io/inovex/attractify/attractify-clickhouse:amd64
	docker push ghcr.io/inovex/attractify/attractify-clickhouse:arm64

	docker manifest rm ghcr.io/inovex/attractify/attractify-clickhouse
	docker manifest create ghcr.io/inovex/attractify/attractify-clickhouse ghcr.io/inovex/attractify/attractify-clickhouse:amd64 ghcr.io/inovex/attractify/attractify-clickhouse:arm64
	docker manifest push ghcr.io/inovex/attractify/attractify-clickhouse

	docker push ghcr.io/inovex/attractify/attractify-postgres
	docker push ghcr.io/inovex/attractify/attractify-server
