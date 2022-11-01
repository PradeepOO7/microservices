package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "8080"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	//Connect to DB
	connect:=ConnectToDB()
	if connect ==nil{
		log.Println("Cannot connect to Postgres! ")
	}

	app := Config{
		DB:connect,
		Models: data.New(connect),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Panicf("Something went wrong %s\n", err)
	}
	log.Println("Started Authentication Service")

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	counter := 0
	for {
		connect, err := openDB(dsn)
		if err != nil {
			log.Printf(" Database not yet ready %s\n", err)
			counter++
		}else{
			log.Printf("Coonected to DB")
			return connect
		}
		if counter > 10 {
			log.Panicf("Can't Establish a DB connection %s\n", err)
			return nil
		}
		log.Println("Waiting 2 sec before reconnecting ")
		time.Sleep(2*time.Second)
		continue
	}

}
