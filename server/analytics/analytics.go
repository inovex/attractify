package analytics

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

type Analytics struct {
	*sqlx.DB
	ClusterArgs *ClusterArgs
}

type ClusterArgs struct {
	ClusterName string
	LocalSuffix string
}

func OpenDB(dsn string) (*Analytics, error) {
	if strings.Contains(dsn, "tls_config=") {
		caCert, err := ioutil.ReadFile("../certs/clickhouse/ca/ca.crt")
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, errors.New("could not append CA cert to pool")
		}

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		clickhouse.RegisterTLSConfig("tls", tlsConfig)
	}

	db, err := sqlx.Connect("clickhouse", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			return nil, fmt.Errorf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	var clusterArgs *ClusterArgs
	parsedDSN, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	if parsedDSN.Query().Has("cluster") && parsedDSN.Query().Has("local_suffix") {
		clusterArgs.ClusterName = parsedDSN.Query().Get("cluster")
		clusterArgs.LocalSuffix = parsedDSN.Query().Get("local_suffix")
	}

	return &Analytics{db, clusterArgs}, nil
}

func (a Analytics) tableName(name string) string {
	if a.ClusterArgs == nil {
		return name
	}
	return fmt.Sprintf("%s_%s ON CLUSTER %s", name, a.ClusterArgs.LocalSuffix, a.ClusterArgs.ClusterName)
}
