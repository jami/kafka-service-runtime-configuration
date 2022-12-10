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

type Configuration struct {
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

func (c *Configuration) LoadFromEnv() {
	c.ListenerPort = getFromEnv("ListenerPort", c.ListenerPort)
	c.BrokerURI = getFromEnv("BrokerURI", c.BrokerURI)
	c.StaticPath = getFromEnv("StaticPath", c.StaticPath)
}

func main() {
	appConfig := Configuration{
		ListenerPort: "7331",
		BrokerURI:    "127.0.0.1:9091",
		StaticPath:   "services/configurator/static/configurator-spa/dist",
	}

	appConfig.LoadFromEnv()
	fmt.Printf("config %#v\n", appConfig)

	rtcSchemaListener := rtc.CreateRuntimeConfigurationSchemaListener()
	rtcConfigListener := rtc.CreateRuntimeConfigurationListener()

	schemaList := rtcSchemaListener.GetLatestSchemaList()
	configList := rtcConfigListener.GetLatestConfigurationList()

	schemaListData, _ := json.MarshalIndent(schemaList, "schemalist", "    ")
	configListData, _ := json.MarshalIndent(configList, "configlist", "    ")

	fmt.Println("configurator")
	fmt.Println("schema")
	fmt.Printf("%s\n", schemaListData)
	fmt.Println("config")
	fmt.Printf("%s\n", configListData)

	staticSPA := http.FileServer(http.Dir(appConfig.StaticPath))
	http.Handle("/", staticSPA)

	err := http.ListenAndServe(":"+appConfig.ListenerPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
