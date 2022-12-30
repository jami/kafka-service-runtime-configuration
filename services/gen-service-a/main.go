package main

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	rtc "github.com/jami/kafka-service-runtime-configuration/component/runtime-configuration"
	"github.com/jami/kafka-service-runtime-configuration/services/gen-service-a/model"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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

	// build schema for later comparison
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
			// set the features

			switch ll := apc.LogLevel; ll {
			case "trace":
				logrus.SetLevel(log.TraceLevel)
			case "debug":
				logrus.SetLevel(log.DebugLevel)
			case "info":
				logrus.SetLevel(log.InfoLevel)
			case "warn":
				logrus.SetLevel(log.WarnLevel)
			case "error":
				logrus.SetLevel(log.ErrorLevel)
			case "fatal":
				logrus.SetLevel(log.FatalLevel)
			}
		},
		true,
	)

	defer func() {
		runtimeConfiguration.CloseListener()
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Trace("called /handler")
		log.Debug("request URL", r.URL)

		if r.URL.Path != "/" {
			log.Warn("more then just home")
		}

		log.Info("just a info")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":"+config.ListenerPort, nil))
}
