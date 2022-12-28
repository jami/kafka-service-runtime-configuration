package runtimeconfiguration

import (
	"fmt"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type SchemaList map[string]Schema

type Schema struct {
	ApplicationID string
	Version       int
	Data          string
}

func (s *SchemaList) Append(key string, value string, headers []kafka.Header) {
	schemaVersion := 1

	for _, h := range headers {
		if h.Key == "version" {
			if newVersion, err := strconv.ParseUint(string(h.Value), 10, 32); err == nil {
				schemaVersion = int(newVersion)
			}
		}
	}

	(*s)[key] = Schema{
		ApplicationID: key,
		Version:       schemaVersion,
		Data:          value,
	}

	fmt.Printf("schema store %v\n", *s)
}
