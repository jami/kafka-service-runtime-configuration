
build/local/configurator:
	go build -o bin/configurator  services/configurator/main.go 

build/local: build/local/configurator