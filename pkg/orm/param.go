package goliborm

import (
	config "github.com/vivekab/golib/pkg/config"
	golibtypes "github.com/vivekab/golib/pkg/types"
	"github.com/vivekab/golib/protobuf/protoroot"
)

func getConnectionParamsBySource(source protoroot.DbSource, mode protoroot.DbMode) ConnectionParams {
	params := ConnectionParams{
		Env:        config.GetString("API_ENV"),
		ConnectEnv: config.GetString("CONNECT_ENV"),
	}
	switch source {
	case protoroot.DbSource_DB_SOURCE_API:
		if mode == protoroot.DbMode_DB_MODE_WRITE {
			params.Host = config.GetString("coreWriteDB.host")
			params.DbName = config.GetString("coreWriteDB.dbname")
			params.Port = config.GetInt("coreWriteDB.port")
			params.SslMode = config.GetString("coreWriteDB.sslmode")
			params.MaxIdleConnections = config.GetInt("coreWriteDB.maxIdleConnections")
			params.MaxOpenConnections = config.GetInt("coreWriteDB.maxOpenConnections")
			params.MaxConnectionLifetime = config.GetInt("coreWriteDB.maxConnectionLifetime")
		} else {
			params.Host = config.GetString("apiReadDB.host")
			params.DbName = config.GetString("apiReadDB.dbname")
			params.Port = config.GetInt("apiReadDB.port")
			params.SslMode = config.GetString("apiReadDB.sslmode")
			params.MaxIdleConnections = config.GetInt("apiReadDB.maxIdleConnections")
			params.MaxOpenConnections = config.GetInt("apiReadDB.maxOpenConnections")
			params.MaxConnectionLifetime = config.GetInt("apiReadDB.maxConnectionLifetime")
		}
	case protoroot.DbSource_DB_SOURCE_API_LIST:
		params.Host = config.GetString("apiListReadDB.host")
		params.DbName = config.GetString("apiListReadDB.dbname")
		params.Port = config.GetInt("apiListReadDB.port")
		params.SslMode = config.GetString("apiListReadDB.sslmode")
		params.MaxIdleConnections = config.GetInt("apiListReadDB.maxIdleConnections")
		params.MaxOpenConnections = config.GetInt("apiListReadDB.maxOpenConnections")
		params.MaxConnectionLifetime = config.GetInt("apiListReadDB.maxConnectionLifetime")
	}
	return params
}

func getConnectionParams(service golibtypes.ServiceName, dbSource protoroot.DbSource, mode protoroot.DbMode) ConnectionParams {
	connectionParams := getConnectionParamsBySource(dbSource, mode)
	switch service {
	case golibtypes.ServiceNameAuth:
		if mode == protoroot.DbMode_DB_MODE_WRITE {
			connectionParams.User = "auth_user"
			connectionParams.Password = config.GetString("RDS_AUTH_WRITE_PASSWORD")
		} else {
			connectionParams.User = "auth_ro_user"
			connectionParams.Password = config.GetString("RDS_AUTH_READ_PASSWORD")
		}

	default:
		if mode == protoroot.DbMode_DB_MODE_WRITE {
			connectionParams.User = "internal_user"
			connectionParams.Password = config.GetString("RDS_INTERNAL_WRITE_PASSWORD")
		} else {
			connectionParams.User = "internal_ro_user"
			connectionParams.Password = config.GetString("RDS_INTERNAL_READ_PASSWORD")
		}
	}
	return connectionParams
}
