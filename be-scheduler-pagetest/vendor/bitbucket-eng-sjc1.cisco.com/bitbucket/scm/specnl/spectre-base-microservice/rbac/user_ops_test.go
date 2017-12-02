package rbac

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestToken_GetUserInfo(t *testing.T) {
	tmp := oauthUserInfoService
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("Form: %#v %s %s\n", r.Form, r.FormValue("client_id"), r.PostFormValue("client_id"))
	}))
	defer ts.Close()
	oauthUserInfoService = ts.URL

	tokens := []*APIToken{&APIToken{}, &APIToken{
		AccessToken: "asdfds",
	}, &APIToken{
		AccessToken:                "asdfds",
		TryToRefreshTokenOnFailure: true,
		NumberOfAttempts:           2,
	}}
	for _, to := range tokens {
		err := to.GetUserInfo(context.TODO())
		assert.NotNil(t, err)
	}

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"error":             "blah",
			"error_description": "blah blah",
		}
		w.WriteHeader(400)
		fmt.Println(json.NewEncoder(w).Encode(resp))
	}))
	defer ts.Close()
	oauthUserInfoService = ts.URL

	to := &APIToken{
		AccessToken:                "fdsfs",
		TryToRefreshTokenOnFailure: true,
		NumberOfAttempts:           2,
	}
	err := to.GetUserInfo(context.TODO())
	assert.NotNil(t, err)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"accesslevel":   "1",
			"city":          "DETROIT",
			"company":       "ALLY FINANCIAL INC",
			"country":       "UNITED STATES",
			"emailverified": "mppod4dev@yopmail.com",
			"familyname":    "Mishra",
			"givenname":     "Sunil",
			"orgId":         "8517088f-3144-46b4-9f04-adaa054ab797",
			"phonenumber":   "+1 453543534",
			"postalCode":    "48202",
			"roles": []string{
				"PARTNER_ADMIN",
				"PARTNER_USER",
			},
			"state":   "MI",
			"street":  "3044 W GRAND BLVD",
			"street2": "0",
			"sub":     "mppod4dev",
			"title":   "0",
		}
		fmt.Println(json.NewEncoder(w).Encode(resp))
	}))
	defer ts.Close()
	oauthUserInfoService = ts.URL

	to = &APIToken{
		AccessToken: "fdsfs",
	}
	err = to.GetUserInfo(context.TODO())
	assert.Nil(t, err)
	oauthUserInfoService = tmp
}

func TestToken_GetAppRoleScopes(t *testing.T) {
	tmp := appRoleScopeUrl
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//fmt.Printf("Form: %#v %s %s\n", r.Form, r.FormValue("client_id"), r.PostFormValue("client_id"))
	}))
	defer ts.Close()
	appRoleScopeUrl = ts.URL

	tokens := []*APIToken{&APIToken{}, &APIToken{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		err := to.GetAppRoleScopes(context.TODO())
		assert.NotNil(t, err)
	}

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"error":             "blah",
			"error_description": "blah blah",
		}
		w.WriteHeader(400)
		fmt.Println(json.NewEncoder(w).Encode(resp))
	}))
	defer ts.Close()

	to1 := &APIToken{
		AccessToken:                "fdsfs",
		TryToRefreshTokenOnFailure: true,
		NumberOfAttempts:           2,
	}
	err1 := to1.GetAppRoleScopes(context.TODO())
	assert.NotNil(t, err1)

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := []map[string]interface{}{
			{
				"role":   "r1",
				"object": "o1",
				"scopes": []string{"t1"},
			},
		}
		fmt.Println(json.NewEncoder(w).Encode(resp))
	}))
	defer ts1.Close()
	appRoleScopeUrl = ts1.URL

	to1 = &APIToken{
		AccessToken: "fdsfs",
	}
	err1 = to1.GetAppRoleScopes(context.TODO())
	assert.Nil(t, err1)
	appRoleScopeUrl = tmp
}

func TestToken_GetAppScopes(t *testing.T) {
	tmp := appScopeUrl
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("Form: %#v %s %s\n", r.Form, r.FormValue("client_id"), r.PostFormValue("client_id"))
	}))
	defer ts.Close()
	appScopeUrl = ts.URL

	tokens := []*APIToken{&APIToken{}, &APIToken{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		err := to.GetAppScopes(context.TODO())
		assert.NotNil(t, err)
	}

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := []map[string]interface{}{
			{
				"object": "o1",
				"scopes": []string{"t1"},
			},
		}
		fmt.Println(json.NewEncoder(w).Encode(resp))
	}))
	defer ts.Close()
	appScopeUrl = ts.URL

	to := &APIToken{
		AccessToken: "fdsfs",
	}
	err := to.GetAppScopes(context.TODO())
	assert.Nil(t, err)
	appScopeUrl = tmp
}

func TestToken_GetRolesApps(t *testing.T) {
	tmp := appRoleUrl
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("Form: %#v %s %s\n", r.Form, r.FormValue("client_id"), r.PostFormValue("client_id"))
	}))
	defer ts.Close()
	appRoleUrl = ts.URL

	tokens := []*APIToken{&APIToken{}, &APIToken{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		err := to.GetRolesApps(context.TODO())
		assert.NotNil(t, err)
	}

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := []map[string]interface{}{
			{
				"role":   "r1",
				"object": "o1",
			},
		}
		fmt.Println(json.NewEncoder(w).Encode(resp))
	}))
	defer ts.Close()
	appRoleUrl = ts.URL

	to := &APIToken{
		AccessToken: "fdsfs",
	}
	err := to.GetRolesApps(context.TODO())
	assert.Nil(t, err)
	appRoleUrl = tmp
}
