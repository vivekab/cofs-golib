package golibgrpc

import (
	"context"
	"fmt"
	"strings"

	golibconstants "github.com/vivekab/golib/pkg/constants"
	golibtypes "github.com/vivekab/golib/pkg/types"
)

func isServiceRunningLocally(srvName golibtypes.ServiceName, localServices string) bool {
	for _, srv := range strings.Split(localServices, ",") {
		if strings.TrimSpace(srv) == string(srvName) {
			return true
		}
	}
	return false
}
func getConnectionString(serviceName golibtypes.ServiceName, env string, port string) string {
	var r string
	baseUrl := golibtypes.GetBaseUrl(serviceName)
	switch env {
	case golibconstants.EnvLocal:
		r = fmt.Sprintf("%s.dev.svc.cluster.local:%s", baseUrl, port)
	case golibconstants.EnvDev:
		r = fmt.Sprintf("%s.dev.svc.cluster.local:%s", baseUrl, port)
	case golibconstants.EnvQA:
		r = fmt.Sprintf("%s.qa.svc.cluster.local:%s", baseUrl, port)
	}
	return r
}

func GetConnectionStringForService(ctx context.Context, options ClientOptions) (string, error) {
	var (
		connStr string
		err     error
	)
	// TODO Validate options
	connStr = getConnectionString(options.Name, options.Env, options.Port)
	if options.ConnectEnv != "" {
		if options.Env == golibconstants.EnvLocal && !strings.Contains(connStr, "localhost") {
			testStr := strings.TrimSuffix(connStr, ".dev:3001")
			switch strings.ToLower(options.ConnectEnv) {
			case golibconstants.EnvDev:
				connStr = testStr + ".dev:3001"
			case golibconstants.EnvQA:
				connStr = testStr + ".qa:3001"
			}
		}
		if isServiceRunningLocally(options.Name, options.LocalRunningServices) {
			connStr = "localhost:" + localServicePortMap[options.Name]
		}
	}
	if connStr != "" {
		return connStr, err
	}
	return "", fmt.Errorf("service %s must be registered", options.Name)
}
