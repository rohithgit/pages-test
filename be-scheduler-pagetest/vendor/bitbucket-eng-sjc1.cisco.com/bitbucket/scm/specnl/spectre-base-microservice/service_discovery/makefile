.PHONY: default help runtest regsvc consul
.DEFAULT_GOAL := help

default: help

consul: ## run consul container
	@echo "run consul container"
	@echo "see: https://hub.docker.com/r/progrium/consul/ for image details"
	docker run -d -p 8400:8400 -p 8500:8500 -p 8600:53/udp --name consul -h node1 progrium/consul -server -bootstrap

regsvc: ## register service with consul
	@echo register google as a service
	curl -X PUT -d '{"Datacenter": "dc1", "Node": "google", "Address": "www.google.com", "Service": {"Service": "search", "Port": 80}}' http://127.0.0.1:8500/v1/catalog/register

runtest: ## run consol test
	@echo running consol test ...
	go run test/consultest.go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

