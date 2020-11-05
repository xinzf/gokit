package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/xinzf/gokit/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dblog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type database struct {
	//options    DbOptions
	connections map[string]*gorm.DB
}

var DB *database

//var dbOnce sync.Once

func init() {
	DB = new(database)
	DB.connections = map[string]*gorm.DB{}
}

func (db *database) Register(opts ...DbOption) error {
	var err error

	ops := newDbOptions(opts...)
	if ops.Addr == "" {
		err = errors.New("dbconfig's addr is empty")
		return err
	}
	if ops.User == "" {
		err = errors.New("dbconfig's user is empty")
		return err
	}
	if ops.Name == "" {
		err = errors.New("dbconfig's name is empty")
		return err
	}
	if ops.Pswd == "" {
		err = errors.New("dbconfig's pswd is empty")
		return err
	}

	cfg := &gorm.Config{
		AllowGlobalUpdate: true,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableAutomaticPing: false,
	}
	if ops.Log {
		cfg.Logger = newDBLogger()
	}
	var d *gorm.DB
	d, err = gorm.Open(mysql.Open(ops.String()), cfg)
	if err != nil {
		return err
	}

	_db, _ := d.DB()
	_db.SetConnMaxLifetime(time.Duration(300) * time.Second)
	_db.SetMaxIdleConns(ops.MaxIdleConns)
	_db.SetMaxOpenConns(ops.MaxOpenConns)

	if len(ops.initSql) > 0 {
		for _, sql := range ops.initSql {
			if err = d.Exec(sql).Error; err != nil {
				return err
			}
		}
	}

	DB.connections[ops.Name] = d
	if cfg.Logger != nil {
		cfg.Logger.Info(context.TODO(), fmt.Sprintf("DB: %s init", ops.Name))
	}
	return nil
}

func (db *database) Use(dbName ...string) *gorm.DB {
	if len(db.connections) == 0 {
		panic("has none db connections")
	}
	if len(db.connections) > 1 && len(dbName) == 0 {
		panic("db connection name must be specified")
	}

	if len(dbName) == 0 {
		var d *gorm.DB
		for _, d = range db.connections {
			return d
		}
	}

	if len(dbName) > 0 {
		d, found := db.connections[dbName[0]]
		if !found {
			panic(fmt.Sprintf("not found db connection with dbName: %s", dbName[0]))
		}
		return d
	}
	return nil
}

type dbLogger struct {
	parentLogger logger.Logger
}

func newDBLogger() *dbLogger {
	return &dbLogger{
		parentLogger: logger.NewZapLogger(
			logger.Level(logger.DebugLevel),
			logger.ProjectName("Mysql"),
			logger.Type(logger.LogText),
		),
	}
}

func (this *dbLogger) LogMode(level dblog.LogLevel) dblog.Interface {
	return this
}

func (this *dbLogger) Info(ctx context.Context, s string, i ...interface{}) {
	slice := make([]interface{}, 0)
	slice = append(slice, "dbContext", ctx, i)
	this.parentLogger.Info(s, slice)
}

func (this *dbLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	slice := make([]interface{}, 0)
	slice = append(slice, "dbContext", ctx, i)
	this.parentLogger.Info(s, slice)
}

func (this *dbLogger) Error(ctx context.Context, s string, i ...interface{}) {
	slice := make([]interface{}, 0)
	slice = append(slice, "dbContext", ctx, i)
	this.parentLogger.Info(s, slice)
}

func (this *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	slice := make([]interface{}, 0)
	sql, _ := fc()
	slice = append(slice, "SQL", sql, "context", ctx)
	if err != nil {
		slice = append(slice, err)
		this.parentLogger.Error(slice)
	} else {
		this.parentLogger.Debug(slice)
	}
}
