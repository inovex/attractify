FROM postgres
LABEL org.opencontainers.image.source https://github.com/inovex/attractify/postgres
ADD /docker/postgres/1_schema.sql /docker-entrypoint-initdb.d/
ADD /docker/postgres/2_fixtures.sql /docker-entrypoint-initdb.d/