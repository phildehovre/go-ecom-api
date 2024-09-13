package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/phildehovre/go-complete-api/cmd/api"
	"github.com/phildehovre/go-complete-api/config"
	"github.com/phildehovre/go-complete-api/db"
)

func main() {
	// Create new instance of DB
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Initialise DB
	initStorage(db)

	// Create new server instance
	server := api.NewAPIServer(":8080", db)
	// Initialise server
	if err := server.Run(); err != nil {
		log.Fatalf("server failed to run  %v", err)
	}

}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
