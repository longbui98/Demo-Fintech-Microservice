package internal

import (
	"micro-project/order_service/internal/handler"
	"micro-project/order_service/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterRoutes(server *rest.Server, ctx *svc.ServiceContext) {
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/order",
		Handler: handler.CreateOrder(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/orders",
		Handler: handler.CreateOrders(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/order",
		Handler: handler.DeleteOrder(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/order",
		Handler: handler.GetOrderByEmail(ctx),
	})
}
