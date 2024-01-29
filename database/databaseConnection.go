package database

import (
	"database/sql"
	f "fmt"
	"log"
	"os"
	"task/api/models"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//	func SqlDBI(database string) *sql.DB {
//		pswd := os.Getenv("MYSQL_PWD")
//		dsn := "root:" + pswd + "@tcp(127.0.0.1:3306)/" + database
//		db, err := sql.Open("mysql", dsn)
//		if err != nil {
//			f.Println(err)
//			panic(err)
//		}
//		f.Println("connected to database")
//		return db
//	}

var DB *gorm.DB

func GormDBI() {
	dsn := os.Getenv("DSN")
	f.Println(dsn)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		f.Println("Failed to connect to database")
		panic(err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	DefConn(sqlDB)
	DB.AutoMigrate(&models.User{}, &models.Picture{})

	f.Println("Connected to database")
}

func DefConn(db *sql.DB) {
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetConnMaxLifetime(time.Minute * 15)

}
