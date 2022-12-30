package runtimeconfiguration

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	AllApplications                      = ""
	waitForLatestConfigDurationInSeconds = 2.0
)

type OnMessageEvent func(key string, value []byte, headers []kafka.Header)

type preloaderEntry struct {
	value   []byte
	headers []kafka.Header
}

type preloader struct {
	appIDFilter string
	isReady     bool
	isActive    bool
	latestTick  time.Time
	store       map[string]preloaderEntry
	onMessage   OnMessageEvent
}

func (pl *preloader) update(key string, value []byte, headers []kafka.Header) {
	pl.store[key] = preloaderEntry{
		value:   value,
		headers: headers,
	}
	pl.latestTick = time.Now()
}

func (pl *preloader) isDone() bool {
	if pl.isActive == false {
		return true
	}

	return pl.isReady
}

func createPreloader(isActive bool, appId string, onMessage OnMessageEvent) *preloader {
	pl := &preloader{
		appIDFilter: appId,
		isReady:     false,
		isActive:    isActive,
		latestTick:  time.Now(),
		store:       map[string]preloaderEntry{},
	}

	if isActive == true {
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			stop := make(chan bool, 1)

			for {
				select {
				case <-ticker.C:
					//fmt.Println("tick")
					currentTime := time.Now()
					if currentTime.Sub(pl.latestTick).Seconds() > waitForLatestConfigDurationInSeconds {
						stop <- true
					}
				case <-stop:
					if pl.appIDFilter == AllApplications {
						for appId, entry := range pl.store {
							onMessage(appId, entry.value, entry.headers)
						}
					} else {
						if entry, ok := pl.store[pl.appIDFilter]; ok {
							onMessage(pl.appIDFilter, entry.value, entry.headers)
						}
					}
					pl.isReady = true
					return
				}
			}
		}()
	}

	return pl
}

func (rc *RuntimeConfiguration) DataChangeListener(appId string, onMessage OnMessageEvent, waitForLatestAppId bool) {
	rc.DataConsumer = createConsumer(rc.BrokerURI, "runtimeconfiguration.data.group")

	pl := createPreloader(waitForLatestAppId, appId, onMessage)

	listen(
		rc.DataConsumer,
		RuntimeConfigurationDataTopic,
		func(key string, value []byte, headers []kafka.Header) {
			//fmt.Printf("Receive message from %s key: %s wait: %v appId: %s\n", RuntimeConfigurationDataTopic, key, waitForLatestAppId, appId)
			pl.update(key, value, headers)

			if pl.isDone() {
				onMessage(key, value, headers)
			}
		},
		rc.doneChannel,
	)
}

func (rc *RuntimeConfiguration) SchemaChangeListener(onMessage OnMessageEvent) {
	rc.SchemaConsumer = createConsumer(rc.BrokerURI, "runtimeconfiguration.schema.group")
	listen(
		rc.SchemaConsumer,
		RuntimeConfigurationSchemaTopic,
		func(key string, value []byte, headers []kafka.Header) {
			//fmt.Printf("Receive message from %s key: %s\n", RuntimeConfigurationSchemaTopic, key)
			onMessage(key, value, headers)
		},
		rc.doneChannel,
	)
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
