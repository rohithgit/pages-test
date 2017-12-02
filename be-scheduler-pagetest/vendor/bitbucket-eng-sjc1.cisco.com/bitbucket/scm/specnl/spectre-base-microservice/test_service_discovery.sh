#!/bin/bash
#---- Consul Docker RUN--------------------
(docker stop TestDis) || true
(docker rm TestDis) || true
(docker rmi test1) || true
(docker stop consul1) || true
(docker rm consul1) || true

docker run -d -p 8400:8400 -p 8500:8500 -p 8600:53/udp --name consul1 progrium/consul -server -bootstrap
sleep 5
#-----Test Service Registry------------------------
curl -X PUT -d '{"Datacenter": "dc1", "Node": "google", "Address": "127.0.0.1", "Service": {"Service": "search", "Port": 80}}' http://127.0.0.1:8500/v1/catalog/register

#------Grab IP of consul1
IP=$(docker inspect -f '{{.NetworkSettings.IPAddress }}' consul1)

#-------Check Test Service Registry-----------------
#dig @$IP search.service.consul

#------Test Code Docker Files-----------------------
docker build -t test1 -f Dockerfile_service_discovery .

docker run -it --name TestDis --dns=$IP --dns-search=service.consul test1 go test -v -cover

(docker stop TestDis) || true
(docker rm TestDis) || true
(docker rmi test1) || true
(docker stop consul1) || true
(docker rm consul1) || true