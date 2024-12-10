package goliblocalise

import (
	"context"
	"testing"

	"github.com/vivekab/golib/protobuf/protoroot"
)

func TestLocalise_GetAttributes(t *testing.T) {
	ctx := context.Background()
	InitializeTest()
	Start()
	attr := GetErrorAtrributes(ctx, protoroot.ErrorCode_ERROR_CODE_AUTH_FORBIDDEN)
	if attr.HttpCode != 401 {
		t.Log("error get localise attributes")
	}
}
