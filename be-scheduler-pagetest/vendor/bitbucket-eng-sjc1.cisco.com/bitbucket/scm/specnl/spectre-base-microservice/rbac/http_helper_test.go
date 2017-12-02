package rbac

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log"
	"golang.org/x/net/context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"time"
)

func init() {
	logger.Level = log.DebugLevel
}

func TestToken_getAndPostCall(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("Form: %#v %s %s\n", r.Form, r.FormValue("client_id"), r.PostFormValue("client_id"))
	}))
	defer ts.Close()

	tokens := []*APIToken{&APIToken{},  &APIToken{
		UseYourClientCreds: true,
	}}
	for _, to := range tokens {
		_, err := to.getCall(context.TODO(), ts.URL)
		assert.Nil(t, err)
		_, err = to.postCall(context.TODO(), ts.URL, nil, nil)
		assert.Nil(t, err)
	}
}
func TestToken_getAndPostCallError(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "KABOOM!", http.StatusInternalServerError)
	}))
	defer ts.Close()

	tokens := []*APIToken{&APIToken{},  &APIToken{
		UseYourClientCreds: true,
	}}
	for _, to := range tokens {
		_, err := to.getCall(context.TODO(), ts.URL)
		assert.NotNil(t, err)
		_, err = to.postCall(context.TODO(), ts.URL, nil, nil)
		assert.NotNil(t, err)
	}
}

func TestToken_Timeout(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep((TIMEOUT + 1) * time.Second)
	}))
	defer ts.Close()

	tokens := []*APIToken{&APIToken{},  &APIToken{
		UseYourClientCreds: true,
	}}
	for _, to := range tokens {
		_, err := to.getCall(context.TODO(), ts.URL)
		assert.NotNil(t, err)
		_, err = to.postCall(context.TODO(), ts.URL, nil, nil)
		assert.NotNil(t, err)
	}
}
