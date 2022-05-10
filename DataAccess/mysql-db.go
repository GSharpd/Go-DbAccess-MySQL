package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mySQLConnectionString := "golang:golang@/devbook?charset=utf8&parseTime=true&loc=Local"

	db, err := sql.Open("mysql", mySQLConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection established")

	lines, err := db.Query("select * from users")
	if err != nil {
		log.Fatal(err)
	}
	defer lines.Close()
}
