
build/local/configurator:
	cd services/configurator/static/configurator-spa && npm i && npm run build
	go build -o bin/configurator  services/configurator/main.go 

build/local/genservicea:
	go build -o bin/genservicea  services/gen-service-a/main.go 

build/local: build/local/configurator build/local/genservicea