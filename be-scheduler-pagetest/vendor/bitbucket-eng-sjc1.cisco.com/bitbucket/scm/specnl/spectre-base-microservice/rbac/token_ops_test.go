package rbac

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestToken_GetAccessToken(t *testing.T) {
	tmp := oauthAccessTokenService
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("Form: %#v %s %s\n", r.Form, r.FormValue("client_id"), r.PostFormValue("client_id"))
	}))
	defer ts.Close()
	oauthAccessTokenService = ts.URL

	tokens := []*APIToken{&APIToken{}, &APIToken{
		ClientSecret:       "fdssdf",
		UseYourClientCreds: true,
	}, &APIToken{
		ClientId:           "dfgf",
		UseYourClientCreds: true,
	}, &APIToken{
		ClientId:           "fdsfs",
		ClientSecret:       "sdfdsf",
		UseYourClientCreds: true,
	}}
	for _, to := range tokens {
		err := to.GetAccessToken(context.TODO())
		assert.NotNil(t, err)
	}

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"access_token":  "sfdfsfsfsfsdfsf",
			"refresh_token": "fsasdfsdgdffghfg",
			"expires_in":    450,
		}
		fmt.Println(json.NewEncoder(w).Encode(resp))
	}))
	defer ts.Close()
	oauthAccessTokenService = ts.URL

	to := &APIToken{
		ClientId:           "fdsfs",
		ClientSecret:       "sdfdsf",
		UseYourClientCreds: true,
	}
	err := to.GetAccessToken(context.TODO())
	assert.Nil(t, err)

	oauthAccessTokenService = tmp
}

func TestToken_RefreshAccessToken(t *testing.T) {
	tmp := oauthRefreshTokenService
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("Form: %#v %s %s\n", r.Form, r.FormValue("client_id"), r.PostFormValue("client_id"))
	}))
	defer ts.Close()
	oauthRefreshTokenService = ts.URL

	tokens := []*APIToken{&APIToken{}, &APIToken{
		ClientSecret: "fdssdf",
	}, &APIToken{
		ClientId: "dfgf",
	}, &APIToken{
		RefreshToken: "asdfsfs",
	}}
	for _, to := range tokens {
		err := to.RefreshAccessToken(context.TODO())
		assert.NotNil(t, err)
	}

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"access_token":  "sfdfsfsfsfsdfsf",
			"refresh_token": "fsasdfsdgdffghfg",
			"expires_in":    450,
		}
		fmt.Println(json.NewEncoder(w).Encode(resp))
	}))
	defer ts.Close()
	oauthRefreshTokenService = ts.URL

	to := &APIToken{
		RefreshToken: "asdfsfs",
	}
	err := to.RefreshAccessToken(context.TODO())
	assert.Nil(t, err)

	oauthRefreshTokenService = tmp
}

func TestToken_RevokeToken(t *testing.T) {
	tmp := oauthRevokeTokenService
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("Form: %#v %s %s\n", r.Form, r.FormValue("client_id"), r.PostFormValue("client_id"))
	}))
	defer ts.Close()
	oauthRevokeTokenService = ts.URL

	tokens := []*APIToken{&APIToken{}, &APIToken{
		ClientSecret: "fdssdf",
	}, &APIToken{
		ClientId: "dfgf",
	}, &APIToken{
		RefreshToken: "asdfsfs",
	}}
	for _, to := range tokens {
		err := to.RevokeToken(context.TODO())
		assert.NotNil(t, err)
	}

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "KABOOM!", http.StatusInternalServerError)
	}))
	defer ts.Close()
	oauthRevokeTokenService = ts.URL

	to := &APIToken{
		AccessToken: "asdfsfs",
	}
	err := to.RevokeToken(context.TODO())
	assert.NotNil(t, err)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	defer ts.Close()
	oauthRevokeTokenService = ts.URL

	to = &APIToken{
		AccessToken: "asdfsfs",
		ExpiresIn: 450,
		User: UserInfo{
			Username: "hello",
		},
	}
	err = to.RevokeToken(context.TODO())
	assert.Nil(t, err)
	oauthRevokeTokenService = tmp
}

func TestToken_ValidateTokenAndScope(t *testing.T) {
	tmp := oauthValidationService
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("Form: %#v %s %s\n", r.Form, r.FormValue("client_id"), r.PostFormValue("client_id"))
	}))
	defer ts.Close()
	oauthValidationService = ts.URL

	tokens := []*APIToken{&APIToken{}, &APIToken{
		ClientSecret: "fdssdf",
	}, &APIToken{
		ClientId: "dfgf",
	}, &APIToken{
		RefreshToken: "asdfsfs",
	}, &APIToken{
		AccessToken: "asdfsfs",
	}}
	for _, to := range tokens {
		_, err := to.ValidateToken(context.TODO())
		assert.NotNil(t, err)

		_, err = to.ValidateScope(context.TODO(), "t1")
		assert.NotNil(t, err)

		_, err = to.ValidateScope(context.TODO(), "")
		assert.NotNil(t, err)
	}

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"access_token":  map[string]interface{}{
				"uid": "test",
				"token": "sdfdsfsdfsfsdfsdf",
			},
			"refresh_token": "fsasdfsdgdffghfg",
			"scope": "t1 t2",
			"expires_in":    450,
		}
		fmt.Println(json.NewEncoder(w).Encode(resp))
	}))
	defer ts.Close()
	oauthValidationService = ts.URL

	to := &APIToken{
		AccessToken: "asdfsfs",
	}
	v, err := to.ValidateToken(context.TODO())
	assert.Nil(t, err)
	assert.True(t, v)

	v, err = to.ValidateScope(context.TODO(), "t1")
	assert.Nil(t, err)
	assert.True(t, v)

	v, err = to.ValidateScope(context.TODO(), "t3")
	assert.Nil(t, err)
	assert.False(t, v)

	oauthValidationService = tmp
}