package kafka

import "github.com/confluentinc/confluent-kafka-go/kafka"

type ConsumerGroupConfig struct {
	Brokers          string `yaml:"brokers"`
	GroupID          string `yaml:"group_id"`
	Topic            string `yaml:"topic"`
	AutoOffsetReset  string `yaml:"auto_offset_reset"`
	EnableAutoCommit bool   `yaml:"enable_auto_commit"`
}

func NewKafkaConsumer(cfg ConsumerGroupConfig) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        cfg.Brokers,
		"group.id":                 cfg.GroupID,
		"auto.offset.reset":        cfg.AutoOffsetReset,
		"enable.auto.commit":       cfg.EnableAutoCommit,
		"enable.auto.offset.store": false, // Ручное сохранение offset'ов
		"session.timeout.ms":       6000,
		"max.poll.interval.ms":     300000,
	})

	if err != nil {
		return nil, err
	}

	err = consumer.SubscribeTopics([]string{cfg.Topic}, nil)
	if err != nil {
		consumer.Close()
		return nil, err
	}

	return consumer, nil
}
