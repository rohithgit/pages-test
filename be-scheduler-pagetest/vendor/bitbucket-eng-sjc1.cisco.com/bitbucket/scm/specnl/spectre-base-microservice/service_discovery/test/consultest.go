package main

import (
	"fmt"
	"github.com/benschw/dns-clb-go/clb"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/service_discovery"
)

const (
	name string = "be-tagging"
)

func main() {
	saddr := service_discovery.NewServiceAddress(clb.NewClb("127.0.0.1", "8600", clb.RoundRobin))
	err := saddr.GetServiceAddress(name)
	if err != nil {
		fmt.Printf("GetServiceAddress Errored: %v\n", err)
	}
	fmt.Printf("Service info %s - IP: %s Port: %d Address: %s Name: %s\n", name, saddr.IP, saddr.Port, saddr.Address, saddr.Name)

}
