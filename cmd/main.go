package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/satyendra001/mdm-oauth/cmd/api"
	core "github.com/satyendra001/mdm-oauth/core/utils"
	"github.com/satyendra001/mdm-oauth/utils"
)

// Function to intiate the DB connection
func initDB(conn *pgxpool.Pool, ctx context.Context) {
	if err := conn.Ping(ctx); err != nil {
		log.Fatal("DB Connection error => ", err)
	}

	log.Println("DB Connection Successfull...")
}

func main() {
	dbConf := utils.DBEnvs
	ctx := context.Background()
	connPool, err := core.NewPostgreSQLConnection(ctx, 2, 5, &dbConf)
	if err != nil {
		log.Fatal("Unable to create PostreSQL connection :(")
	}

	initDB(connPool, ctx)

	defer connPool.Close()

	server := api.NewAPIServer("dmt.localhost:3000", connPool)

	if err := server.Run(ctx); err != nil {
		log.Fatal("Error while runnning the server ==> ", err.Error())
	}
}
