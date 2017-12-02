package hystrix

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strings"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/service_discovery"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// This function loads default configs for hystrix "command"
func loadDefaultsForHystrix(command string) {
	hystrix.ConfigureCommand(command, hystrix.CommandConfig{
		Timeout:               DEFAULT_TIMEOUT,
		MaxConcurrentRequests: hystrix.DefaultMaxConcurrent,
		ErrorPercentThreshold: hystrix.DefaultErrorPercentThreshold,
	})
}

// HTTPClientCall function wraps the client making the req in hystrix.
func HTTPClientCall(ctx context.Context, req *http.Request, client *http.Client, command string, config *hystrix.CommandConfig, resolveAddress bool) (chan *http.Response, chan error) {
	if strings.TrimSpace(command) == "" {
		command = uuid.New()
	}

	outputChan := make(chan *http.Response)
	errorsChan := make(chan error)

	if req == nil {
		go func() { errorsChan <- errors.New("Request is nil.") }()
		return nil, errorsChan
	}
	if client == nil {
		go func() { errorsChan <- errors.New("Client is nil.") }()
		return nil, errorsChan
	}

	if resolveAddress {
		sd := service_discovery.NewServiceAddress(nil)
		err := sd.ParseServerAddress(req.URL.Host)
		if err == nil {
			req.URL.Host = sd.Address
		}
		//req.URL.Host = service_discovery.ParseServerAddress(req.URL.Host, returnIP)
		log.Info("Updated Address with Port: ", req.URL.Host)
	}

	if strings.HasPrefix(req.URL.String(), "https") {
		client.Transport = &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true}, // TODO: change this to proper verification after we have the strategy in place.
			DisableCompression: true,
		}
	}

	trackingIdI := ctx.Value(TRACKING_ID_FIELD_FOR_LOG)
	if trackingIdI != nil {
		trackingId, _ := trackingIdI.(string)
		req.Header.Set("TrackingId", trackingId)
	}

	if config == nil {
		loadDefaultsForHystrix(command)
	} else {
		hystrix.ConfigureCommand(command, *config)
	}

	errorsChan = hystrix.Go(command, func() error {
		//resp, err := client.Do(req)
		resp, err := ctxhttp.Do(ctx, client, req)
		if err == nil {
			if resp != nil {
				outputChan <- resp
				return nil
			} else {
				return errors.New("Response is empty.")
			}
		}
		return err
	}, nil)
	return outputChan, errorsChan
}

// HTTPClientCallSync is the sync version of the call wrapped with hystrix
func HTTPClientCallSync(ctx context.Context, req *http.Request, client *http.Client, command string, config *hystrix.CommandConfig, resolveAddress bool) (*http.Response, error) {
	outputChan, errorsChan := HTTPClientCall(ctx, req, client, command, config, resolveAddress)
	select {
	case resp := <-outputChan:
		return resp, nil
	case err := <-errorsChan:
		return nil, err
	}
}
