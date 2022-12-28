package runtimeconfiguration

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (rc *RuntimeConfiguration) Update(appID string, value string) error {
	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": rc.BrokerURI,
		},
	)

	if err != nil {
		return err
	}

	dataTopic := RuntimeConfigurationDataTopic

	return producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &dataTopic,
				Partition: kafka.PartitionAny,
			},
			Value:   []byte(value),
			Key:     []byte(appID),
			Headers: []kafka.Header{},
		},
		nil,
	)
}
