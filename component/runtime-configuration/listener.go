package runtimeconfiguration

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const runtimeConfigurationTopic = "hotconfig.data.json.v1"
const runtimeConfigurationSchemaTopic = "hotconfig.schema.json.v1"

type RuntimeConfigurationListener struct {
	brokerURI    string
	consumer     *kafka.Consumer
	configTopic  string
	partitionCnt int
	eofCnt       int
}

type RuntimeConfigurationSchemaListener struct {
	brokerURI    string
	consumer     *kafka.Consumer
	configTopic  string
	partitionCnt int
	eofCnt       int
}

func CreateListener(brokerURI string) *RuntimeConfigurationListener {
	listener := RuntimeConfigurationListener{
		brokerURI:   brokerURI,
		configTopic: runtimeConfigurationTopic,
	}

	return &listener
}

func (l *RuntimeConfigurationListener) Listen() {
	go func() {
		var err error

		l.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":  l.brokerURI,
			"group.id":           "configurator.group",
			"auto.offset.reset":  "earliest",
			"enable.auto.commit": "false",
		})

		if err != nil {
			panic(err)
		}

		if err = l.consumer.SubscribeTopics([]string{l.configTopic}, l.rebalanceCallback); err != nil {
			panic(err)
		}

		for {
			message, merr := l.consumer.ReadMessage(10 * time.Second)
			fmt.Println("Message")
			fmt.Printf("err %v %#v\n", merr, message)
		}
	}()
}

func (l *RuntimeConfigurationListener) Close() {
	if l.consumer == nil {
		return
	}

	l.consumer.Close()
}

func (l *RuntimeConfigurationListener) rebalanceCallback(c *kafka.Consumer, ev kafka.Event) error {
	var err error = nil
	switch e := ev.(type) {
	case kafka.AssignedPartitions:
		fmt.Fprintf(os.Stderr, "%% %v\n", e)
		err = l.consumer.Assign(e.Partitions)
		l.partitionCnt = len(e.Partitions)
		l.eofCnt = 0
	case kafka.RevokedPartitions:
		fmt.Fprintf(os.Stderr, "%% %v\n", e)
		err = l.consumer.Unassign()
		l.partitionCnt = 0
		l.eofCnt = 0
	}
	return err
}

func CreateSchemaListener(brokerURI string) *RuntimeConfigurationSchemaListener {
	listener := RuntimeConfigurationSchemaListener{
		brokerURI:   brokerURI,
		configTopic: RuntimeConfigurationSchemaTopic,
	}

	return &listener
}

func (l *RuntimeConfigurationSchemaListener) Listen() {
	go func() {
		var err error

		l.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":  l.brokerURI,
			"group.id":           "configurator.schema.group",
			"auto.offset.reset":  "earliest",
			"enable.auto.commit": "false",
		})

		if err != nil {
			panic(err)
		}

		if err = l.consumer.SubscribeTopics([]string{l.configTopic}, l.rebalanceCallback); err != nil {
			panic(err)
		}

		for {
			message, merr := l.consumer.ReadMessage(10 * time.Second)
			fmt.Println("Schema Message")
			fmt.Printf("err %v %#v\n", merr, message)
		}
	}()
}

func (l *RuntimeConfigurationSchemaListener) Close() {
	if l.consumer == nil {
		return
	}

	l.consumer.Close()
}

func (l *RuntimeConfigurationSchemaListener) rebalanceCallback(c *kafka.Consumer, ev kafka.Event) error {
	var err error = nil
	switch e := ev.(type) {
	case kafka.AssignedPartitions:
		fmt.Fprintf(os.Stderr, "%% %v\n", e)
		//err = l.consumer.Assign(e.Partitions)
		l.partitionCnt = len(e.Partitions)
		l.eofCnt = 0
	case kafka.RevokedPartitions:
		fmt.Fprintf(os.Stderr, "%% %v\n", e)
		//err = l.consumer.Unassign()
		l.partitionCnt = 0
		l.eofCnt = 0
	}
	return err
}
