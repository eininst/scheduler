package data

import (
	"database/sql"
	"encoding/json"
	"github.com/eininst/scheduler/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func NewDB() *gorm.DB {
	var dbconfig struct {
		Dsn          string        `json:"dsn"`
		MaxIdleCount int           `json:"maxIdleCount"`
		MaxOpenCount int           `json:"maxOpenCount"`
		MaxLifetime  time.Duration `json:"maxLifetime"`
	}

	mstr := configs.Get("mysql").String()
	_ = json.Unmarshal([]byte(mstr), &dbconfig)

	sqlDB, err := sql.Open("mysql", dbconfig.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   false,
		},
		CreateBatchSize: 100,
	})

	if err != nil {
		panic(err)
	}
	perr := sqlDB.Ping()
	if perr != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(dbconfig.MaxIdleCount)
	sqlDB.SetMaxOpenConns(dbconfig.MaxOpenCount)
	sqlDB.SetConnMaxLifetime(dbconfig.MaxLifetime * time.Second)

	log.Println("Connected to Mysql server...")

	return gormDB
}
