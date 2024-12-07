package core

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/satyendra001/mdm-oauth/utils"
)

const (
	// DatabaseConnectionTimeOut is the default timeout for established connections.
	DatabaseConnectionTimeOut = 2 * time.Second
)

func NewPostgreSQLConnection(ctx context.Context, minPoolSize int32, maxPoolSize int32, conf *utils.DBConfig) (*pgxpool.Pool, error) {
	log.Println("Started creating a DB connection")
	dsn := fmt.Sprintf(
		// "postgresql://%s:%s@%s:%d/%s?statement_cache_mode=describe&sslmode=disable",
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database,
	)

	context, cancel := context.WithTimeout(ctx, DatabaseConnectionTimeOut)
	defer cancel()

	connConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Println("Error while parsing the config for DB connection => ", err.Error())
		return nil, err
	}

	if maxPoolSize == 0 {
		maxPoolSize = 1 // Default pool size is set to 1
	}
	if minPoolSize == 0 {
		minPoolSize = 1 // Default pool size is set to 1
	}
	if conf.Scheme == "" {
		conf.Scheme = "public" // Default scheme is public.
	}
	connConfig.MaxConns = maxPoolSize
	connConfig.MinConns = minPoolSize
	// connConfig.ConnConfig.PreferSimpleProtocol = true

	connPool, err := pgxpool.NewWithConfig(context, connConfig)
	if err != nil {
		log.Println("Error while creating the New pgxpool connection ==> ", err.Error())
		return nil, err
	}

	log.Println("Successfully created DB Pool. Returning the conn pool")

	return connPool, nil
}
