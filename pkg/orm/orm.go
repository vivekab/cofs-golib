package goliborm

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	golibconfig "github.com/vivekab/golib/pkg/config"
	"gorm.io/gorm/logger"

	"gorm.io/gorm/clause"

	apmPostgres "go.elastic.co/apm/module/apmgormv2/driver/postgres"

	golibconstants "github.com/vivekab/golib/pkg/constants"
	golibtypes "github.com/vivekab/golib/pkg/types"
	golibutils "github.com/vivekab/golib/pkg/utils"
	"github.com/vivekab/golib/protobuf/protoroot"
	"gorm.io/gorm"
)

const (
	POSTGRES = "postgres"
	MYSQL    = "mysql"

	VerifyCa         = "verify-ca"
	VerifyFull       = "verify-full"
	SslRootCert      = "sslrootcert"
	GlobalBundlePath = "cert/psql-cert.pem"

	FieldNameVersion    = "Version"
	FieldNameDeletedAt  = "DeletedAt"
	FieldNameModifiedAt = "ModifiedAt"
)

type Orm struct {
	Error            error
	dialect          string
	connectionString string
	db               *gorm.DB
}

type ConnectionParams struct {
	Host                  string
	User                  string
	Password              string
	DbName                string
	Port                  int
	SslMode               string
	MaxOpenConnections    int
	MaxIdleConnections    int
	MaxConnectionLifetime int
	Env                   string
	ConnectEnv            string
}

func (c ConnectionParams) getConnectionString(dialect string) (string, error) {
	connStr := ""
	if dialect == POSTGRES {
		connStr = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			c.Host, c.User, c.Password, c.DbName, c.Port, c.SslMode)
		if c.SslMode == VerifyCa || c.SslMode == VerifyFull {
			connStr = connStr + fmt.Sprintf(" %s=%s", SslRootCert, GlobalBundlePath)
		}
	} else if dialect == MYSQL {
		connStr = fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
			c.User, c.Password, c.Host, c.Port, c.DbName)
	} else {
		return connStr, errors.New("not supported dialect")
	}
	return connStr, nil
}

type Result struct {
	Value        interface{}
	Error        error
	RowsAffected int64
}

func Open(service golibtypes.ServiceName, dbSource protoroot.DbSource, mode protoroot.DbMode) (*Orm, error) {
	connectionParams := getConnectionParams(service, dbSource, mode)
	dialect := "postgres"
	if connectionParams.Env == golibconstants.EnvLocal && connectionParams.ConnectEnv != "" {
		prefix := ""
		switch connectionParams.ConnectEnv {
		case golibconstants.EnvDev:
			prefix = "dev_"
		case golibconstants.EnvQA:
			prefix = "qa_"
		}
		if prefix != "" {
			nameSplit := strings.Split(connectionParams.DbName, "_")
			connectionParams.DbName = prefix + strings.Join(nameSplit[1:], "_")
			userSplit := strings.Split(connectionParams.User, "_")
			connectionParams.User = prefix + strings.Join(userSplit[1:], "_")
		}
	}
	dsn, err := connectionParams.getConnectionString(dialect)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(apmPostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqldb.SetMaxOpenConns(connectionParams.MaxOpenConnections)
	sqldb.SetMaxIdleConns(connectionParams.MaxIdleConnections)
	sqldb.SetConnMaxLifetime(time.Duration(connectionParams.MaxConnectionLifetime) * time.Second)

	o := &Orm{
		dialect:          dialect,
		connectionString: dsn,
		db:               db,
		Error:            db.Error,
	}

	env := golibconfig.GetString("API_ENV")
	if env != golibconstants.EnvProd {
		db.Logger = logger.Default.LogMode(logger.Info)
	}

	return o, nil
}

func (o *Orm) GetUnderLyingDb() *gorm.DB {
	return o.db
}

func (o *Orm) Close(ctx context.Context) error {
	sqlDb, err := o.db.DB()
	if err != nil {
		return err
	}
	return sqlDb.Close()
}

func (o *Orm) Where(query interface{}, args ...interface{}) *Orm {
	orm := o.clone()
	orm.db = o.db.Where(query, args...)
	orm.Error = orm.db.Error
	return orm
}

func (o *Orm) Or(query interface{}, args ...interface{}) *Orm {
	orm := o.clone()
	orm.db = o.db.Or(query, args...)
	orm.Error = orm.db.Error
	return orm
}

func (o *Orm) Not(query interface{}, args ...interface{}) *Orm {
	orm := o.clone()
	orm.db = o.db.Not(query, args...)
	orm.Error = orm.db.Error
	return orm
}

func (o *Orm) Select(query interface{}, args ...interface{}) *Orm {
	orm := o.clone()
	orm.db = o.db.Select(query, args...)
	orm.Error = orm.db.Error
	return orm
}

func (o *Orm) Omit(columns ...string) *Orm {
	orm := o.clone()
	orm.db = o.db.Omit(columns...)
	orm.Error = orm.db.Error
	return orm
}

func (o *Orm) Group(query string) *Orm {
	orm := o.clone()
	orm.db = o.db.Group(query)
	orm.Error = orm.db.Error
	return orm
}

func (o *Orm) Having(query interface{}, values ...interface{}) *Orm {
	orm := o.clone()
	orm.db = o.db.Having(query, values...)
	orm.Error = orm.db.Error
	return orm
}

func (o *Orm) Joins(query string, args ...interface{}) *Orm {
	orm := o.clone()
	orm.db = o.db.Joins(query, args...)
	orm.Error = orm.db.Error
	return orm
}

func (o *Orm) Assign(attrs ...interface{}) *Orm {
	orm := o.clone()
	orm.db = o.db.Assign(attrs...)
	orm.Error = orm.db.Error
	return orm
}

func (o *Orm) First(ctx context.Context, out interface{}, where ...interface{}) *Result {
	db := o.db.First(out, where...)
	return getResult(db)
}

func (o *Orm) Take(ctx context.Context, out interface{}, where ...interface{}) *Result {
	return getResult(o.db.Take(out, where...))
}

func (o *Orm) Last(ctx context.Context, out interface{}, where ...interface{}) *Result {
	return getResult(o.db.Last(out, where...))
}

// required
func (o *Orm) Find(ctx context.Context, out interface{}, where ...interface{}) *Result {
	return getResult(o.db.WithContext(ctx).Find(out, where...))
}

// required
func (o *Orm) Scan(ctx context.Context, dest interface{}) *Result {
	res := getResult(o.db.Scan(dest))
	if res.Error == nil && res.RowsAffected == 0 {
		res.Error = gorm.ErrRecordNotFound
	}
	return res
}

// require
func (o *Orm) Row(ctx context.Context) *sql.Row {
	result := o.db.Row()
	return result
}

func (o *Orm) Rows(ctx context.Context) (*sql.Rows, error) {
	result, err := o.db.WithContext(ctx).Rows()
	return result, err
}

func (o *Orm) ScanRows(ctx context.Context, rows *sql.Rows, result interface{}) error {
	return o.db.WithContext(ctx).ScanRows(rows, result)
}

func (o *Orm) Pluck(ctx context.Context, column string, value interface{}) *Result {
	return getResult(o.db.WithContext(ctx).Pluck(column, value))
}

func (o *Orm) Count(ctx context.Context, value *int64) *Result {
	var result *Result
	r := o.db.WithContext(ctx).Count(value)
	result = getResult(r)
	return result
}

func (o *Orm) FirstOrInit(ctx context.Context, out interface{}, where ...interface{}) *Result {
	return getResult(o.db.WithContext(ctx).FirstOrInit(out, where))
}

func (o *Orm) FirstOrCreate(ctx context.Context, out interface{}, where ...interface{}) *Result {
	return getResult(o.db.WithContext(ctx).FirstOrCreate(out, where))
}

// Updates takes updateByFieldValue to be passed in string seperated by comma
// for example
// updateByFieldValue should be -> 'act-1234','act-1235','act-1236'
func (o *Orm) Updates(ctx context.Context, model interface{}, updateByFieldName string, updateByFieldValue interface{}, toUpdate map[string]interface{}, tx *gorm.DB, version ...int) *Result {
	return o.update(ctx, updateFields{model, updateByFieldName, "%s in (?)", updateByFieldValue, toUpdate}, tx, version...)
}

func (o *Orm) Update(ctx context.Context, model interface{}, updateByFieldName string, updateByFieldValue interface{}, toUpdate map[string]interface{}, tx *gorm.DB, version ...int) *Result {
	return o.update(ctx, updateFields{model, updateByFieldName, "%s = ?", updateByFieldValue, toUpdate}, tx, version...)
}

type updateFields struct {
	model              interface{}
	updateByFieldName  string
	queryKey           string
	updateByFieldValue interface{}
	toUpdate           map[string]interface{}
}

func (o *Orm) update(ctx context.Context, u updateFields, tx *gorm.DB, version ...int) *Result {
	var (
		db                      = getDbWithCtx(ctx, o.db, tx)
		columnNameVersion       = golibutils.GetTableColumnFromModel(u.model, FieldNameVersion)
		columnNameModifiedAt    = golibutils.GetTableColumnFromModel(u.model, FieldNameModifiedAt)
		columnNameDeletedAt     = golibutils.GetTableColumnFromModel(u.model, FieldNameDeletedAt)
		columnNameUpdateByField = golibutils.GetTableColumnFromModel(u.model, u.updateByFieldName)
	)

	toUpdateWithColumnName := make(map[string]interface{})
	for key, value := range u.toUpdate {
		if golibutils.IsJsonBField(u.model, key) {
			value, _ = json.Marshal(value)
		}
		toUpdateWithColumnName[golibutils.GetTableColumnFromModel(u.model, key)] = value
	}

	// increase version on every update
	toUpdateWithColumnName[columnNameVersion] = gorm.Expr(fmt.Sprintf("%s + 1", columnNameVersion))
	toUpdateWithColumnName[columnNameModifiedAt] = time.Now().UTC()

	// if version is passed, update should only happen if version matches with DB value
	whereClause := fmt.Sprintf(u.queryKey+" and %s IS NULL", columnNameUpdateByField, columnNameDeletedAt)
	if len(version) == 1 {
		whereClause += fmt.Sprintf(" and %s = %d", columnNameVersion, version[0])
	}

	r := db.WithContext(ctx).Model(u.model).Clauses(clause.Returning{}).Where(whereClause, u.updateByFieldValue).Updates(toUpdateWithColumnName)
	return getResult(r)
}
func (o *Orm) UpdateColumns(ctx context.Context, values interface{}) *Result {
	return getResult(o.db.WithContext(ctx).UpdateColumns(values))
}

func (o *Orm) Save(ctx context.Context, value interface{}) *Result {
	return getResult(o.db.WithContext(ctx).Save(value))
}

func (o *Orm) Create(ctx context.Context, model interface{}, tx *gorm.DB) *Result {
	return getResult(getDbWithCtx(ctx, o.db, tx).Create(model))
}

func (o *Orm) Delete(ctx context.Context, model interface{}, deleteByFieldName string, deleteByFieldValue interface{}, tx *gorm.DB) *Result {
	var (
		db                      = getDbWithCtx(ctx, o.db, tx)
		timeNow                 = time.Now().UTC()
		columnNameDeletedAt     = golibutils.GetTableColumnFromModel(model, FieldNameDeletedAt)
		columnNameDeleteByField = golibutils.GetTableColumnFromModel(model, deleteByFieldName)
	)

	// set deleted_at timestamp
	toUpdate := map[string]interface{}{
		columnNameDeletedAt: &timeNow,
	}

	r := db.WithContext(ctx).Model(model).Where(fmt.Sprintf("%s = ? and %s IS NULL", columnNameDeleteByField, columnNameDeletedAt), deleteByFieldValue).Updates(toUpdate)
	return getResult(r)
}
func (o *Orm) Raw(ctx context.Context, sql string, values ...interface{}) *Orm {
	orm := o.clone()
	orm.db = o.db.WithContext(ctx).Raw(sql, values...)
	orm.Error = orm.db.Error
	return orm
}

func (o *Orm) Exec(ctx context.Context, sql string, values ...interface{}) *Result {
	return getResult(o.db.WithContext(ctx).Exec(sql, values...))
}

func SetTransactionTimeouts(ctx context.Context, tx *gorm.DB, statementTimeout, idleInTransactionSessionTimeout int64) error {
	if err := tx.Exec(fmt.Sprintf("SET LOCAL statement_timeout = %d", statementTimeout)).Error; err != nil {
		return err
	}
	if err := tx.Exec(fmt.Sprintf("SET LOCAL idle_in_transaction_session_timeout = %d", idleInTransactionSessionTimeout)).Error; err != nil {
		return err
	}
	return nil
}

func (o *Orm) Transaction(ctx context.Context, callback func(tx *gorm.DB) error) error {
	return o.db.Transaction(callback)
}

func (o *Orm) FetchRows(ctx context.Context, query string, output interface{}) *Result {
	r := o.db.WithContext(ctx).Raw(query).Scan(output)
	return getResult(r)
}

func (o *Orm) Begin(ctx context.Context) *Orm {
	orm := o.clone()
	orm.db = orm.db.Begin()
	return orm
}

func (o *Orm) Commit(ctx context.Context) *Result {
	return getResult(o.db.Commit())
}

func (o *Orm) Rollback(ctx context.Context) *Result {
	return getResult(o.db.Rollback())
}

func (o *Orm) AutoMigrate(values ...interface{}) error {
	return o.db.AutoMigrate(values...)
}

func (o *Orm) AddIndex(model interface{}, column string) error {
	err := o.db.Migrator().CreateIndex(&model, column)
	if err != nil {
		return err
	}
	return nil
}

func (o *Orm) Model(value interface{}) *Orm {
	orm := o.clone()
	orm.db = o.db.Model(value)
	return orm
}

func (o *Orm) Table(name string) *Orm {
	orm := o.clone()
	orm.db = o.db.Table(name)
	return orm
}

func (o *Orm) Association(column string) *gorm.Association {
	return o.db.Association(column)
}

func (o *Orm) clone() *Orm {
	return &Orm{
		dialect:          o.dialect,
		connectionString: o.connectionString,
		db:               o.db,
	}
}

func getResult(result *gorm.DB) *Result {
	var model interface{}
	if result.Statement != nil {
		model = result.Statement.Model
	}
	return &Result{
		Error:        result.Error,
		Value:        model,
		RowsAffected: result.RowsAffected,
	}
}

func getDbWithCtx(ctx context.Context, db, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return db
}
