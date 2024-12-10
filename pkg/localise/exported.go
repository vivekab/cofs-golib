package goliblocalise

import (
	"context"
	"net/http"

	golibconfig "github.com/vivekab/golib/pkg/config"
	golibtypes "github.com/vivekab/golib/pkg/types"

	goliblogging "github.com/vivekab/golib/pkg/logging"
	"github.com/vivekab/golib/protobuf/protoroot"
	grpcRoot "github.com/vivekab/golib/protobuf/protoroot"
)

var (
	store *localizationStore
)

const (
	refreshInterval = 600
	repoUrl         = "https://raw.githubusercontent.com/solidfi/errors/master"
)

func InitializeTest() {
	golibconfig.SetupTestConfig("API_ENV")
	Initialize(context.Background(), golibtypes.ServiceNameAuth)
	Start()
}

func Initialize(ctx context.Context, serviceName golibtypes.ServiceName) {
	client := NewLocaliseService(ctx, repoUrl)
	store = newStore(refreshInterval, serviceName, client)
}

func Start() {
	store.Start()
}

func Stop() {
	store.Stop()
}

func GetErrorAtrributes(ctx context.Context, errCode protoroot.ErrorCode) *LocaliseAttributes {
	if store == nil {
		goliblogging.PanicD(ctx, "error getting localise Translation Struct since the store is not initialised", goliblogging.Fields{
			"err":     "store not initialised",
			"errCode": errCode.String(),
		})
	}
	key := getLocaliseStoreMapKey(grpcRoot.Language_LANGUAGE_EN, errCode.String())
	if value, ok := store.store[key]; ok {
		return &value
	}
	goliblogging.WarnD(ctx, "errLocaliseErrorNotFound", goliblogging.Fields{
		"errCode": errCode.String(),
	})
	return &LocaliseAttributes{Value: "", HttpCode: http.StatusInternalServerError}
}
