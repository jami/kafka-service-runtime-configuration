package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	rtc "github.com/jami/kafka-service-runtime-configuration/component/runtime-configuration"
	"github.com/jami/kafka-service-runtime-configuration/services/configurator/schema"
)

const applicationId = "configurator"

type AppConfiguration struct {
	BrokerURI    string
	StaticPath   string
	ListenerPort string
}

func getFromEnv(k string, def string) string {
	e := os.Getenv(strings.ToUpper(applicationId + "_" + k))
	if e == "" {
		e = def
	}

	return e
}

func main() {
	config := AppConfiguration{
		ListenerPort: getFromEnv("ListenerPort", "8090"),
		BrokerURI:    getFromEnv("BrokerURI", "127.0.0.1:9092"),
		StaticPath:   getFromEnv("StaticPath", "services/configurator/static/configurator-spa/dist"),
	}

	schemaCache := schema.CreateCache()

	fmt.Printf("config %#v\n", config)

	runtimeConfiguration := rtc.CreateRuntimeConfiguration(config.BrokerURI)

	runtimeConfiguration.DataChangeListener(
		rtc.AllApplications,
		func(key string, value []byte, headers []kafka.Header) {
			fmt.Println("Received config change for key: " + key)
			schemaCache.SetDefaultValues(key, value)
		},
		false,
	)

	runtimeConfiguration.SchemaChangeListener(func(key string, value []byte, headers []kafka.Header) {
		fmt.Println("Received schema change for key: " + key)
		schemaCache.Set(key, value)
	})

	/*

		rtcSchemaListener := rtc.CreateSchemaListener(config.BrokerURI)
		rtcSchemaListener.Listen()

		rtcListener := rtc.CreateListener(config.BrokerURI, configurationHandler)
		rtcListener.Listen()

		defer func() {
			rtcSchemaListener.Close()
			rtcListener.Close()
		}()
	*/
	//schemaList := rtcSchemaListener.GetLatestSchemaList()
	//configList := rtcListener.GetLatestConfigurationList()
	/*
		schemaListData, _ := json.MarshalIndent(schemaList, "schemalist", "    ")
		configListData, _ := json.MarshalIndent(configList, "configlist", "    ")

		fmt.Println("configurator")
		fmt.Println("schema")
		fmt.Printf("%s\n", schemaListData)
		fmt.Println("config")
		fmt.Printf("%s\n", configListData)
	*/

	http.HandleFunc("/api/schema/test", func(w http.ResponseWriter, r *http.Request) {
		schemaListData, _ := json.MarshalIndent(struct{}{}, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(schemaListData)

		runtimeConfiguration.Update(
			"genservicea",
			`
			{
				"logLevel": "info"
			}
			`,
		)
	})

	staticSPA := http.FileServer(http.Dir(config.StaticPath))
	http.Handle("/", staticSPA)

	err := http.ListenAndServe(":"+config.ListenerPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
