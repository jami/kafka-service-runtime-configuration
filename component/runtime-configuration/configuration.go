package runtimeconfiguration

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const RuntimeConfigurationSchemaTopic string = "hotconfig.schema.json.v1"
const RuntimeConfigurationDataTopic string = "hotconfig.data.json.v1"

type RuntimeConfiguration struct {
	BrokerURI               string
	DataConsumer            *kafka.Consumer
	SchemaConsumer          *kafka.Consumer
	doneChannel             chan struct{}
	waitForLatestDataChange bool
}

func CreateRuntimeConfiguration(brokerURI string) *RuntimeConfiguration {
	rtc := &RuntimeConfiguration{
		BrokerURI:               brokerURI,
		DataConsumer:            nil,
		SchemaConsumer:          nil,
		waitForLatestDataChange: false,
	}

	return rtc
}
