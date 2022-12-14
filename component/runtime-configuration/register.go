package runtimeconfiguration

import (
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (rc *RuntimeConfiguration) Register(appID string, schema string, schemaVersion int) error {
	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": rc.BrokerURI,
		},
	)

	if err != nil {
		return err
	}

	schemaTopic := RuntimeConfigurationSchemaTopic

	return producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &schemaTopic,
				Partition: kafka.PartitionAny,
			},
			Value: []byte(schema),
			Key:   []byte(appID),
			Headers: []kafka.Header{
				{
					Key:   "version",
					Value: []byte(strconv.Itoa(schemaVersion)),
				},
			},
		},
		nil,
	)
}
