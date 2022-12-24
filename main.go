package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nguyentruyen98/go-be/api"
	db "github.com/nguyentruyen98/go-be/db/sqlc"
	"github.com/nguyentruyen98/go-be/util"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}

}
