package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	rtc "github.com/jami/kafka-service-runtime-configuration/component/runtime-configuration"
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

func configurationHandler(id string, data map[string]interface{}) {

}

func main() {
	config := AppConfiguration{
		ListenerPort: getFromEnv("ListenerPort", "8090"),
		BrokerURI:    getFromEnv("BrokerURI", "127.0.0.1:9092"),
		StaticPath:   getFromEnv("StaticPath", "services/configurator/static/configurator-spa/dist"),
	}

	fmt.Printf("config %#v\n", config)

	rtcSchemaListener := rtc.CreateSchemaListener(config.BrokerURI)
	rtcSchemaListener.Listen()

	rtcListener := rtc.CreateListener(config.BrokerURI, configurationHandler)
	rtcListener.Listen()

	defer func() {
		rtcSchemaListener.Close()
		rtcListener.Close()
	}()

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

	http.HandleFunc("/api/schema/list", func(w http.ResponseWriter, r *http.Request) {
		schemaListData, _ := json.MarshalIndent(rtcSchemaListener.GetSchemaList(), "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(schemaListData)
	})

	staticSPA := http.FileServer(http.Dir(config.StaticPath))
	http.Handle("/", staticSPA)

	err := http.ListenAndServe(":"+config.ListenerPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
