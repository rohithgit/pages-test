package utils

import (
	"fmt"
	"net"
	"strings"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/service_discovery"
	"golang.org/x/net/context"
)

func ParseServerAddress(ctx context.Context, server string, returnIP bool) string {
	SpectreLog.Debug("Entering ParseServerAddress for service discovery")

	sd := service_discovery.NewServiceAddress(nil)
	err := sd.ParseServerAddress(server)
	if err != nil {
		SpectreLog.Errorf("Service discovery error: %v\n", err)
		return server
	}

	if returnIP {
		SpectreLog.Info("Retrieved address from service discovery", sd.Address)
		return sd.Address
	} else {
		serviceIP := strings.Split(sd.Address, ":")
		SpectreLog.Infof("Retrieving address service ip %s and port %d",serviceIP[0], sd.Port )
		return net.JoinHostPort(serviceIP[0], fmt.Sprint(sd.Port))
	}
}
