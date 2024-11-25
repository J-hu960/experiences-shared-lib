package logger

import (
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
		"level":      entry.Level,
		"message":    entry.Message,
		"fields":     entry.Data,
	}
	_, err := hook.client.Index().
		Index(hook.index).
		BodyJson(doc).
		Do(entry.Context)

	return err

}

func (hook *ElasticHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
