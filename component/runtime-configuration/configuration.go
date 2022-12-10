package runtimeconfiguration

import (
	"fmt"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var RuntimeConfigurationSchemaTopic string = "hotconfig.schema.json.v1"

type ConfigurationList []Configuration

type Configuration struct {
	ApplicationID string
}

type Handler struct {
	BrokerURI string
}

func CreateHandler(brokerURI string) *Handler {
	handler := &Handler{
		BrokerURI: brokerURI,
	}

	return handler
}

func (h *Handler) PropagateSchema(appID string, schema string, schemaVersion int) error {
	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": h.BrokerURI,
		},
	)

	fmt.Printf("Error %v", err)

	if err != nil {
		return err
	}

	return producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &RuntimeConfigurationSchemaTopic,
				Partition: kafka.PartitionAny,
			},
			Value: []byte(schema),
			Key:   []byte(appID),
			Headers: []kafka.Header{
				{Key: "version", Value: []byte(strconv.Itoa(schemaVersion))},
			},
		},
		nil,
	)
}
