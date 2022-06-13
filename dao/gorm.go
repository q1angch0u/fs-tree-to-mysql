package dao

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"my-tree/config"
	"os"

	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// db连接
var db *gorm.DB

// Setup 初始化连接
func Setup() {
	// db = newConnection()
	var dbURI string
	var dialector gorm.Dialector

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Duration(100) * time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Silent,                    // 日志级别
			IgnoreRecordNotFoundError: true,                             // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                            // 禁用彩色打印
		},
	)

	dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		config.DatabaseSetting.User,
		config.DatabaseSetting.Password,
		config.DatabaseSetting.Host,
		config.DatabaseSetting.Port,
		config.DatabaseSetting.Name)
	dialector = mysql.New(mysql.Config{
		DSN:                       dbURI, // data source name
		DefaultStringSize:         1024,  // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	})
	conn, err := gorm.Open(dialector, &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Print(err.Error())
	}
	sqlDB, err := conn.DB()
	if err != nil {
		fmt.Errorf("connect db server failed.")
	}
	sqlDB.SetMaxIdleConns(100)                  // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxOpenConns(100)                  // SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetConnMaxLifetime(time.Second * 600) // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db = conn
}

// GetDB 开放给外部获得db连接
func GetDB() *gorm.DB {
	if db == nil {
		Setup()
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Errorf("connect db server failed.")
		Setup()
	}
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		Setup()
	}

	return db
}
