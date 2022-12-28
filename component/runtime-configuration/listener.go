package runtimeconfiguration

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const runtimeConfigurationTopic = "hotconfig.data.json.v1"
const runtimeConfigurationSchemaTopic = "hotconfig.schema.json.v1"

type RuntimeConfigurationListener struct {
	brokerURI      string
	consumer       *kafka.Consumer
	configTopic    string
	partitionCnt   int
	eofCnt         int
	configurations ConfigurationList
	handlerFunc    func(string, map[string]interface{})
}

type RuntimeConfigurationSchemaListener struct {
	brokerURI    string
	consumer     *kafka.Consumer
	configTopic  string
	partitionCnt int
	eofCnt       int
	schemas      SchemaList
}

func CreateListener(brokerURI string, handler func(string, map[string]interface{})) *RuntimeConfigurationListener {
	listener := RuntimeConfigurationListener{
		brokerURI:   brokerURI,
		configTopic: runtimeConfigurationTopic,
		handlerFunc: handler,
	}

	return &listener
}

func (l *RuntimeConfigurationListener) Listen() {
	go func() {
		var err error

		l.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":               l.brokerURI,
			"group.id":                        "configurator.config.group",
			"auto.offset.reset":               "latest",
			"enable.auto.commit":              false,
			"go.application.rebalance.enable": true,
		})

		if err != nil {
			panic(err)
		}

		metadata, err := l.consumer.GetMetadata(&l.configTopic, false, 1000)

		if err != nil {
			panic(err)
		}

		partitions := []kafka.TopicPartition{}

		for index, _ := range metadata.Topics[l.configTopic].Partitions {
			partitions = append(
				partitions,
				kafka.TopicPartition{
					Topic:     &l.configTopic,
					Partition: int32(index),
					Offset:    0,
				},
			)
		}

		l.consumer.Assign(partitions)

		for {
			ev := l.consumer.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				var m map[string]interface{}
				json.Unmarshal(e.Value, &m)
				l.handlerFunc(string(e.Key), m)
			case kafka.PartitionEOF:
				fmt.Fprintf(os.Stderr, "%% Reached %v\n", e)
				l.eofCnt++
				/*if exitEOF && l.eofCnt >= l.partitionCnt {
					//run = false
				}*/
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover.
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			case kafka.OffsetsCommitted:
				//if verbosity >= 2 {
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				//}
			case nil:
				// Ignore, Poll() timed out.
			default:
				fmt.Fprintf(os.Stderr, "%% Unhandled event %T ignored: %v\n", e, e)
			}
		}
	}()
}

func (l *RuntimeConfigurationListener) Close() {
	if l.consumer == nil {
		return
	}

	l.consumer.Close()
}

func CreateSchemaListener(brokerURI string) *RuntimeConfigurationSchemaListener {
	listener := RuntimeConfigurationSchemaListener{
		brokerURI:   brokerURI,
		configTopic: RuntimeConfigurationSchemaTopic,
		schemas:     SchemaList{},
	}

	return &listener
}

func (l *RuntimeConfigurationSchemaListener) GetSchemaList() SchemaList {
	return l.schemas
}

func (l *RuntimeConfigurationSchemaListener) Listen() {
	go func() {
		var err error

		l.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":               l.brokerURI,
			"group.id":                        "configurator.schema.group",
			"auto.offset.reset":               "latest",
			"enable.auto.commit":              false,
			"go.application.rebalance.enable": true,
		})

		if err != nil {
			panic(err)
		}

		metadata, err := l.consumer.GetMetadata(&l.configTopic, false, 1000)

		if err != nil {
			panic(err)
		}

		partitions := []kafka.TopicPartition{}

		for index, _ := range metadata.Topics[l.configTopic].Partitions {
			partitions = append(
				partitions,
				kafka.TopicPartition{
					Topic:     &l.configTopic,
					Partition: int32(index),
					Offset:    0,
				},
			)
		}

		l.consumer.Assign(partitions)

		for {
			ev := l.consumer.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				l.schemas.Append(string(e.Key), string(e.Value), e.Headers)
			case kafka.PartitionEOF:
				fmt.Fprintf(os.Stderr, "%% Reached %v\n", e)
				l.eofCnt++
				/*if exitEOF && l.eofCnt >= l.partitionCnt {
					//run = false
				}*/
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover.
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			case kafka.OffsetsCommitted:
				//if verbosity >= 2 {
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				//}
			case nil:
				// Ignore, Poll() timed out.
			default:
				fmt.Fprintf(os.Stderr, "%% Unhandled event %T ignored: %v\n", e, e)
			}
		}
	}()
}

func (l *RuntimeConfigurationSchemaListener) Close() {
	if l.consumer == nil {
		return
	}

	l.consumer.Close()
}
