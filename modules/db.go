package module

import (
	"fmt"
	"log"
	"os"
	"time"

	common "api-test/common"

	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var DBPolicy *gorm.DB

func DbConnection() {
	common.LoadEnv()
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	PORT := os.Getenv("DB_PORT")
	DBNAME := os.Getenv("DB_NAME")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)

	Database, err := gorm.Open(mysql.Open(url), &gorm.Config{Logger: newLogger})

	if err != nil {
		common.WriteLogWithoutContext("Couldnt connect to Database "+err.Error(), zapcore.FatalLevel, "FATAL")
		panic("Couldnt connect to Database")
	}
	common.WriteLogWithoutContext("Connection to Central DB was established ", zapcore.InfoLevel, "INFO")
	DB = Database
}
