package logger

import (
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

var (
	logger      *logrus.Logger
	initialized bool
)

// Utilizamos una sola intsancia, un Singleton.
func InitializeLogger(elasticSearchURL string) error {
	if initialized {
		return nil
	}

	client, err := elastic.NewClient(
		elastic.SetURL("https://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}),
	)
	if err != nil {
		return err
	}

	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.AddHook(&ElasticHook{
		client:   client,
		index:    "logs",
		logLevel: logrus.InfoLevel,
	})

	initialized = true
	return nil
}

// Aqui devolvemos la instancia de logger si ya tenemos uno (un Singleton).
func GetLogger() (*logrus.Logger, error) {
	if !initialized {
		return nil, errors.New("logger has not been initialized")
	}
	return logger, nil
}
