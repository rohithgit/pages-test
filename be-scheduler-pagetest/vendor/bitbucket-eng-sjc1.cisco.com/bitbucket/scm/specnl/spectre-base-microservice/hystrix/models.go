package hystrix

import (
	"net/http"

	"golang.org/x/net/context"
)

type Hystrixer interface {
	NewSpectreHystrixInstance(req *http.Request, client *http.Client, resolveAddress bool) Hystrix
}

type HystrixerImpl struct{}

type Hystrix interface {
	HTTPClientCall(context.Context) (chan *http.Response, chan error)
	HTTPClientCallSync(context.Context) (*http.Response, error)
}

type SpectreHystrix struct {
	Timeout               int
	MaxConcurrentRequests int
	ErrorPercentThreshold int
	Client                *http.Client
	Request               *http.Request
	Command               string
	DoResolve             bool
}

// HystrixerImplMock mock struct for unit testing
type HystrixerImplMock struct {
	MockNewSpectreHystrixInstance func(req *http.Request, client *http.Client, resolveAddress bool) Hystrix
}

// SpectreHystrixMock mock struct for unit testing
type SpectreHystrixMock struct {
	MockHTTPClientCall     func(context.Context) (chan *http.Response, chan error)
	MockHTTPClientCallSync func(context.Context) (*http.Response, error)
}
