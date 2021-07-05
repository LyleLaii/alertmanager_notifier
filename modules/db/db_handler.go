package db

import (
	"alertmanager_notifier/config"
	anLog "alertmanager_notifier/log"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"log"
	"os"
	"regexp"
	"time"

	// import postgres driver dependent
	"gorm.io/driver/postgres"

	"github.com/spf13/viper"
)

// DBInstance  DB Instance
var DBInstance *gorm.DB

// InitDBHandler init db hander
func InitDBHandler(logger anLog.Logger, runConf *config.RunningConfig) {
	var err error
	var defaultLogger gormlogger.Interface
	if runConf.RunMode.String() == "dev" {
		defaultLogger = gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormlogger.Config{
				SlowThreshold: time.Duration(100) * time.Millisecond, // 慢 SQL 阈值
				LogLevel:      gormlogger.Info,    // Log loglevel
				Colorful:      true,    // 禁用彩色打印
			},
		)
	} else {
		defaultLogger = NewDBLogger(logger, runConf)
	}

	//TODO: ugly, fix it
	if viper.IsSet("database.postgres") {
		var dbConfig PostgresConfig
		if e := viper.UnmarshalKey("database.postgres",&dbConfig); e != nil {
			logger.Panic("DBHandler", fmt.Sprintf("Error load databese config file: %s", e))
		}
		dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.Dbname,
			dbConfig.Password,
			dbConfig.Sslmode)
		reg := regexp.MustCompile(`password=(.*)[\s]?`)
		infostr := reg.ReplaceAllString(dsn, "password=******")
		logger.Info("DBHandler", fmt.Sprintf("database connect info: %v", infostr))

		DBInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: defaultLogger})
		if err != nil {
			logger.Panic("DBHandler", fmt.Sprintf("con not connect database cause: %v", err))
		}
	} else if viper.IsSet("database.sqlite") {
		var dbConfig SqliteConfig
		if e := viper.UnmarshalKey("database.sqlite",&dbConfig); e != nil {
			logger.Panic("DBHandler", fmt.Sprintf("Error load databese config file: %s", e))
		}
		logger.Info("DBHandler", fmt.Sprintf("database connect info: %v", dbConfig.DataPath))
		DBInstance, err = gorm.Open(sqlite.Open(dbConfig.DataPath), &gorm.Config{Logger: defaultLogger})
		if err != nil {
			logger.Panic("DBHandler", fmt.Sprintf("con not connect database cause: %v", err))
		}
	} else {
		logger.Panic("DBHandler", fmt.Sprintf("con not find database config"))
	}

	DBInstance.AutoMigrate(&Admin{}, &UserInfo{}, &ReceiverInfo{}, &UserRota{})

}

// CloseDB close db instance
// func CloseDB() {
// 	log.Info("DBHandler", "close db")
// 	defer DBInstance.Close()
// }
