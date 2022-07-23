package analytics

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

type Analytics struct {
	*sqlx.DB
	IsCluster   bool
	ClusterArgs ClusterArgs
}

type ClusterArgs struct {
	Cluster     string
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

	var clusterArgs ClusterArgs
	isCluster := false

	if strings.Contains(dsn, "&cluster") && strings.Contains(dsn, "&local_suffix") {
		isCluster = true
		dsnArgs := strings.Split(dsn, "&")

		for _, v := range dsnArgs {
			splitted := strings.Split(v, "=")
			key := splitted[0]
			value := splitted[1]
			if key == "cluster" {
				clusterArgs.Cluster = value
			} else if key == "local_suffix" {
				clusterArgs.LocalSuffix = value
			}
		}
	}

	return &Analytics{db, isCluster, clusterArgs}, nil
}
