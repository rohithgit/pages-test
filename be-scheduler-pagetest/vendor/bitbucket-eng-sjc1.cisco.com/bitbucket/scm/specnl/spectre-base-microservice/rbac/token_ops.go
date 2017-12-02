package rbac

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/context"
)

// This method validates if the token is valid and also verify if the token has the authorization for the provided scope
func (u *APIToken) ValidateScope(ctx context.Context, scope string) (bool, error) {
	isValid, err := u.ValidateToken(ctx)
	if !isValid {
		return isValid, err
	}
	for _, v := range u.Scopes {
		if scope == v {
			return true, nil
		}
	}
	return false, nil
}

// This method validates if the token is valid
func (u *APIToken) ValidateToken(ctx context.Context) (bool, error) {
	logger := GetLoggerWithTrackingIdFromContext(ctx, logger)
	//func ValidateToken(ctx context.Context) (context.Context, bool, error) {
	logger.Debugln("validation service url:", oauthValidationService)

	if !govalidator.StringMatches(u.AccessToken, SPECTRE_COMMON_TOKEN_FILTER_REGEX) || len(u.AccessToken) > SPECTRE_COMMON_TOKEN_LENGTH {
		return false, errors.New("Access Token is invalid")
	}

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+u.AccessToken)
	profileResponse, err := u.postCall(ctx, oauthValidationService, nil, headers)
	if err != nil {
		return false, err
	}
	logger.Debugf("Profile response from validation: %#v", string(profileResponse))

	// Validation has a slightly different format, access_token will be a map.
	var resp map[string]interface{}
	err = json.Unmarshal(profileResponse, &resp)
	if err != nil {
		return false, err
	}

	access_token, ok := resp["access_token"]
	if ok {
		tokenData, ok := access_token.(map[string]string)
		if ok {
			u.User.Username, _ = tokenData["uid"]
		}
	}

	scope, ok := resp["scope"]
	if ok {
		u.Scope, _ = scope.(string)
	}
	expIn, ok := resp["expires_in"]
	if ok {
		u.ExpiresIn, _ = expIn.(float64)
	}

	if u.ExpiresIn != 0 {
		u.ExpiryTime = time.Now().Add(time.Duration(u.ExpiresIn) * time.Second)
	}

	if u.Scope != "" {
		u.Scopes = strings.Split(u.Scope, " ")
	}
	return true, nil
}

// This method is used to get an access token with the client id and secret using client credentials grant type
func (u *APIToken) GetAccessToken(ctx context.Context) error {
	if !govalidator.StringMatches(u.ClientId, SPECTRE_COMMON_TOKEN_FILTER_REGEX) || len(u.ClientId) > SPECTRE_COMMON_TOKEN_LENGTH{
		return errors.New("Client id is invalid")
	}
	if !govalidator.StringMatches(u.ClientSecret, SPECTRE_COMMON_TOKEN_FILTER_REGEX) || len(u.ClientSecret) > SPECTRE_COMMON_TOKEN_LENGTH{
		return errors.New("Client secret is invalid")
	}
	params := url.Values{}
	if len(u.Scopes) > 0 {
		params.Set("scope", strings.Join(u.Scopes, " "))
	}
	profileResponse, err := u.postCall(ctx, oauthAccessTokenService, params, nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(profileResponse, u)
	if err != nil {
		return err
	}

	if u.ExpiresIn != 0 {
		u.ExpiryTime = time.Now().Add(time.Duration(u.ExpiresIn) * time.Second)
	}

	return nil
}

// This method is used to refresh the access token using a related refresh token
func (u *APIToken) RefreshAccessToken(ctx context.Context) error {
	if !govalidator.StringMatches(u.RefreshToken, SPECTRE_COMMON_TOKEN_FILTER_REGEX) || len(u.RefreshToken) > SPECTRE_COMMON_TOKEN_LENGTH{
		return errors.New("Refresh Token is invalid")
	}
	params := url.Values{}
	params.Set("refresh_token", u.RefreshToken)
	profileResponse, err := u.postCall(ctx, oauthRefreshTokenService, params, nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(profileResponse, u)
	if err != nil {
		return err
	}

	if u.ExpiresIn != 0 {
		u.ExpiryTime = time.Now().Add(time.Duration(u.ExpiresIn) * time.Second)
	}

	return nil
}

// This method is used to revoke the token
func (u *APIToken) RevokeToken(ctx context.Context) error {
	logger := GetLoggerWithTrackingIdFromContext(ctx, logger)

	logger.Debugln("Revoke service url:", oauthRevokeTokenService)

	if !govalidator.StringMatches(u.AccessToken, SPECTRE_COMMON_TOKEN_FILTER_REGEX) || len(u.AccessToken) > SPECTRE_COMMON_TOKEN_LENGTH{
		return errors.New("Access Token is invalid")
	}

	params := url.Values{}
	params.Set("token", u.AccessToken)

	_, err := u.postCall(ctx, oauthRevokeTokenService, params, nil)
	if err != nil {
		//return err
		logger.Error("Revoke call error: ", err)
	}

	u.AccessToken = ""
	u.RefreshToken = ""
	u.ExpiresIn = 0
	u.ExpiryTime = time.Now().Add(-1 * time.Second)
	u.Scope = ""
	u.Scopes = []string{}
	u.User = UserInfo{}

	return err
}

func (u *APIToken) RetrieveAccessToken() string {
	return u.AccessToken
}