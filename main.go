package main

import (
	"database/sql"
	f "fmt"
	"task/api/database"
	"task/api/initializers"
	"task/api/router"
)

func init() {
	f.Println("Init starting")
	initializers.LoadEnvVariables()
	f.Println("Try to make connection to server")
	database.GormDBI()
}

func main() {
	var sqlDB *sql.DB
	f.Println("Main starting")
	defer sqlDB.Close()
	router.StartRouter()
}
