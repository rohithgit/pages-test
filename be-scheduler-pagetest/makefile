.PHONY: default mac linux container
default: mac

mac:
	@echo "building page test service (mac)" ;\
	go build

linux:
	@echo "building page test service (linux)" ;\
	env GOOS=linux GOARCH=amd64 go build

windows:
	@echo "building page test service (windows)" ;\
	env GOOS=windows GOARCH=386 go build

container:
	@echo "building docker container" ;\
	docker build --rm=true --tag=toposvc .
