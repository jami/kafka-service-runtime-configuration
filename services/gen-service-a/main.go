package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"

	rtc "github.com/jami/kafka-service-runtime-configuration/component/runtime-configuration"
	"github.com/jami/kafka-service-runtime-configuration/services/gen-service-a/model"
)

const applicationId = "genservicea"

func handleRuntimeConfigurationChange() {

}

func getFromEnv(k string, def string) string {
	e := os.Getenv(strings.ToUpper(applicationId + "_" + k))
	if e == "" {
		e = def
	}

	return e
}

func main() {
	config := model.AppConfiguration{
		BrokerURI:    getFromEnv("BrokerURI", "127.0.0.1:9092"),
		ListenerPort: getFromEnv("ListenerPort", "9091"),
	}

	rtConfigHandler := rtc.CreateHandler(config.BrokerURI)
	err := rtConfigHandler.PropagateSchema(
		applicationId,
		model.RuntimeConfigurationJSONSchema,
		model.RuntimeConfigurationJSONSchemaVersion,
	)

	println(err)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	log.Fatal(http.ListenAndServe(":"+config.ListenerPort, nil))
}
