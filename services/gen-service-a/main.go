package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	rtc "github.com/jami/kafka-service-runtime-configuration/component/runtime-configuration"
	"github.com/jami/kafka-service-runtime-configuration/services/gen-service-a/model"
	"github.com/xeipuuv/gojsonschema"
)

const applicationId = "genservicea"

func getFromEnv(k string, def string) string {
	e := os.Getenv(strings.ToUpper(applicationId + "_" + k))
	if e == "" {
		e = def
	}

	return e
}

func onDataChangeHandler(data []byte) {

}

func main() {
	config := model.AppConfiguration{
		BrokerURI:    getFromEnv("BrokerURI", "127.0.0.1:9092"),
		ListenerPort: getFromEnv("ListenerPort", "9091"),
	}

	jsonSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(model.RuntimeConfigurationJSONSchema))

	if err != nil {
		fmt.Printf("Invalid json schema: %s\n", err)
	}

	runtimeConfiguration := rtc.CreateRuntimeConfiguration(config.BrokerURI)
	runtimeConfiguration.Register(
		applicationId,
		model.RuntimeConfigurationJSONSchema,
		model.RuntimeConfigurationJSONSchemaVersion,
	)

	runtimeConfiguration.DataChangeListener(
		applicationId,
		func(key string, value []byte, headers []kafka.Header) {
			fmt.Println("Received config change for key: " + key)

			var apc model.AppRuntimeConfiguration
			if valid, errorList := rtc.Validate(jsonSchema, value, &apc); !valid {
				for _, e := range errorList {
					fmt.Printf("Validator error: %s", e.Error())
				}
				return
			}

			fmt.Println("apc has a valid configuration")
		},
		true,
	)

	defer func() {
		runtimeConfiguration.CloseListener()
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	log.Fatal(http.ListenAndServe(":"+config.ListenerPort, nil))
}
