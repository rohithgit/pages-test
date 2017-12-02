package service_discovery

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log"
	"github.com/benschw/dns-clb-go/clb"
	"os"
	"time"
)

const HOSTED_ZONE = "ecs.internal"
const AWS_METADATA_URL = "http://169.254.169.254/latest/meta-data/iam"

// ServiceAddress contains IP, Port number and address of service
// It also provides service address and server utility functions

type ServiceAddress struct {
	IP          string
	Port        uint16
	Address     string
	Name        string
	NameAndPort string
	lb          clb.LoadBalancer
}

func getDiscoveryProperties() map[string]string {

	pers_properties := map[string]string{
		"zookeeper": "zookeeper.datalake.internal:2181",
		"kafka":     "kafka.datalake.internal:9092",
		"mongodb":   "mongodb.datalake.internal:27017",
		"postgres":  "postgres.datalake.internal:5432",
		"elasticache":"elasticache.datalake.internal:6379",
		"vault" : "vault.datalake.internal:8200",
		"elasticsearch" : "elasticsearch.datalake.internal:9200",
	}
	return pers_properties
}

func GetServiceName(targetService string) string {
	log.Info("GetServiceName : with target service : " + targetService)
	var serviceRet string = ""
	var tier string = ""
	parts := strings.Split(targetService, "-")

	//println(len(parts))
	if len(parts) > 1 {
		tier = parts[0]
		log.Info("microservice tier : " + tier)
		serviceRet = GetServiceNameWithTier(targetService, tier)
	} else {
		serviceRet = GetServiceNameWithTier(targetService, "pers")
	}
	//println(" serviceRet" , serviceRet)

	log.Info("Service Name returned --> " + serviceRet)
	return serviceRet
}

func GetServiceNameWithTier(targetService string, tier string) string {
	log.Infof("GetServiceNamewithTier() with targetService : %s  and tier : %s",targetService,tier)
	var serviceRet string = ""
	var sysEnvDomain = os.Getenv("domain")
	var sysEnvDomainStateful = os.Getenv("domain_stateful")
	var sysEnvDomainStateless = os.Getenv("domain_stateless")
	discoveryProps := getDiscoveryProperties()
	serviceDetails, ok := discoveryProps[targetService]
	if ok {
		if sysEnvDomainStateful != "" {
			var parts = strings.Split(serviceDetails, ":")
			if len(parts) >1 {
				var persPort = parts[1]
				serviceRet = targetService + "." + sysEnvDomainStateful + ":" + persPort
			}

		} else {

			serviceRet = serviceDetails
		}
	} else if tier != "pers" {
		log.Info("Tier not equal to pers")
		if sysEnvDomainStateless != "" {
			log.Info("Got Domain " + sysEnvDomain)
			serviceRet = targetService + "." + tier + "." + sysEnvDomainStateless + ":443"
		} else {
			log.Info("Using default domain " + HOSTED_ZONE)
			serviceRet = targetService + "." + tier + "." + HOSTED_ZONE + ":443"

		}
	}
	log.Info("GetServiceNamewithTier returns : " + serviceRet)
	return serviceRet
}

// NewServiceAddress returns a new ServiceAddress object
// If no CLB is supplied then a new CLB is instantiated

func NewServiceAddress(lb clb.LoadBalancer) *ServiceAddress {
	saddr := new(ServiceAddress)
	if lb == nil {
		saddr.lb = clb.New()
	} else {
		saddr.lb = lb
	}
	return saddr
}

// GetServiceAddress parses the service name, port number and address from
// consul address of the service. We will check for the environment here and request the service accordingly.
func (ref *ServiceAddress) GetServiceAddress(name string) error {
	log.Info("GetServiceAddress called with : " + name)
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(AWS_METADATA_URL)
	log.Infof("resp: %#v, err: %v", resp, err)
	if err != nil || resp.StatusCode != 200 {
		// handle err
		ref.Name = name
		address, errMantl := ref.lb.GetAddress(name)
		if errMantl != nil {
			log.Error("Service discovery lookup error for service: ", name, " - ", errMantl.Error())
			return errMantl
		}

		log.Infof("Mantl Info obtained from Service discovery for service: %s - Address: %s Port: %d\n", name, address.Address, address.Port)

		ref.IP = address.Address
		ref.Port = address.Port
		ref.Address = net.JoinHostPort(address.Address, fmt.Sprintf("%d", address.Port))
		ref.NameAndPort = net.JoinHostPort(ref.Name, fmt.Sprintf("%d", address.Port))
		return nil
	} else {
		log.Info("AWS environment detected")
		awsRet := GetServiceName(name)
		np := strings.Split(awsRet, ":")
		ref.Name = np[0]
		if len(np) > 1 {
			p, err := strconv.ParseUint(np[1], 10, 16)
			if err != nil {
				return err
			}
			ref.Port = uint16(p)
		}
		ref.NameAndPort = awsRet
		ref.Address = awsRet
		log.Infof("Printing out structure values  ref.Address : %s  , ref.Name : %s , ref.NameAndPort : %s", ref.Address, ref.Name, ref.NameAndPort )
		log.Infof("AWS Info obtained from Service discovery for service: %s - Address: %s Port: %d", name, ref.Address, ref.Port)
		return nil
	}
	return nil
}

// ParseServerAddress parses the server address into IP(IP), Port, Addr
func (ref *ServiceAddress) ParseServerAddress(server string) error {
	log.Info("ParseServerAddress() called with :" + server)
	ref.Address = server
	serviceIP := strings.Split(server, ":")
	ref.IP = serviceIP[0]
	if len(serviceIP) > 1 {
		p, err := strconv.ParseUint(serviceIP[1], 10, 16)
		if err != nil {
			return err
		}
		ref.Port = uint16(p)
	}
	checkIP := net.ParseIP(serviceIP[0])
	if checkIP.To4() == nil {
		log.Info("Service provided is not an IP: ", server)
		log.Info("Calling ref.GetServiceAddress() with serviceIP[0] : " + serviceIP[0]  )
		if err := ref.GetServiceAddress(serviceIP[0]); err != nil {
			return err
		}
	} else {
		log.Info("Service provided is an IP: ", server)
	}
	return nil
}