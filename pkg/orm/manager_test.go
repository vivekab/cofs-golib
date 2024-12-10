package goliborm

import (
	"context"
	"os"
	"testing"

	goliblocalise "github.com/vivekab/golib/pkg/localise"

	"github.com/stretchr/testify/assert"
	config "github.com/vivekab/golib/pkg/config"
	goliberror "github.com/vivekab/golib/pkg/error"
	golibtypes "github.com/vivekab/golib/pkg/types"
	"github.com/vivekab/golib/protobuf/protoroot"
)

func setup() DBManager {
	os.Setenv("RDS_QLDB_CONSUMER_WRITE_PASSWORD", "test")
	config.SetupTestConfig("API_ENV", "RDS_QLDB_CONSUMER_WRITE_PASSWORD")
	goliblocalise.InitializeTest()
	return NewDBManager()
}

func TestAddConnection(t *testing.T) {
	mgr := setup()
	err := mgr.AddConnection(context.TODO(), &ConnectionConfig{
		ServiceName:    golibtypes.ServiceNameAuth,
		ConnectionName: "test_write",
		DbSource:       protoroot.DbSource_DB_SOURCE_API,
		DbMode:         protoroot.DbMode_DB_MODE_WRITE,
		IsDefault:      false,
	})

	if err != nil {
		assert.ErrorContains(t, err, "no such host")
	} else {
		t.Log("connected")
	}
}

func TestGetConnection(t *testing.T) {
	mgr := setup()
	_, err := mgr.GetConnection(context.TODO(), "test_write")
	if err != nil {
		assert.Equal(t, goliberror.GetErrCode(err).String(), protoroot.ErrorCode_ERROR_CODE_RESOURCE_NOT_FOUND.String())
	} else {
		t.Log("retrieved connection")
	}
}
