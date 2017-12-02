package hystrix

import (
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/net/context"
)

var TRACKING_ID_FIELD_FOR_LOG = "trackingId"
var DEFAULT_TIMEOUT = 10 * 1000

func (h HystrixerImpl) NewSpectreHystrixInstance(req *http.Request, client *http.Client, resolveAddress bool) Hystrix {
	return &SpectreHystrix{
		Timeout:               DEFAULT_TIMEOUT,
		MaxConcurrentRequests: hystrix.DefaultMaxConcurrent,
		ErrorPercentThreshold: hystrix.DefaultErrorPercentThreshold,

		Request:   req,
		Client:    client,
		DoResolve: resolveAddress,
	}
}

func (s *SpectreHystrix) loadDefaults() {
	if s.Timeout <= 0 {
		s.Timeout = DEFAULT_TIMEOUT
	}
	if s.MaxConcurrentRequests <= 0 {
		s.MaxConcurrentRequests = hystrix.DefaultMaxConcurrent
	}
	if s.ErrorPercentThreshold <= 0 {
		s.ErrorPercentThreshold = hystrix.DefaultErrorPercentThreshold
	}
}

func (s *SpectreHystrix) HTTPClientCall(ctx context.Context) (chan *http.Response, chan error) {
	s.loadDefaults()
	return HTTPClientCall(ctx, s.Request, s.Client, s.Command, &hystrix.CommandConfig{
		Timeout:               s.Timeout,
		MaxConcurrentRequests: s.MaxConcurrentRequests,
		ErrorPercentThreshold: s.ErrorPercentThreshold,
	}, s.DoResolve)
}

func (s *SpectreHystrix) HTTPClientCallSync(ctx context.Context) (*http.Response, error) {
	s.loadDefaults()
	return HTTPClientCallSync(ctx, s.Request, s.Client, s.Command, &hystrix.CommandConfig{
		Timeout:               s.Timeout,
		MaxConcurrentRequests: s.MaxConcurrentRequests,
		ErrorPercentThreshold: s.ErrorPercentThreshold,
	}, s.DoResolve)
}

// HTTPClientCall
func (s *SpectreHystrixMock) HTTPClientCall(ctx context.Context) (chan *http.Response, chan error) {
	return s.MockHTTPClientCall(ctx)
}

// HTTPClientCallSync
func (s *SpectreHystrixMock) HTTPClientCallSync(ctx context.Context) (*http.Response, error) {
	return s.MockHTTPClientCallSync(ctx)
}

// NewSpectreHystrixInstance
func (h HystrixerImplMock) NewSpectreHystrixInstance(req *http.Request, client *http.Client, resolveAddress bool) Hystrix {
	return h.MockNewSpectreHystrixInstance(req, client, resolveAddress)
}
