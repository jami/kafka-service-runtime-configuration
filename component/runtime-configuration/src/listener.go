package src

type RuntimeConfigurationListener struct {
}

type RuntimeConfigurationSchemaListener struct {
}

func CreateRuntimeConfigurationListener() *RuntimeConfigurationListener {
	return &RuntimeConfigurationListener{}
}

func CreateRuntimeConfigurationSchemaListener() *RuntimeConfigurationSchemaListener {
	return &RuntimeConfigurationSchemaListener{}
}

func (r RuntimeConfigurationSchemaListener) GetLatestSchemaList() SchemaList {
	retList := SchemaList{}

	return retList
}

func (r RuntimeConfigurationListener) GetLatestConfigurationList() ConfigurationList {
	retList := ConfigurationList{}

	return retList
}
