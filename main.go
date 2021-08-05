package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/rossi1/coding-challenge/cmd"
	"github.com/rossi1/coding-challenge/operation"
)

var (
	host     = os.Getenv("DB_HOST")
	port, _  = strconv.Atoi(os.Getenv("DB_PORT"))
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
)

func main() {
	db, err := DBInstance()

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	ops := operation.NewToken() // instiantate token

	err = cmd.Execute(db, ops)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)

	}

}

// DBInstance connects to database instance
func DBInstance() (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err

	}

	return db, nil

}
