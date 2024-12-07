package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/satyendra001/mdm-oauth/routes"
)

type APIServer struct {
	addr     string
	connPool *pgxpool.Pool
}

func NewAPIServer(addr string, conn *pgxpool.Pool) *APIServer {
	return &APIServer{
		addr:     addr,
		connPool: conn,
	}
}

// This type of method is called method with a receiver. Here Run receives a pointer to APIServer
func (server *APIServer) Run(ctx context.Context) error {
	router := mux.NewRouter()

	routeHandler := routes.NewRouteHandler(server.connPool, ctx)
	routeHandler.HandleRoutes(router)

	log.Println("Listening on", server.addr)
	return http.ListenAndServe(server.addr, router)
}
