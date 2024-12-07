package routes

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/satyendra001/mdm-oauth/user"
)

type RouteHandler struct {
	conn *pgxpool.Pool
	ctx  context.Context
}

func NewRouteHandler(connPool *pgxpool.Pool, context context.Context) *RouteHandler {
	return &RouteHandler{
		conn: connPool,
		ctx:  context,
	}
}

func (h *RouteHandler) HandleRoutes(mux *mux.Router) {

	userStore := user.NewUserStore(h.conn, h.ctx)
	userStore.RegisterOauthRoutes(mux)

}
