package rbac

import (
	"time"

	"errors"
	"net/http"

	"golang.org/x/net/context"
)

const (
	TIMEOUT = 10 // sec

	uKey privateKey = 1919
	tKey privateKey = 1920
)

type RBAC interface {
	NewTokenFromRequestHeader(ctx context.Context, header http.Header) (context.Context, error)
	NewTokenFromClientId(ctx context.Context, clientId, clientSecret string, scopes []string) (context.Context, error)
	NewContextWithToken(ctx context.Context, u Token) context.Context
	TokenFromContext(ctx context.Context) (Token, bool)
	TokenOp(ctx context.Context, op Auth_Op) (context.Context, error)
	ValidateTokenOp(ctx context.Context) (bool, error)
	ValidateScopeOp(ctx context.Context, scope string) (bool, error)
}

type RBACImpl struct {
}

type AuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type Token interface {
	RetrieveUser(ctx context.Context) (*UserInfo, error)
	RetrieveAccessToken() string
	GetAccessToken(ctx context.Context) error
	ValidateToken(ctx context.Context) (bool, error)
	ValidateScope(ctx context.Context, scope string) (bool, error)
	RefreshAccessToken(ctx context.Context) error
	RevokeToken(ctx context.Context) error

	GetUserInfo(ctx context.Context) error

	GetRolesApps(ctx context.Context) error
	GetAppScopes(ctx context.Context) error
	GetAppRoleScopes(ctx context.Context) error
}

type APIToken struct {
	ClientId     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiryTime   time.Time `json:"-"`
	ExpiresIn    float64   `json:"expires_in"`
	Scopes       []string  `json:"-"`
	Scope        string    `json:"scope"`
	User         UserInfo  `json:"-"`

	UseYourClientCreds         bool `json:"-"`
	TryToRefreshTokenOnFailure bool `json:"-"`
	NumberOfAttempts           int  `json:"-"`
}

type UserInfo struct {
	Username string `json:"sub"`

	City       string `json:"city"`
	Company    string `json:"company"`
	Country    string `json:"country"`
	Email      string `json:"emailverified"`
	LastName   string `json:"familyname"`
	FirstName  string `json:"givenname"`
	OrgID      string `json:"orgId"`
	Phone      string `json:"phonenumber"`
	PostalCode string `json:"postalCode"`
	State      string `json:"state"`
	Street     string `json:"street"`
	Street2    string `json:"street2"`

	SpectreRoles []AppRole `json:"-"`
	Roles        []string  `json:"roles"`
}

type AppRole struct {
	RoleName string   `json:"role"`
	ObjectId string   `json:"object"`
	Scopes   []string `json:"scopes"`
}

type privateKey int

type Auth_Op int

const (
	GET_TOKEN Auth_Op = iota
	VALIDATE_TOKEN
	REFRESH_TOKEN
	REVOKE_TOKEN
	USER_INFO
	APP_ROLE_SCOPES
	APP_ROLES
	APP_SCOPES
)

// RBACImplMock mock struct for unit testing
type RBACImplMock struct {
	MockNewTokenFromRequestHeader func(ctx context.Context, header http.Header) (context.Context, error)
	MockNewTokenFromClientId      func(ctx context.Context, clientId, clientSecret string, scopes []string) (context.Context, error)
	MockNewContextWithToken       func(ctx context.Context, u Token) context.Context
	MockTokenFromContext          func(ctx context.Context) (Token, bool)
	MockTokenOp                   func(ctx context.Context, op Auth_Op) (context.Context, error)
	MockValidateTokenOp           func(ctx context.Context) (bool, error)
	MockValidateScopeOp           func(ctx context.Context, scope string) (bool, error)
}

// NewTokenFromRequestHeader
func (r RBACImplMock) NewTokenFromRequestHeader(ctx context.Context, header http.Header) (context.Context, error) {
	return r.MockNewTokenFromRequestHeader(ctx, header)
}

// NewTokenFromClientId
func (r RBACImplMock) NewTokenFromClientId(ctx context.Context, clientId, clientSecret string, scopes []string) (context.Context, error) {
	return r.MockNewTokenFromClientId(ctx, clientId, clientSecret, scopes)
}

// NewContextWithToken
func (r RBACImplMock) NewContextWithToken(ctx context.Context, u Token) context.Context {
	return r.MockNewContextWithToken(ctx, u)
}

// TokenFromContext
func (r RBACImplMock) TokenFromContext(ctx context.Context) (Token, bool) {
	return r.MockTokenFromContext(ctx)
}

// TokenOp
func (r RBACImplMock) TokenOp(ctx context.Context, op Auth_Op) (context.Context, error) {
	return r.MockTokenOp(ctx, op)
}

// ValidateTokenOp
func (r RBACImplMock) ValidateTokenOp(ctx context.Context) (bool, error) {
	return r.MockValidateTokenOp(ctx)
}

// ValidateScopeOp
func (r RBACImplMock) ValidateScopeOp(ctx context.Context, scope string) (bool, error) {
	return r.MockValidateScopeOp(ctx, scope)
}

// for testing purposes
type RBACMockTest1 struct{}

func (r RBACMockTest1) NewTokenFromRequestHeader(ctx context.Context, header http.Header) (context.Context, error) {
	return context.WithValue(ctx, tKey, &MockToken1{
		AccessToken: "sdfsd",
	}), nil
}
func (r RBACMockTest1) NewTokenFromClientId(ctx context.Context, clientId, clientSecret string, scopes []string) (context.Context, error) {
	return context.WithValue(ctx, tKey, &MockToken1{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}), nil
}
func (r RBACMockTest1) NewContextWithToken(ctx context.Context, u Token) context.Context {
	return context.WithValue(ctx, tKey, u)
}
func (r RBACMockTest1) TokenFromContext(ctx context.Context) (Token, bool) {
	t, ok := ctx.Value(tKey).(Token)
	return t, ok
}
func (r RBACMockTest1) TokenOp(ctx context.Context, op Auth_Op) (context.Context, error) {
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
func (r RBACMockTest1) ValidateTokenOp(ctx context.Context) (bool, error) {
	t, _ := r.TokenFromContext(ctx)
	return t.ValidateToken(ctx)
}
func (r RBACMockTest1) ValidateScopeOp(ctx context.Context, scope string) (bool, error) {
	t, _ := r.TokenFromContext(ctx)
	return t.ValidateScope(ctx, scope)
}

type MockToken1 struct {
	AccessToken  string `json:"access_token"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (u *MockToken1) RetrieveUser(ctx context.Context) (*UserInfo, error) {
	return nil, errors.New("")
}
func (t *MockToken1) RetrieveAccessToken() string {
	return t.AccessToken
}
func (t *MockToken1) GetAccessToken(ctx context.Context) error {
	return errors.New("error")
}

func (t *MockToken1) ValidateToken(ctx context.Context) (bool, error) {
	return false, errors.New("error")
}

func (t *MockToken1) ValidateScope(ctx context.Context, scope string) (bool, error) {
	return false, errors.New("error")
}

func (t *MockToken1) RefreshAccessToken(ctx context.Context) error {
	return errors.New("error")
}

func (t *MockToken1) RevokeToken(ctx context.Context) error {
	return errors.New("error")
}

func (t *MockToken1) GetUserInfo(ctx context.Context) error {
	return errors.New("error")
}

func (t *MockToken1) GetRolesApps(ctx context.Context) error {
	return errors.New("error")
}

func (t *MockToken1) GetAppScopes(ctx context.Context) error {
	return errors.New("error")
}

func (t *MockToken1) GetAppRoleScopes(ctx context.Context) error {
	return errors.New("error")
}

type RBACMockTest2 struct{}

func (r RBACMockTest2) NewTokenFromRequestHeader(ctx context.Context, header http.Header) (context.Context, error) {
	return context.WithValue(ctx, tKey, &MockToken2{
		AccessToken: "sdfsd",
	}), nil
}
func (r RBACMockTest2) NewTokenFromClientId(ctx context.Context, clientId, clientSecret string, scopes []string) (context.Context, error) {
	return context.WithValue(ctx, tKey, &MockToken2{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}), nil
}
func (r RBACMockTest2) NewContextWithToken(ctx context.Context, u Token) context.Context {
	return context.WithValue(ctx, tKey, u)
}
func (r RBACMockTest2) TokenFromContext(ctx context.Context) (Token, bool) {
	t, ok := ctx.Value(tKey).(Token)
	return t, ok
}
func (r RBACMockTest2) TokenOp(ctx context.Context, op Auth_Op) (context.Context, error) {
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
func (r RBACMockTest2) ValidateTokenOp(ctx context.Context) (bool, error) {
	t, _ := r.TokenFromContext(ctx)
	return t.ValidateToken(ctx)
}
func (r RBACMockTest2) ValidateScopeOp(ctx context.Context, scope string) (bool, error) {
	t, _ := r.TokenFromContext(ctx)
	return t.ValidateScope(ctx, scope)
}

type MockToken2 struct {
	AccessToken  string `json:"access_token"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (u *MockToken2) RetrieveUser(ctx context.Context) (*UserInfo, error) {
	return &UserInfo{}, nil
}
func (t *MockToken2) RetrieveAccessToken() string {
	return t.AccessToken
}
func (t *MockToken2) GetAccessToken(ctx context.Context) error {
	return nil
}

func (t *MockToken2) ValidateToken(ctx context.Context) (bool, error) {
	return true, nil
}

func (t *MockToken2) ValidateScope(ctx context.Context, scope string) (bool, error) {
	return true, nil
}

func (t *MockToken2) RefreshAccessToken(ctx context.Context) error {
	return nil
}

func (t *MockToken2) RevokeToken(ctx context.Context) error {
	return nil
}

func (t *MockToken2) GetUserInfo(ctx context.Context) error {
	return nil
}

func (t *MockToken2) GetRolesApps(ctx context.Context) error {
	return nil
}

func (t *MockToken2) GetAppScopes(ctx context.Context) error {
	return nil
}

func (t *MockToken2) GetAppRoleScopes(ctx context.Context) error {
	return nil
}
