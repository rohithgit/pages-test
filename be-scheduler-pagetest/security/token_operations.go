package security

import (
	"errors"
	"net/http"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/rbac"

	"golang.org/x/net/context"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"strconv"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"strings"
	"github.com/asaskevich/govalidator"
)

type rbacModelKey string

var (
	CtxRbacModelKey rbacModelKey = "rbacmodel"
	CtxRbacModelV2Key rbacModelKey = "rbacmodelv2"
)

type Validator struct {
	Handler http.Handler
}

func (validator Validator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Println( "Inside Serve HTTP")
	if !(r.URL.EscapedPath() == strings.ToLower(constants.PATH_PREFIX + constants.HEALTH_CHECK)) {
		//if !ValidateRequest(w, r, constants.SCOPE_READ) {
		//	utils.PrintError(w, 400, "Invalid Authorization.")
		//	return
		//}

		//utils.SpectreLog.Debugln("Token is Valid")

		//input validation
		appId := r.URL.Query().Get("applicationId")
		if appId != "" && !(govalidator.IsAlphanumeric(appId) || govalidator.IsUUID(appId)) {
			utils.PrintError(w, 400, "Invalid application ID format.")
			return
		}

		name := r.URL.Query().Get("name")
		if name != "" && (govalidator.IsJSON(name) || govalidator.Contains(name, "<") || govalidator.Contains(name, ">")) {
			utils.PrintError(w, 400, "Invalid name.")
			return
		}
		url := r.URL.Query().Get("url")
		if url != "" && (!govalidator.IsURL(url) || govalidator.Contains(url, "<") || govalidator.Contains(url, ">")) {
			utils.PrintError(w, 400, "Invalid url.")
			return
		}
		lookupurl := r.URL.Query().Get("lookupurl")
		if lookupurl != "" && (!govalidator.IsURL(lookupurl) || govalidator.Contains(lookupurl, "<") || govalidator.Contains(lookupurl, ">")) {
			utils.PrintError(w, 400, "Invalid lookupurl.")
			return
		}
		numRows := r.URL.Query().Get("numRows")
		if numRows != "" && !govalidator.IsNumeric(numRows) {
			utils.PrintError(w, 400, "Invalid numRows attribute.")
			return
		}
		startPage := r.URL.Query().Get("startPage")
		if startPage != "" && !govalidator.IsNumeric(startPage) {
			utils.PrintError(w, 400, "Invalid startPage attribute.")
			return
		}
	}
	validator.Handler.ServeHTTP(w, r)
}

// RbacModelFromContext returns the RBAC value stored in ctx, if any.
func RbacModelFromContext(ctx context.Context) (rbac.RBAC, bool) {
	rbacModel, ok := ctx.Value(CtxRbacModelKey).(rbac.RBAC)
	return rbacModel, ok
}

func RbacModelV2FromContext(ctx context.Context) (rbac.RBACV2, bool) {
	rbacModel, ok := ctx.Value(CtxRbacModelV2Key).(rbac.RBACV2)
	return rbacModel, ok
}

// ValidateAccessToken from header to validate authorization to perform CRUD operations
func ValidateAccessToken(ctx context.Context, header http.Header) (context.Context, error) {
	var (
		err          error
		isValidToken bool
	)

	rbacModelObj, ok := RbacModelFromContext(ctx)

	if !ok {
		utils.SpectreLog.Error("RbacModel not found in context")
		return ctx, errors.New("RbacModel not found in context")
	}

	//validate Access token before invoking any CRUD operations
	if ctx, err = rbacModelObj.NewTokenFromRequestHeader(ctx, header); err != nil {
		utils.SpectreLog.Errorf("Failed to get the Access token from header: %s", err.Error())
		return ctx, err
	}

	if isValidToken, err = rbacModelObj.ValidateTokenOp(ctx); err != nil {
		utils.SpectreLog.Errorf("Validate TokenOp errored -> %v", err)
		return ctx, err
	}

	if !isValidToken {
		utils.SpectreLog.Errorf("Token is not valid in context")
		return ctx, err
	}
	return ctx, err
}

// ValidateAccessTokenScope to valudate access token scope for CRUD operations
// ValidateAccessTokenScope to valudate access token scope for CRUD operations
func ValidateTokenScope(ctx context.Context, scope string) (context.Context, error) {

	var (
		err          error
		isValidScope bool
	)

	rbacModelObj, ok := RbacModelFromContext(ctx)

	if !ok {
		utils.SpectreLog.Error("RbacModel not found in context")
		return ctx, errors.New("RbacModel not found in context")
	}

	//validate scope of the token for CRUD operations
	if isValidScope, err = rbacModelObj.ValidateScopeOp(ctx, scope); err != nil {
		utils.SpectreLog.Errorf("scope op validate errored -> %s", err.Error())
		return ctx, err
	}

	if !isValidScope {
		utils.SpectreLog.Errorf("Scope is not valid for given ctx")
		return ctx, errors.New("Invalid scope.")
	}

	utils.SpectreLog.Debugln("isValidScope :", isValidScope)
	return ctx, nil
}

func ValidateAppScope(ctx context.Context, appId string, scope string) (context.Context, error) {

	var (
		err          error
		isValidScope bool
	)

	rbacModelObj, ok := RbacModelV2FromContext(ctx)

	if !ok {
		utils.SpectreLog.Error("RbacModel not found in context")
		return ctx, errors.New("RbacModel not found in context")
	}
	if appId == ""  {

		utils.SpectreLog.Error("application is not provided")
		return ctx, errors.New("Application ID not provided")
	}

	utils.SpectreLog.Debugf("Validating appId %s and scope %s\n", appId, scope)

	//validate scope of the token for CRUD operations
	if isValidScope, err = rbacModelObj.ValidateAppScopeOp(ctx, appId, scope); err != nil {
		utils.SpectreLog.Errorf("scope op validate errored -> %s", err.Error())
		return ctx, err
	}

	if !isValidScope {
		utils.SpectreLog.Errorf("Scope is not valid for given ctx")
		return ctx, errors.New("Invalid scope.")
	}

	utils.SpectreLog.Debugln("isValidScope :", isValidScope)
	return ctx, nil
}


// GetAccessToken API to get the access token from RBAC service using clientID and clientSecret
func GetAccessToken(ctx context.Context, scope []string) (*rbac.APIToken, error) {
	var (
		err error
	)

	rbacModelObj, ok := RbacModelFromContext(ctx)

	if !ok {
		utils.SpectreLog.Error("RbacModel not found in context")
		return nil, errors.New("RbacModel not found in context")
	}

	ctx, err = rbacModelObj.NewTokenFromClientId(ctx, global.Options.ClientId, global.Options.ClientSecret, scope)
	if err != nil {
		utils.SpectreLog.Errorf("Failed to get new token from clientId %v", err)
		return nil, err
	}

	ctx, err = rbacModelObj.TokenOp(ctx, rbac.GET_TOKEN)
	if err != nil {
		utils.SpectreLog.Errorf("topkenop error -> %v", err)
		return nil, err
	}

	token, ok := rbacModelObj.TokenFromContext(ctx)
	if !ok {
		utils.SpectreLog.Errorln("Failed to get token from context ")
		return nil, errors.New("get access token error")
	}

	apiToken, ok := token.(*rbac.APIToken)
	if !ok {
		utils.SpectreLog.Errorf("type assertion errored")
		return nil, errors.New("type assertion error")
	}

	utils.SpectreLog.Debugln("GetAccesstoken completed")
	return apiToken, err
}

func ValidateRequest(w http.ResponseWriter, r *http.Request, scopeOp string) bool {
	var (
		err error
		ctx context.Context
	)

	scope := constants.SCOPE_SUFFIX
	//TODO: remove condition once RBAC is implemented all the UI services and token is available in header
	//Check env option to validate the token
	if !global.Options.ValidateToken {
		utils.SpectreLog.Debugln("RBAC is turned off. Set env. varaible RBAC_ON=TRUE to turn it on.")
		return true
	}

	scope += scopeOp

	appId := r.URL.Query().Get("applicationId")
	rbacModelv2 := new(rbac.RBACImplV2)
	ctx = context.WithValue(context.TODO(), CtxRbacModelKey, rbacModelv2)

	ctx, err = ValidateAccessToken(ctx, r.Header)
	if err != nil {
		utils.SpectreLog.Errorf("Token is not valid -> %s", err.Error())
		if global.Options.ValidateToken {
			w.Header().Add("Status", strconv.Itoa(http.StatusForbidden) + " " + constants.CODE_403_MESSAGE)
			w.Header().Add("content-type", constants.JSON_CONTENT_TYPE)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			return false
		}
	}

	utils.SpectreLog.Info("Token validated successfully...")

	ctx = context.WithValue(ctx, CtxRbacModelV2Key, rbacModelv2 )
	if err == nil { // check only if there is no error in ValidateAccessToken
		utils.SpectreLog.Info("Validating scope..")
		ctx, err = ValidateAppScope(ctx, appId, scope)

		if err != nil {
			utils.SpectreLog.Errorln("ValidateAppScope is errored %v", err)
			if global.Options.ValidateToken {
				w.Header().Add("Status", strconv.Itoa(http.StatusInternalServerError)+" "+constants.CODE_500_MESSAGE)
				w.Header().Add("content-type", constants.JSON_CONTENT_TYPE)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server error"))
				return false
			}

		}
	}
	utils.SpectreLog.Infof("Scope validated successfully %s", scopeOp)
	return true
}

func GetUserId(r *http.Request) string {
	var (
		err      error
		ctx      context.Context
		token    rbac.Token
		ok       bool
		userInfo *rbac.UserInfo
	)

	rbacModel := new(rbac.RBACImpl)
	ctx = context.WithValue(context.TODO(), CtxRbacModelKey, rbacModel)

	//retrieve Access token before invoking any CRUD operations
	if ctx, err = rbacModel.NewTokenFromRequestHeader(ctx, r.Header); err != nil {
		utils.SpectreLog.Errorf("Failed to get the Access token from header: %s", err.Error())
		return ""
	}

	if token, ok = rbacModel.TokenFromContext(ctx); !ok {
		utils.SpectreLog.Errorf("Failed to get Token obj %v", err)
		return ""
	}

	userInfo, err = token.RetrieveUser(ctx)

	if err != nil {
		utils.SpectreLog.Errorf("Failed to get user info %v", err)
		return ""
	}

	return userInfo.Username
}
