package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var ErrDatabaseInit = errors.New("database initialize fail")

var (
	once sync.Once
	db   *gorm.DB
)

// GetInstance
func GetInstance() *gorm.DB {
	once.Do(func() {
		db = gormDb()
	})
	return db
}

func gormDb() *gorm.DB {
	if config.CONFIG.System.DatabaseType == "postgres" {
		return gormPostgres()
	} else if config.CONFIG.System.DatabaseType == "sqlite" {
		return gormSqlite()
	} else {
		return gormMysql()
	}
}

func gormPostgres() *gorm.DB {
	if CONFIG_POSTGRES.DbName == "" {
		fmt.Println("conf dbname is empty")
		return nil
	}

	postgresConfig := postgres.Config{
		DSN:                  CONFIG_POSTGRES.Dsn(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}

	db, err := gorm.Open(postgres.New(postgresConfig), gormConfig(CONFIG_POSTGRES.LogMode, CONFIG_POSTGRES.LogZap))
	if err != nil {
		fmt.Printf("open postgres is failed %v \n", err)
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(CONFIG_POSTGRES.MaxIdleConns)
	sqlDB.SetMaxOpenConns(CONFIG_POSTGRES.MaxOpenConns)

	return db
}

func gormMysql() *gorm.DB {
	if CONFIG_MYSQL.DbName == "" {
		fmt.Println("conf dbname is empty")
		return nil
	}
	err := createTable(CONFIG_MYSQL.BaseDsn(), "mysql", CONFIG_MYSQL.DbName)
	if err != nil {
		fmt.Printf("create database %s is failed %v \n", CONFIG_MYSQL.DbName, err)
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       CONFIG_MYSQL.Dsn(),
		DefaultStringSize:         191,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(CONFIG_MYSQL.LogMode, CONFIG_MYSQL.LogZap))
	if err != nil {
		fmt.Printf("open mysql is failed %v \n", err)
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(CONFIG_MYSQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(CONFIG_MYSQL.MaxOpenConns)

	return db
}

func gormSqlite() *gorm.DB {
	if CONFIG_SQLITE.DbName == "" {
		fmt.Println("conf dbname is empty")
		return nil
	}

	conn := filepath.Join(consts.ExecDir, CONFIG_SQLITE.DbName+".db") // ?cache=shared&mode=rwc&_journal_mode=WAL")
	_logUtils.Infof(conn)
	dialector := sqlite.Open(conn)

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: false,
		Logger:                 logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
		},
	})

	if err != nil {
		_logUtils.Error(err.Error())
	}

	db.Exec("PRAGMA journal_mode=WAL;")

	err = db.Use(
		dbresolver.Register(
			dbresolver.Config{}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(CONFIG_SQLITE.MaxIdleConns).
			SetMaxOpenConns(CONFIG_SQLITE.MaxOpenConns),
	)
	if err != nil {
		_logUtils.Error(err.Error())
	}

	db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: false})

	return db
}

// gormConfig get gorm conf
func gormConfig(mod bool, logZap string) *gorm.Config {
	var config = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	switch logZap {
	case "silent", "Silent":
		config.Logger = Default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = Default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = Default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = Default.LogMode(logger.Info)
	case "zap", "Zap":
		config.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)

	default:
		if mod {
			config.Logger = Default.LogMode(logger.Info)
			break
		}
		config.Logger = Default.LogMode(logger.Silent)
	}

	return config
}

// createTable create database(mysql)
func createTable(dsn, driver, dbName string) error {
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", dbName)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)

	return err
}

func DorpDB(dsn, driver, dbName string) error {
	execSql := fmt.Sprintf("DROP database if exists `%s`;", dbName)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	if db == nil {
		return errors.New("db is nil")
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	if err = db.Ping(); err != nil {

		return err
	}

	_, err = db.Exec(execSql)
	if err != nil {
		return err
	}

	_logUtils.Debug(execSql)

	return nil
}
