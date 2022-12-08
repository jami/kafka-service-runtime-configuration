package main

import (
	"encoding/json"
	"fmt"

	rtc "github.com/jami/kafka-service-runtime-configuration/component/runtime-configuration/src"
)

func main() {
	rtcSchemaListener := rtc.CreateRuntimeConfigurationSchemaListener()
	rtcConfigListener := rtc.CreateRuntimeConfigurationListener()

	schemaList := rtcSchemaListener.GetLatestSchemaList()
	configList := rtcConfigListener.GetLatestConfigurationList()

	schemaListData, _ := json.MarshalIndent(schemaList, "schemalist", "    ")
	configListData, _ := json.MarshalIndent(configList, "configlist", "    ")

	fmt.Println("configurator")
	fmt.Println("schema")
	fmt.Println(schemaListData)
	fmt.Println("config")
	fmt.Println(configListData)
}
