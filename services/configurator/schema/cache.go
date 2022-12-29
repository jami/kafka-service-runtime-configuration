package schema

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type JSONSchemaListEntry struct {
	Schema        interface{} `json:"schema"`
	Values        interface{} `json:"values"`
	ApplicationID string      `json:"id"`
}

type JSONSchemaList struct {
	Schemas []JSONSchemaListEntry `json:"schemas"`
}

type CacheEntry struct {
	schema    *gojsonschema.Schema
	plainJSON string
	value     interface{}
}

type Cache struct {
	Store map[string]CacheEntry
}

func CreateCache() *Cache {
	return &Cache{
		Store: map[string]CacheEntry{},
	}
}

func (c *Cache) Set(key string, data []byte) {
	jsonString := string(data)
	loader := gojsonschema.NewStringLoader(jsonString)
	schema, err := gojsonschema.NewSchema(loader)

	if err != nil {
		fmt.Printf("Error creating schema from loader %v\n", err)
	}

	c.Store[key] = CacheEntry{
		schema:    schema,
		plainJSON: jsonString,
	}
}

func (c *Cache) SetDefaultValues(key string, data []byte) {
	cacheEntry, ok := c.Store[key]

	if ok != true {
		return
	}

	var value interface{}
	var err error

	if err = json.Unmarshal(data, &value); err != nil {
		fmt.Printf("Error unmarshaling json %v\n", err)
		return
	}

	valueLoader := gojsonschema.NewGoLoader(value)
	validationResult, err := cacheEntry.schema.Validate(valueLoader)

	if err != nil {
		fmt.Printf("Error unmarshaling json %v\n", err)
		return
	}

	if validationResult.Valid() {
		fmt.Printf("The document is valid\n")
		cacheEntry.value = value
		c.Store[key] = cacheEntry
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, err := range validationResult.Errors() {
			// Err implements the ResultError interface
			fmt.Printf("- %s\n", err)
		}
	}

}

func (c *Cache) GetJSONSchemaList() JSONSchemaList {
	res := JSONSchemaList{
		Schemas: []JSONSchemaListEntry{},
	}

	for k, v := range c.Store {
		var target interface{}

		json.Unmarshal([]byte(v.plainJSON), &target)

		res.Schemas = append(
			res.Schemas,
			JSONSchemaListEntry{
				Schema:        target,
				ApplicationID: k,
				Values:        v.value,
			},
		)
	}

	return res
}
