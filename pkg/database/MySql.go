package database

import (
	"context"
	"database/sql"
	"fmt"
	hi "github.com/denisenkom/go-mssqldb"
	"log"
)

// Replace with your own connection parameters
var server = "localhost"
var port = 1433
var user = "sa"
var password = "your_password"

func ConnectToDB(db *sql.DB) {
	var err error

	fmt.Println(hi.Error{
		Number:     0,
		State:      0,
		Class:      0,
		Message:    "",
		ServerName: "",
		ProcName:   "",
		LineNo:     0,
		All:        nil,
	})
	// Create connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d",
		server, user, password, port)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	log.Printf("Connected!\n")

	// Close the database connection pool after program executes
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
}

func CloseConnection(db *sql.DB){
	db.Close()
}

// SelectVersion Gets and prints SQL Server version
func SelectVersion(db *sql.DB){
	// Use background context
	ctx := context.Background()

	// Ping database to see if it's still alive.
	// Important for handling network issues and long queries.
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}

	var result string

	// Run query and scan for result
	err = db.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	fmt.Printf("%s\n", result)
}