package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"

	"etcd"
	"logging"
)

var (
	certFile = "/certs/etcd-client.crt"
	keyFile  = "/certs/etcd-client.key"
	caFile   = "/certs/ca.crt"
)

var log = logging.Log.Logger

func init() {
	// Output to stderr instead of stdout, could also be a file.
	// Only log the warning severity or above.
	logging.Log.LoggingInit("Text", os.Stdout, "Info")
}

func main() {
	var err error
	var SpectreEtcdClient *etcd.SpectreEtcd

	log.Info("Trying to establish etcd connection.")

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	SpectreEtcdClient, err = etcd.NewETCDClient(
		[]string{"https://etcd0.cisco.com:2379", "https://etcd1.cisco.com:2379", "https://etcd2.cisco.com:2379"},
		//[]string{"https://127.0.0.1:4001"},
		transport)
	if err != nil {
		log.Error("Connection error:", err)
	} else {
		log.Info("etcd connection established")

		_, err = SpectreEtcdClient.SetKeyValue("/foo1", "bar1")
		if err == nil {
			_, err := SpectreEtcdClient.GetStringValue("/foo1")
			if err == nil {
				_, err = SpectreEtcdClient.DeleteKey("/foo1")
			}
		}
	}
	if err != nil {
		log.Error("Error:", err)
	}
}
