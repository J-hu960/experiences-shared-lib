package logger

import (
	"errors"

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

	client, err := elastic.NewClient(elastic.SetURL(elasticSearchURL))
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

// GetLogger devuelve la instancia del logger inicializado
func GetLogger() (*logrus.Logger, error) {
	if !initialized {
		return nil, errors.New("logger has not been initialized")
	}
	return logger, nil
}
