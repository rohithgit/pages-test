package rbac

import (
	"encoding/json"
	"errors"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/context"
)

func (u *APIToken) RetrieveUser(ctx context.Context) (*UserInfo, error) {
	err := u.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &u.User, nil
}

// This method is used to get user info using the access_token using OpenID Connect
func (u *APIToken) GetUserInfo(ctx context.Context) error {
	logger := GetLoggerWithTrackingIdFromContext(ctx, logger)

	var retryCount int
	if u.TryToRefreshTokenOnFailure {
		retryCount = u.NumberOfAttempts
	} else {
		retryCount = 1
	}
	for i := 0; i < retryCount; i++ {
		//func GetUserInfo(ctx context.Context, useYourClientCreds string) (context.Context, error) {
		logger.Debugln("Get user info url:", oauthUserInfoService)

		if !govalidator.StringMatches(u.AccessToken, SPECTRE_COMMON_TOKEN_FILTER_REGEX) || len(u.AccessToken) > SPECTRE_COMMON_TOKEN_LENGTH {
			return errors.New("Access Token is invalid")
		}

		profileResponse, err := u.getCall(ctx, oauthUserInfoService)
		if err != nil {
			logger.Error("There is an error: ", err)
			var authError AuthError
			err1 := json.Unmarshal(profileResponse, &authError)
			if err1 == nil {
				logger.Error("Error is an IT error")
				logger.Debug("Attempting to refresh access token and retrying the call.")
				err2 := u.RefreshAccessToken(ctx)
				if err2 == nil {
					logger.Error("Access token was successfully refreshed. Reattempting the original call.")
					if i < retryCount {
						continue
					}
				}
			}
			return err
		}
		logger.Debugf("User info response: %#v", string(profileResponse))

		var user UserInfo
		err = json.Unmarshal(profileResponse, &user)
		if err != nil {
			return err
		}
		u.User = user
		logger.Debugf("User updated successfully: %#v", user)
		return err
	}
	return nil
}

// This method is used to get spectre roles - apps - scopes combination of data based on the requested operation
func (u *APIToken) getSpectreRoles(ctx context.Context, op Auth_Op) error {
	logger := GetLoggerWithTrackingIdFromContext(ctx, logger)

	var retryCount int
	if u.TryToRefreshTokenOnFailure {
		retryCount = u.NumberOfAttempts
	} else {
		retryCount = 1
	}
	for i := 0; i < retryCount; i++ {
		finalUrl := ""

		switch op {
		case APP_ROLES:
			finalUrl = appRoleUrl
		case APP_SCOPES:
			finalUrl = appScopeUrl
		default:
			finalUrl = appRoleScopeUrl
		}

		logger.Debugln("Get roles url:", finalUrl)

		if !govalidator.StringMatches(u.AccessToken, SPECTRE_COMMON_TOKEN_FILTER_REGEX) || len(u.AccessToken) > SPECTRE_COMMON_TOKEN_LENGTH {
			return errors.New("Access Token is invalid")
		}

		profileResponse, err := u.getCall(ctx, finalUrl)
		if err != nil {
			logger.Error("There is an error: ", err)
			var authError AuthError
			err1 := json.Unmarshal(profileResponse, &authError)
			if err1 == nil {
				logger.Error("Error is an IT error")
				logger.Debug("Attempting to refresh access token and retrying the call.")
				err2 := u.RefreshAccessToken(ctx)
				if err2 == nil {
					logger.Error("Access token was successfully refreshed. Reattempting the original call.")
					if i < retryCount {
						continue
					}
				}
			}
			return err
		}
		logger.Debugf("Get roles response: %#v", string(profileResponse))

		if u.User.SpectreRoles == nil {
			u.User.SpectreRoles = []AppRole{}
		}
		err = json.Unmarshal(profileResponse, &u.User.SpectreRoles)
		if err != nil {
			return err
		}
		logger.Debugf("Updated roles: %v", u.User.SpectreRoles)
		return nil
	}
	return nil
}

// This method is used for getting all the role - app combinations the user has access to
func (u *APIToken) GetRolesApps(ctx context.Context) error {
	return u.getSpectreRoles(ctx, APP_ROLES)
}

// This method is used for getting all the app - scopes combinations the user has access to
func (u *APIToken) GetAppScopes(ctx context.Context) error {
	return u.getSpectreRoles(ctx, APP_SCOPES)
}

// This method is used for getting all the app - role - scopes combinations the user has access to
func (u *APIToken) GetAppRoleScopes(ctx context.Context) error {
	return u.getSpectreRoles(ctx, APP_ROLE_SCOPES)
}
