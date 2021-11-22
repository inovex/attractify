FROM postgres
LABEL org.opencontainers.image.source https://github.com/inovex/attractify/postgres
ADD /server/schema/postgres.sql /docker-entrypoint-initdb.d/1_schema.sql
ADD /server/testdata/fixtures/postgres.sql /docker-entrypoint-initdb.d/2_fixtures.sql