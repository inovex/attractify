FROM clickhouse/clickhouse-server:latest
RUN apt-get update \
  && apt-get install -y unixodbc odbcinst odbc-postgresql \
  && apt-get clean autoclean \
  && apt-get autoremove --yes \
  && rm -rf /var/lib/{apt,dpkg,cache,log}

# Add custom config files
ADD /docker/clickhouse/config/config.xml /etc/clickhouse-server/config.d/config.xml
ADD /docker/clickhouse/config/users.xml /etc/clickhouse-server/config.d/users.xml

# Add DB dependencies
ADD /docker/clickhouse/config/odbc.ini /etc/
ADD /docker/clickhouse/config/odbcinst_amd64.ini /etc/odbcinst.ini

ADD /server/schema/clickhouse.sql /docker-entrypoint-initdb.d/
