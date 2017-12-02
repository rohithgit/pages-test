package rbac

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log"
	"github.com/asaskevich/govalidator"
	"golang.org/x/net/context"
)

var (
	logger                            = log.NewEntry(log.New())
	SPECTRE_COMMON_TOKEN_FILTER_REGEX = `^[0-9_\w\-:\.]+$`
	SPECTRE_COMMON_TOKEN_LENGTH       = 20000
	SPECTRE_COMMON_FILTER_REGEX       = `^[0-9,\w\-\.]+$`
	TRACKING_ID_FIELD_FOR_LOG         = "trackingId"
	CUSTOMER_ID_FIELD_FOR_LOG         = "customerId"
)

var rbacService, oauthAccessTokenService, oauthValidationService, oauthRefreshTokenService, oauthRevokeTokenService, oauthUserInfoService, appRoleUrl, appScopeUrl, appRoleScopeUrl string
var proto = "https"

func init() {
	rbacService = os.Getenv("RBAC_SERVICE")
	oauthAccessTokenService = proto + "://" + rbacService + "/api/v1/access_token"
	oauthRefreshTokenService = proto + "://" + rbacService + "/api/v1/refresh_token"
	oauthValidationService = proto + "://" + rbacService + "/api/v1/validate_token"
	oauthRevokeTokenService = proto + "://" + rbacService + "/api/v1/revoke_token"
	oauthUserInfoService = proto + "://" + rbacService + "/api/v1/userauthinfo"

	appRoleUrl = proto + "://" + rbacService + "/api/v1/user_app_roles"
	appScopeUrl = proto + "://" + rbacService + "/api/v1/user_app_scopes"
	appRoleScopeUrl = proto + "://" + rbacService + "/api/v1/user_app_detailed_roles"
}

// This function helps with capturing the access_token header from the header instance, create a Token object and loads it into the context
func (r RBACImpl) NewTokenFromRequestHeader(ctx context.Context, header http.Header) (context.Context, error) {
	var tokenParam string
	authorization := header.Get("Authorization")
	if strings.TrimSpace(authorization) != "" && strings.Contains(authorization, "Bearer") {
		auths := strings.Split(authorization, " ")
		if len(auths) == 2 {
			tokenParam = auths[1]
		}
	}
	if tokenParam == "" {
		return ctx, errors.New("No token in the header.")
	}

	token := &APIToken{
		AccessToken: tokenParam,
	}

	trackingId := header.Get("TrackingId")
	if strings.TrimSpace(trackingId) != "" {
		ctx = context.WithValue(ctx, TRACKING_ID_FIELD_FOR_LOG, trackingId)
	}

	ctx = r.NewContextWithToken(ctx, token)
	return ctx, nil
}

// This function helps to create a Token object using the provided client id and secret, and loads it into the context
func (r RBACImpl) NewTokenFromClientId(ctx context.Context, clientId, clientSecret string, scopes []string) (context.Context, error) {
	if !govalidator.StringMatches(clientId, SPECTRE_COMMON_TOKEN_FILTER_REGEX) || len(clientId) > SPECTRE_COMMON_TOKEN_LENGTH {
		return ctx, errors.New("Client id is invalid")
	}
	if !govalidator.StringMatches(clientSecret, SPECTRE_COMMON_TOKEN_FILTER_REGEX) || len(clientId) > SPECTRE_COMMON_TOKEN_LENGTH {
		return ctx, errors.New("Client secret is invalid")
	}
	token := &APIToken{
		ClientId:     clientId,
		ClientSecret: clientSecret,

		Scopes: scopes,

		UseYourClientCreds: true,
	}

	ctx = r.NewContextWithToken(ctx, token)
	return ctx, nil
}
