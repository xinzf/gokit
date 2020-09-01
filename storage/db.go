package storage

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/xinzf/gokit/logger"
	"sync"
	"time"
)

type database struct {
	options    DbOptions
	connection *gorm.DB
}

var DB *database
var dbOnce sync.Once

func (db *database) Init(opts ...DbOption) error {
	var err error

	dbOnce.Do(func() {

		DB = &database{options: newDbOptions(opts...)}

		if DB.options.Addr == "" {
			err = errors.New("dbconfig's addr is empty")
			return
		}
		if DB.options.User == "" {
			err = errors.New("dbconfig's user is empty")
			return
		}
		if DB.options.Name == "" {
			err = errors.New("dbconfig's name is empty")
			return
		}
		if DB.options.Pswd == "" {
			err = errors.New("dbconfig's pswd is empty")
			return
		}

		var d *gorm.DB
		d, err = gorm.Open("mysql", DB.options.String())
		if err != nil {
			return
		}

		d.LogMode(DB.options.Log)
		if DB.options.logger != nil {
			d.SetLogger(DB.options.logger)
		}

		unixMilli := func(t time.Time) int64 {
			return t.UnixNano() / 1e6
		}

		d.DB().SetMaxIdleConns(DB.options.MaxIdleConns)
		d.DB().SetMaxOpenConns(DB.options.MaxOpenConns)
		d.DB().SetConnMaxLifetime(time.Duration(300) * time.Second)

		d.Callback().Create().Before("gorm:create").Register("set_created_updated", func(scope *gorm.Scope) {
			if scope.HasColumn("created") {
				scope.SetColumn("created", unixMilli(time.Now()))
			}
			if scope.HasColumn("updated") {
				scope.SetColumn("updated", unixMilli(time.Now()))
			}
		})
		d.Callback().Update().Before("gorm:update").Register("set_updated", func(scope *gorm.Scope) {
			if scope.HasColumn("updated") {
				scope.SetColumn("updated", unixMilli(time.Now()))
			}
		})

		d.SingularTable(true)

		if len(DB.options.initSql) > 0 {
			for _, sql := range DB.options.initSql {
				if err = d.Exec(sql).Error; err != nil {
					return
				}
			}
		}

		DB.connection = d
		DB.options.logger.Info("Mysql","DB inited.")

		go func() {
			defer func() {
				logger.DefaultLogger.Debug("Mysql","db exit.")
			}()
			for {
				time.Sleep(10*time.Second)
				if err := DB.connection.DB().Ping(); err != nil {
					logger.DefaultLogger.Error(err)
				} else {
					logger.DefaultLogger.Info("Mysql","db ping")
				}
			}
		}()
	})
	return err
}

func (db *database) Use() *gorm.DB {
	return db.connection
}

func (db *database) Close() {
	db.connection.Close()
}
