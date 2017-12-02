package rbac

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"errors"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/hystrix"
	hyx "github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/net/context"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log"
)

// This is a helper method to perform HTTP GET call with the token
func (u *APIToken) getCall(ctx context.Context, urlToCall string) ([]byte, error) {
	params := url.Values{}
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+u.AccessToken)
	return u.createRequest(ctx, http.MethodGet, urlToCall, params, headers)
}

// This is a helper method to perform HTTP POST call with the token
func (u *APIToken) postCall(ctx context.Context, urlToCall string, params url.Values, headers http.Header) ([]byte, error) {
	return u.createRequest(ctx, http.MethodPost, urlToCall, params, headers)
}

// This is a helper method for creating a HTTP Request of a provided method and return the response body
func (u *APIToken) createRequest(ctx context.Context, method string, urlToCall string, params url.Values, headers http.Header) ([]byte, error) {
	logger := GetLoggerWithTrackingIdFromContext(ctx, logger)
	logger.Debugf("Method: %s, url: %s, params: %v, headers: %v\n", method, urlToCall, params, headers)

	if params == nil {
		params = url.Values{}
	}
	//params.Set("token", u.AccessToken)
	if u.UseYourClientCreds {
		params.Set("client_id", u.ClientId)
		params.Set("client_secret", u.ClientSecret)
	}
	req, err := http.NewRequest(method, urlToCall, strings.NewReader(params.Encode()))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if headers != nil {
		for k, _ := range headers {
			req.Header.Set(k, headers.Get(k))
		}
	}
	//req.Header.Set("Authorization", "Bearer "+u.AccessToken)
	if method == http.MethodPost {
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	}
	return httpCallHelper(ctx, req)
}

// This is a helper function for making the HTTP call with the provided request instance and returns the response body
func httpCallHelper(ctx context.Context, req *http.Request) ([]byte, error) {
	logger := GetLoggerWithTrackingIdFromContext(ctx, logger)

	client := cleanhttp.DefaultClient()
	client.Timeout = TIMEOUT * time.Second

	if req.URL.Scheme == "https" {
		client.Transport = &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true}, // TODO: remove this after we have the strategy in place.
			DisableCompression: true,
		}
	}

	// moved to hystrix base code
	//trackingIdI := ctx.Value(TRACKING_ID_FIELD_FOR_LOG)
	//if trackingIdI != nil {
	//	trackingId, _ := trackingIdI.(string)
	//	req.Header.Set("TrackingId", trackingId)
	//}

	//resp, err := ctxhttp.Do(ctx, client, req)
	resp, err := hystrix.HTTPClientCallSync(ctx, req, client, "", &hyx.CommandConfig{
		Timeout:               1000 * TIMEOUT,
		MaxConcurrentRequests: hyx.DefaultMaxConcurrent,
		ErrorPercentThreshold: hyx.DefaultErrorPercentThreshold,
	}, true)
	if err != nil {
		logger.Error("http call error: ", err)
		return nil, err
	}
	logger.Debug("http status code: ", resp.StatusCode)

	defer resp.Body.Close()
	rawResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Response body read error: ", err)
		return nil, err
	}
	logger.Debug("Response body: ", string(rawResp))
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(rawResp))
	}
	return rawResp, nil
}


func GetLoggerWithTrackingIdFromContext(ctx context.Context, log *log.Entry) *log.Entry {
	val := ctx.Value(TRACKING_ID_FIELD_FOR_LOG)
	if val != nil {
		log = log.WithField(TRACKING_ID_FIELD_FOR_LOG, val)
	}
	cus := ctx.Value(CUSTOMER_ID_FIELD_FOR_LOG)
	if cus != nil {
		log = log.WithField(CUSTOMER_ID_FIELD_FOR_LOG, cus)
	}
	return log
}