package model

const RuntimeConfigurationJSONSchema = `
{
	"type": "object",
	"properties": {
		"logLevel": { 
			"enum": ["error", "info", "warn", "debug"]
		}
	}
}`

const RuntimeConfigurationJSONSchemaVersion int = 1

type AppConfiguration struct {
	BrokerURI    string
	ListenerPort string
}

type AppRuntimeConfiguration struct {
	LogLevel string `json:"logLevel"`
	FeatureA struct {
		Enabled bool `json:"enabled"`
	} `json:"featureA"`
	FeatureB struct {
		Enabled bool `json:"enabled"`
	} `json:"featureB"`
	FeatureC struct {
		Enabled bool `json:"enabled"`
	} `json:"featureC"`
}
