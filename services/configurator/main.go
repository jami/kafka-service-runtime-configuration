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
			//fmt.Println("Received config change for key: " + key)
			schemaCache.SetDefaultValues(key, value)
		},
		true,
	)

	runtimeConfiguration.SchemaChangeListener(func(key string, value []byte, headers []kafka.Header) {
		//fmt.Println("Received schema change for key: " + key)
		schemaCache.Set(key, value)
	})

	defer func() {
		runtimeConfiguration.CloseListener()
	}()

	http.HandleFunc("/api/schema/list", func(w http.ResponseWriter, r *http.Request) {
		schemaListData, _ := json.MarshalIndent(schemaCache.GetJSONSchemaList(), "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(schemaListData)
	})

	http.HandleFunc("/api/schema/update", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			ID     string      `json:"id"`
			Values interface{} `json:"values"`
		}

		w.Header().Set("Content-Type", "application/json")

		// decode input or return error
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(400)
			data, _ := json.MarshalIndent([]error{fmt.Errorf("Decode error! please check your JSON formating. %s", err)}, "", "    ")
			w.Write(data)
			return
		}

		if entry, ok := schemaCache.Store[input.ID]; ok {
			if valid, errorList := entry.Validate(input.Values); !valid {
				w.WriteHeader(400)
				//fmt.Printf("Error %#v\n", errorList)
				data, _ := json.MarshalIndent(errorList, "", "    ")
				w.Write(data)
				return
			}
		} else {
			w.WriteHeader(400)
			data, _ := json.MarshalIndent([]string{fmt.Sprintf("Could not find cache entry with app id: %s", input.ID)}, "", "    ")
			w.Write(data)
			return
		}

		fmt.Printf("update %#v\n", input)
		data, _ := json.Marshal(input.Values)
		runtimeConfiguration.Update(input.ID, string(data))

		w.WriteHeader(200)
		w.Write([]byte("[]"))
	})

	staticSPA := http.FileServer(http.Dir(config.StaticPath))
	http.Handle("/", staticSPA)

	err := http.ListenAndServe(":"+config.ListenerPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
