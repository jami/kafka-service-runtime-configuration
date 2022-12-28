package runtimeconfiguration

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/xeipuuv/gojsonschema"
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

func DeserializeSchema(value []byte) (*Schema, error) {
	res := Schema{}

	if err := json.Unmarshal(value, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func Validate(jsonSchema *gojsonschema.Schema, jsonData []byte, v any) (bool, []error) {
	var err error

	if err = json.Unmarshal(jsonData, v); err != nil {
		return false, []error{err}

	}

	valueLoader := gojsonschema.NewGoLoader(v)
	validationResult, err := jsonSchema.Validate(valueLoader)

	if err != nil {
		return false, []error{err}
	}

	if validationResult.Valid() {
		fmt.Printf("The document is valid\n")
		return true, nil
	}

	errList := []error{}
	for _, schemaValidationError := range validationResult.Errors() {
		errList = append(errList, fmt.Errorf(schemaValidationError.String()))
	}
	return false, errList
}
