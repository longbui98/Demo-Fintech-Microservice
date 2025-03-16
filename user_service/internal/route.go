package internal

import (
	"micro-project/user_service/internal/handler"
	"micro-project/user_service/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterRoutes(server *rest.Server, ctx *svc.ServiceContext) {
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/user/register",
		Handler: handler.RegisterHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/user/login",
		Handler: handler.LoginHandler(ctx),
	})
}
