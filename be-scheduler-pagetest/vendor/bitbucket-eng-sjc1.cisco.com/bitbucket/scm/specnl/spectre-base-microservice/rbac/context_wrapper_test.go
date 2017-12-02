package rbac

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNewContextWithToken(t *testing.T) {
	r := RBACImpl{}
	ctx := context.TODO()
	assert.NotEqual(t, r.NewContextWithToken(ctx, &MockToken1{}), ctx)

	// This is a repetition with Mock
	r1 := RBACMockTest1{}
	ctx = context.TODO()
	assert.NotEqual(t, r1.NewContextWithToken(ctx, &MockToken1{}), ctx)
}

func TestTokenFromContext(t *testing.T) {
	r := RBACImpl{}
	v, ok := r.TokenFromContext(context.TODO())
	assert.False(t, ok)
	assert.Nil(t, v)

	v, ok = r.TokenFromContext(r.NewContextWithToken(context.TODO(), &MockToken1{}))
	assert.True(t, ok)
	assert.NotNil(t, v)

	// This is a repetition with Mock
	r1 := RBACMockTest1{}
	v, ok = r1.TokenFromContext(context.TODO())
	assert.False(t, ok)
	assert.Nil(t, v)

	v, ok = r.TokenFromContext(r.NewContextWithToken(context.TODO(), &MockToken1{}))
	assert.True(t, ok)
	assert.NotNil(t, v)
}

func TestTokenOp(t *testing.T) {
	r := RBACImpl{}
	ops := []Auth_Op{GET_TOKEN, REFRESH_TOKEN, REVOKE_TOKEN, USER_INFO, APP_ROLES, APP_SCOPES, APP_ROLE_SCOPES}

	tokens := []Token{&MockToken1{}, &MockToken1{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		ctx := r.NewContextWithToken(context.TODO(), to)
		for _, op := range ops {
			ctx1, err := r.TokenOp(ctx, op)
			assert.NotNil(t, err)
			assert.Equal(t, ctx, ctx1)
		}
	}

	tokens = []Token{&MockToken2{}, &MockToken2{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		ctx := r.NewContextWithToken(context.TODO(), to)
		for _, op := range ops {
			ctx1, err := r.TokenOp(ctx, op)
			assert.Nil(t, err)
			assert.NotEqual(t, ctx, ctx1)
		}
	}

	// This is a repetition with Mock
	r1 := RBACMockTest1{}

	tokens = []Token{&MockToken1{}, &MockToken1{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		ctx := r1.NewContextWithToken(context.TODO(), to)
		for _, op := range ops {
			ctx1, err := r1.TokenOp(ctx, op)
			assert.NotNil(t, err)
			assert.Equal(t, ctx, ctx1)
		}
	}

	r2 := RBACMockTest2{}
	tokens = []Token{&MockToken2{}, &MockToken2{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		ctx := r2.NewContextWithToken(context.TODO(), to)
		for _, op := range ops {
			ctx1, err := r2.TokenOp(ctx, op)
			assert.Nil(t, err)
			assert.NotEqual(t, ctx, ctx1)
		}
	}
}

func TestValidateToken(t *testing.T) {
	r := RBACImpl{}
	tokens := []Token{&MockToken1{}, &MockToken1{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		ctx := r.NewContextWithToken(context.TODO(), to)
		v, err := r.ValidateTokenOp(ctx)
		assert.NotNil(t, err)
		assert.False(t, v)

		v, err = r.ValidateScopeOp(ctx, "t1")
		assert.NotNil(t, err)
		assert.False(t, v)
	}

	tokens = []Token{&MockToken2{}, &MockToken2{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		ctx := r.NewContextWithToken(context.TODO(), to)
		v, err := r.ValidateTokenOp(ctx)
		assert.Nil(t, err)
		assert.True(t, v)

		v, err = r.ValidateScopeOp(ctx, "t1")
		assert.Nil(t, err)
		assert.True(t, v)
	}


	// this is just a repetition with Mock
	r1 := RBACMockTest1{}
	tokens = []Token{&MockToken1{}, &MockToken1{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		ctx := r1.NewContextWithToken(context.TODO(), to)
		v, err := r1.ValidateTokenOp(ctx)
		assert.NotNil(t, err)
		assert.False(t, v)

		v, err = r1.ValidateScopeOp(ctx, "t1")
		assert.NotNil(t, err)
		assert.False(t, v)
	}

	r2 := RBACMockTest2{}
	tokens = []Token{&MockToken2{}, &MockToken2{
		AccessToken: "asdfds",
	}}
	for _, to := range tokens {
		ctx := r2.NewContextWithToken(context.TODO(), to)
		v, err := r2.ValidateTokenOp(ctx)
		assert.Nil(t, err)
		assert.True(t, v)

		v, err = r1.ValidateScopeOp(ctx, "t1")
		assert.Nil(t, err)
		assert.True(t, v)
	}
}
