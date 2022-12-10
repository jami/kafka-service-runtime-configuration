package runtimeconfiguration

type SchemaProducer struct{}

type RuntimeConfigurationProducer struct {
}

func CreateRuntimeConfigurationProducer() *RuntimeConfigurationProducer {
	return &RuntimeConfigurationProducer{}
}
