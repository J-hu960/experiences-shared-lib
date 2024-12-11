package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type ElasticHook struct {
	client   *elastic.Client
	index    string
	logLevel logrus.Level
}

func (hook *ElasticHook) Fire(entry *logrus.Entry) error {
	doc := map[string]interface{}{
		"@timestamp": time.Now().UTC(),
		"level":      entry.Level.String(),
		"message":    entry.Message,
		"fields":     entry.Data,
	}

	_, err := hook.client.Index().
		Index(hook.index).
		BodyJson(doc).
		Do(context.Background())

	if err != nil {
		return fmt.Errorf("error indexing log entry: %w", err)
	}
	return nil
}

func (hook *ElasticHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
