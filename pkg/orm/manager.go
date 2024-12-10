package goliborm

import (
	"context"
	"sync"

	"gorm.io/gorm"

	golibconstants "github.com/vivekab/golib/pkg/constants"
	golibcontext "github.com/vivekab/golib/pkg/context"
	goliberror "github.com/vivekab/golib/pkg/error"
	golibtypes "github.com/vivekab/golib/pkg/types"
	"github.com/vivekab/golib/protobuf/protoroot"
)

type DBManager interface {
	AddConnection(ctx context.Context, config *ConnectionConfig) error
	GetConnection(ctx context.Context, connectionName string) (*Orm, error)
	GetConnectionBySourceMode(ctx context.Context, dbSource protoroot.DbSource, mode protoroot.DbMode) (*Orm, error)
	Scan(ctx context.Context, dest interface{}, query string, args ...interface{}) *Result
	ScanWithTransaction(ctx context.Context, tx *gorm.DB, dest interface{}, query string, args ...interface{}) *Result
	Close(ctx context.Context) error
}

type manager struct {
	connections       map[string]*Orm
	srcMode           map[string]string
	defaultConnection *Orm
	mu                sync.Mutex
}

func (mgr *manager) Close(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewDBManager() DBManager {
	return &manager{
		connections: make(map[string]*Orm),
		srcMode:     make(map[string]string),
	}
}

type ConnectionConfig struct {
	// defines the service which is connecting to the db
	ServiceName golibtypes.ServiceName
	// defines the connection name
	ConnectionName string
	// define the db source which caller is trying to connect
	DbSource protoroot.DbSource
	// defines the mode in which connection is to be made ( read or write )
	DbMode protoroot.DbMode
	// defines whether the connection being added is default or not
	IsDefault bool
}

// Adds connection to the manager
func (mgr *manager) AddConnection(ctx context.Context, config *ConnectionConfig) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if config.IsDefault && mgr.defaultConnection != nil {
		return goliberror.NewError(ctx,
			&protoroot.GrpcError{Code: protoroot.ErrorCode_ERROR_CODE_DB_ERROR, SysMessage: "default is already set"})
	}

	db, err := Open(config.ServiceName, config.DbSource, config.DbMode)
	if err != nil {
		return err
	}
	srcModeKey := config.DbSource.String() + ":" + config.DbMode.String()
	mgr.srcMode[srcModeKey] = config.ConnectionName
	mgr.connections[config.ConnectionName] = db
	if config.IsDefault {
		mgr.defaultConnection = db
	}
	return nil
}

// Get the connection from manager based on name
func (mgr *manager) GetConnection(ctx context.Context, connectionName string) (*Orm, error) {
	db, ok := mgr.connections[connectionName]
	if !ok {
		return nil, goliberror.NewError(ctx, &protoroot.GrpcError{
			Code:    protoroot.ErrorCode_ERROR_CODE_RESOURCE_NOT_FOUND,
			Message: "db connection not found",
		})
	}
	return db, nil
}

func (mgr *manager) GetConnectionBySourceMode(ctx context.Context, dbSource protoroot.DbSource, mode protoroot.DbMode) (*Orm, error) {
	srcModeKey := dbSource.String() + ":" + mode.String()
	connectionName, ok := mgr.srcMode[srcModeKey]
	if !ok {
		return nil, goliberror.NewError(ctx, &protoroot.GrpcError{
			Code:    protoroot.ErrorCode_ERROR_CODE_RESOURCE_NOT_FOUND,
			Message: "db connection not found",
		})
	}
	return mgr.GetConnection(ctx, connectionName)
}

func (mgr *manager) getConnectionFromContext(ctx context.Context) (*Orm, error) {
	dbSourceModeKey := golibcontext.GetFromContext(ctx, golibconstants.HeaderDBSource)
	if connName, ok := mgr.srcMode[dbSourceModeKey]; ok {
		return mgr.GetConnection(ctx, connName)
	} else {
		if mgr.defaultConnection != nil {
			return mgr.defaultConnection, nil
		}
	}
	return nil, goliberror.NewError(ctx, &protoroot.GrpcError{
		Code:    protoroot.ErrorCode_ERROR_CODE_RESOURCE_NOT_FOUND,
		Message: "db connection not found",
	})
}

func (mgr *manager) getResultError(ctx context.Context, err error) *Result {
	return &Result{Error: err}
}

// Run the select query
func (mgr *manager) Scan(ctx context.Context, dest interface{}, query string, args ...interface{}) *Result {
	db, err := mgr.getConnectionFromContext(ctx)
	if err != nil {
		return mgr.getResultError(ctx, err)
	}
	return db.Raw(ctx, query, args...).Scan(ctx, dest)
}

func (mgr *manager) ScanWithTransaction(ctx context.Context, tx *gorm.DB, dest interface{}, query string, args ...interface{}) *Result {
	db, err := mgr.getConnectionFromContext(ctx)
	if err != nil {
		return mgr.getResultError(ctx, err)
	}
	txn := getDbWithCtx(ctx, db.db, tx)

	return getResult(txn.Raw(query, args...).Scan(dest))
}
