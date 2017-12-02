package hystrix

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

var sampleBody = "Hello World!"

func createRequestAndCall(method, url string, body io.Reader) (*http.Client, *http.Request) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println("Error in creating a test request:", err)
	}
	return &http.Client{}, req
}

func TestAllMethods(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sampleBody)
	}))
	defer ts.Close()
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete} {
		client, req := createRequestAndCall(method, ts.URL, nil)
		respChan, errorChan := HTTPClientCall(context.TODO(), req, client, "testGet1", nil, false)
		select {
		case resp := <-respChan:
			require.NotNil(t, resp)
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			require.NotNil(t, body)
			require.Nil(t, err)
			require.Equal(t, sampleBody, strings.TrimSpace(string(body)))
		case err := <-errorChan:
			require.Nil(t, err)
		}

		client, req = createRequestAndCall(method, ts.URL, nil)
		resp, err := HTTPClientCallSync(context.TODO(), req, client, "testGet1", nil, false)
		require.NotNil(t, resp)
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.NotNil(t, body)
		require.Nil(t, err)
		require.Equal(t, sampleBody, strings.TrimSpace(string(body)))
		require.Nil(t, err)
	}
}

func TestCommandNotGiven(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sampleBody)
	}))
	defer ts.Close()

	client, req := createRequestAndCall(http.MethodGet, ts.URL, nil)
	respChan, errorChan := HTTPClientCall(context.TODO(), req, client, "", nil, false)
	select {
	case resp := <-respChan:
		require.NotNil(t, resp)
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.NotNil(t, body)
		require.Nil(t, err)
		require.Equal(t, sampleBody, strings.TrimSpace(string(body)))
	case err := <-errorChan:
		require.Nil(t, err)
	}
}

func TestConfigGiven(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sampleBody)
	}))
	defer ts.Close()
	client, req := createRequestAndCall(http.MethodPut, ts.URL, nil)
	respChan, errorChan := HTTPClientCall(context.TODO(), req, client, "blah1", &hystrix.CommandConfig{
		Timeout:               hystrix.DefaultTimeout + 2000,
		MaxConcurrentRequests: hystrix.DefaultMaxConcurrent,
		ErrorPercentThreshold: hystrix.DefaultErrorPercentThreshold,
	}, false)
	select {
	case resp := <-respChan:
		require.NotNil(t, resp)
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.NotNil(t, body)
		require.Nil(t, err)
		require.Equal(t, sampleBody, strings.TrimSpace(string(body)))
	case err := <-errorChan:
		require.Nil(t, err)
	}
}

func TestClientNotGiven(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sampleBody)
	}))
	defer ts.Close()
	_, req := createRequestAndCall(http.MethodPut, ts.URL, nil)
	respChan, errorChan := HTTPClientCall(context.TODO(), req, nil, "blah2", nil, false)
	select {
	case resp := <-respChan:
		require.Nil(t, resp)
	case err := <-errorChan:
		require.NotNil(t, err)
	}
}

func TestReqNotGiven(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sampleBody)
	}))
	defer ts.Close()
	client, _ := createRequestAndCall(http.MethodPut, ts.URL, nil)
	respChan, errorChan := HTTPClientCall(context.TODO(), nil, client, "blah3", nil, false)
	select {
	case resp := <-respChan:
		require.Nil(t, resp)
	case err := <-errorChan:
		require.NotNil(t, err)
	}
}

func TestExceedTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(10 * time.Second))
		fmt.Fprintln(w, sampleBody)
		return
	}))
	defer ts.Close()
	for i := 0; i < 5; i++ {
		client, req := createRequestAndCall(http.MethodPut, ts.URL, nil)
		client.Timeout = time.Duration(time.Millisecond * 1)
		respChan, errorChan := HTTPClientCall(context.TODO(), req, client, "blah4", nil, false)
		select {
		case resp := <-respChan:
			resp.Body.Close()
			require.Nil(t, resp)
		case err := <-errorChan:
			require.NotNil(t, err)
		}
	}
}

func TestFailure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "KABOOOOOOM!!!", http.StatusInternalServerError)
		return
	}))
	defer ts.Close()
	for i := 0; i < 5; i++ {
		client, req := createRequestAndCall(http.MethodPut, ts.URL, nil)
		respChan, errorChan := HTTPClientCall(context.TODO(), req, client, "blah5", nil, false)
		select {
		case resp := <-respChan:
			resp.Body.Close()
			require.Equal(t, 500, resp.StatusCode)
		case err := <-errorChan:
			require.NotNil(t, err)
		}
	}
}

func TestUnReachability(t *testing.T) {
	for i := 0; i < 5; i++ {
		client, req := createRequestAndCall(http.MethodGet, "http://blah.blah", nil)
		timeout := time.Duration(time.Millisecond * 10)
		client.Timeout = timeout
		respChan, errorChan := HTTPClientCall(context.TODO(), req, client, "blah5", nil, false)
		select {
		case resp := <-respChan:
			resp.Body.Close()
			require.Equal(t, resp.StatusCode, 500)
		case err := <-errorChan:
			require.NotNil(t, err)
		}

		client, req = createRequestAndCall(http.MethodGet, "http://blah.blah", nil)
		client.Timeout = timeout
		resp, err := HTTPClientCallSync(context.TODO(), req, client, "blah5", nil, false)
		require.Nil(t, resp)
		require.NotNil(t, err)
	}
}

func TestInterface(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sampleBody)
	}))
	defer ts.Close()
	client, req := createRequestAndCall(http.MethodGet, ts.URL, nil)

	h := HystrixerImpl{}
	hy := h.NewSpectreHystrixInstance(req, client, false)
	respChan, errorChan := hy.HTTPClientCall(context.TODO())
	select {
	case resp := <-respChan:
		require.NotNil(t, resp)
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.NotNil(t, body)
		require.Nil(t, err)
		require.Equal(t, sampleBody, strings.TrimSpace(string(body)))
	case err := <-errorChan:
		require.Nil(t, err)
	}

	resp1, err1 := hy.HTTPClientCallSync(context.TODO())
	require.NotNil(t, resp1)
	require.Nil(t, err1)

}