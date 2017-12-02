package rbac

import (
	"errors"

	"golang.org/x/net/context"
)

// This is a wrapper to return a new context with the provided Token instance
func (r RBACImpl) NewContextWithToken(ctx context.Context, u Token) context.Context {
	return context.WithValue(ctx, tKey, u)
}

// This is a wrapper to return the Token instance from the provided context
func (r RBACImpl) TokenFromContext(ctx context.Context) (Token, bool) {
	t, ok := ctx.Value(tKey).(Token)
	return t, ok
}

// This is wrapper for several Token operations. It takes the token from the context, performs the necessary operation and injects it back into the context.
func (r RBACImpl) TokenOp(ctx context.Context, op Auth_Op) (context.Context, error) {
	u, ok := r.TokenFromContext(ctx)
	if !ok {
		return ctx, errors.New("No token in context")
	}
	var err error
	switch op {
	case GET_TOKEN:
		err = u.GetAccessToken(ctx)
	case REFRESH_TOKEN:
		err = u.RefreshAccessToken(ctx)
	case REVOKE_TOKEN:
		err = u.RevokeToken(ctx)
	case USER_INFO:
		err = u.GetUserInfo(ctx)
	case APP_ROLES:
		err = u.GetRolesApps(ctx)
	case APP_SCOPES:
		err = u.GetAppScopes(ctx)
	case APP_ROLE_SCOPES:
		err = u.GetAppRoleScopes(ctx)
	}

	if err != nil {
		return ctx, err
	}
	return r.NewContextWithToken(ctx, u), nil
}

// This is a wrapper for Validate if the token in context is valid
func (r RBACImpl) ValidateTokenOp(ctx context.Context) (bool, error) {
	u, ok := r.TokenFromContext(ctx)
	if !ok {
		return false, errors.New("No token in context")
	}
	return u.ValidateToken(ctx)
}

// This is a wrapper for Validating if a token has authorization for a provided scope
func (r RBACImpl) ValidateScopeOp(ctx context.Context, scope string) (bool, error) {
	u, ok := r.TokenFromContext(ctx)
	if !ok {
		return false, errors.New("No token in context")
	}
	return u.ValidateScope(ctx, scope)
}
