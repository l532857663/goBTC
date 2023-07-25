package db

import (
	"fmt"
	"goBTC/db/internal"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DBS = make(map[string]*gorm.DB)
	LOG *zap.Logger
)

func GetDBByName(dbName string) *gorm.DB {
	return DBS[dbName]
}

func InitDaoByName(m Mysql, config *gorm.Config, dbName string) {
	DBS[dbName] = GormMysql(m, config)
}

// 事务控制
func GetDBTxnByName(dbName string) *gorm.DB {
	return DBS[dbName].Begin()
}

func Gorm(conf Mysql, zapLogger *zap.Logger, dbName string) {
	mysqlConfig := gormConfig(conf, zapLogger)
	InitDaoByName(conf, mysqlConfig, dbName)
	zapLogger.Debug("InitDao:", zap.Any("DbName", dbName))
}

func GormMysql(m Mysql, config *gorm.Config) *gorm.DB {
	//fmt.Println("m.Username: ", m.Username)
	//fmt.Println("m.Password: ", m.Password)
	//fmt.Println("m.Path: ", m.Path)
	//fmt.Println("m.Dbname: ", m.Dbname)
	//fmt.Println("m.Config: ", m.Config)
	dsn := m.Dsn()
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	//fmt.Println("dsn: ", dsn)
	if db, err := gorm.Open(mysql.New(mysqlConfig), config); err != nil {
		fmt.Println("MySQL启动异常 ", err.Error())
		os.Exit(0)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		// fmt.Println("数据库连接成功.")
		return db
	}
}

var ErrorNoDataFind = fmt.Errorf("no data find")

var ErrorValueHasBeenSettled = fmt.Errorf("value has been settled")

var ErrorAffectZeroRow = fmt.Errorf("insert/update affect 0 row")

var ErrorDuplicateValues = fmt.Errorf("duplicate record")

var ErrorNilTransaction = fmt.Errorf("transaction is nil")

var ErrorColumnMustNotBeEmpty = func(column string) error {
	return fmt.Errorf("column '%s' must not be empty", column)
}

// @Description 数据库退出处理函数
// @Author Oracle
// @Version 1.0
// @Update Oracle 2021-12-18 init
func CloseAllGormDBConnections() {
	for k := range DBS {
		if DBS[k] != nil {
			db, err := DBS[k].DB()
			if err != nil {
				fmt.Printf("Close DB[%s] connections error: %+v\n", k, err)
			}
			fmt.Printf("Exiting DB[%s] connections\n", k)
			db.Close()
		}
	}
}

func gormConfig(m Mysql, zapLogger *zap.Logger) *gorm.Config {
	internal.Init(m.LogZap, zapLogger)

	var fig = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	switch m.LogZap {
	case "silent", "Silent":
		fig.Logger = internal.Default.LogMode(logger.Silent)
	case "error", "Error":
		fig.Logger = internal.Default.LogMode(logger.Error)
	case "warn", "Warn":
		fig.Logger = internal.Default.LogMode(logger.Warn)
	case "info", "Info":
		fig.Logger = internal.Default.LogMode(logger.Info)
	case "zap", "Zap":
		fig.Logger = internal.Default.LogMode(logger.Info)
	default:
		if m.LogMode {
			fig.Logger = internal.Default.LogMode(logger.Info)
			break
		}
		fig.Logger = internal.Default.LogMode(logger.Silent)
	}
	return fig
}
