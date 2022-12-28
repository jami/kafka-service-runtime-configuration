package runtimeconfiguration

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	AllApplications                      = ""
	waitForLatestConfigDurationInSeconds = 5.0
)

type OnMessageEvent func(key string, value []byte, headers []kafka.Header)

func (rc *RuntimeConfiguration) DataChangeListener(appId string, onMessage OnMessageEvent, waitForLatestAppId bool) {
	rc.DataConsumer = createConsumer(rc.BrokerURI, "runtimeconfiguration.data.group")
	latestTick := time.Now()

	var latestValue []byte
	var latestHeader []kafka.Header

	reachedLatestOffset := false

	listen(
		rc.DataConsumer,
		RuntimeConfigurationDataTopic,
		func(key string, value []byte, headers []kafka.Header) {
			if waitForLatestAppId && key == appId {
				latestValue = value
				latestHeader = headers
			}

			if key == AllApplications || key == appId {
				if waitForLatestAppId {
					if reachedLatestOffset {
						onMessage(key, value, headers)
					}
				} else {
					onMessage(key, value, headers)
				}
			}
		},
		rc.doneChannel,
	)

	if waitForLatestAppId == true {
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			stop := make(chan bool, 1)

			for {
				select {
				case <-ticker.C:
					currentTime := time.Now()
					if currentTime.Sub(latestTick).Seconds() > waitForLatestConfigDurationInSeconds {
						stop <- true
					}
				case <-stop:
					if latestValue != nil {
						onMessage(appId, latestValue, latestHeader)
					}
					reachedLatestOffset = true
					return
				}
			}
		}()
	}
}

func (rc *RuntimeConfiguration) SchemaChangeListener(onMessage OnMessageEvent) {
	rc.SchemaConsumer = createConsumer(rc.BrokerURI, "runtimeconfiguration.schema.group")
	listen(rc.SchemaConsumer, RuntimeConfigurationSchemaTopic, onMessage, rc.doneChannel)
}

func (rc *RuntimeConfiguration) CloseListener() {
	rc.doneChannel <- struct{}{}
}

func createConsumer(brokerURI string, groupId string) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               brokerURI,
		"group.id":                        groupId,
		"auto.offset.reset":               "latest",
		"enable.auto.commit":              false,
		"go.application.rebalance.enable": true,
	})

	if err != nil {
		panic(err)
	}

	return consumer
}

func listen(consumer *kafka.Consumer, topic string, onMessage OnMessageEvent, done chan struct{}) {
	go func() {
		var err error

		fmt.Println("Consumer started for topic: " + topic)

		if consumer == nil {
			fmt.Println("Consumer is nil and could not be started")
			return
		}

		defer func() {
			consumer.Close()
		}()

		// reset everything to 0

		metadata, err := consumer.GetMetadata(&topic, false, 1000)

		fmt.Printf("topic meta: %#v\n", metadata)

		if err != nil {
			panic(err)
		}

		partitions := []kafka.TopicPartition{}

		for index, _ := range metadata.Topics[topic].Partitions {
			partitions = append(
				partitions,
				kafka.TopicPartition{
					Topic:     &topic,
					Partition: int32(index),
					Offset:    0,
				},
			)
		}

		consumer.Assign(partitions)

		for {
			ev := consumer.Poll(100)

			switch e := ev.(type) {
			case *kafka.Message:
				onMessage(string(e.Key), e.Value, e.Headers)
			case kafka.PartitionEOF:
				fmt.Fprintf(os.Stderr, "%% Reached %v\n", e)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover.
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			case kafka.OffsetsCommitted:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
			case nil:
				// Ignore, Poll() timed out.
			default:
				fmt.Fprintf(os.Stderr, "%% Unhandled event %T ignored: %v\n", e, e)
			}
			select {
			case <-done:
				fmt.Printf("Closing listener for topic %s\n", topic)
				return
			case <-time.After(1 * time.Millisecond):
			}
		}
	}()
}
