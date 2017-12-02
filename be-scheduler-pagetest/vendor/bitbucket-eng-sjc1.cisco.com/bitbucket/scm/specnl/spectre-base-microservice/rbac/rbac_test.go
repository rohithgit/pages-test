package rbac

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNewTokenFromRequestHeader(t *testing.T) {
	r := RBACImpl{}
	_, err := r.NewTokenFromRequestHeader(context.TODO(), http.Header{})
	assert.NotNil(t, err)

	ctx, err := r.NewTokenFromRequestHeader(context.TODO(), http.Header{"Authorization": []string{"Bearer asdf"}})
	assert.Nil(t, err)
	assert.NotNil(t, ctx)
}

func TestNewTokenFromClientId(t *testing.T) {
	r := RBACImpl{}
	_, err := r.NewTokenFromClientId(context.TODO(), "", "sdf", nil)
	assert.NotNil(t, err)

	_, err = r.NewTokenFromClientId(context.TODO(), "asdf", "", nil)
	assert.NotNil(t, err)

	_, err = r.NewTokenFromClientId(context.TODO(), "asdf", "Asdf", nil)
	assert.Nil(t, err)
}
