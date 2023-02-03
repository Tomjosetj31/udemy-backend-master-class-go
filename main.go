package main

import (
	"database/sql"
	"log"

	"github.com/Tomjosetj31/simple-bank/api"
	simplebank "github.com/Tomjosetj31/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)


const (
	dbDriver="postgres"
	dbSource="postgresql://tom:secret@localhost:5432/samplebank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := simplebank.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
