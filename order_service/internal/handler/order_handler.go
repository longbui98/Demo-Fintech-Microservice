package handler

import (
	"encoding/json"
	"micro-project/order_service/internal/logic"
	"micro-project/order_service/internal/model"
	"micro-project/order_service/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateOrder(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.OrderArgs
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			httpx.Error(w, err)
			return
		}

		var res string
		orderLogic := logic.NewOrderLogic(r.Context(), ctx)
		err := orderLogic.CreateOrder(&user, &res)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, map[string]string{"message": "Created order successfully"})
	}
}

func CreateOrders(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []model.OrderArgs

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Error(w, err)
			return
		}

		var res string
		userLogic := logic.NewOrderLogic(r.Context(), ctx)
		err := userLogic.CreateOrders(&req, &res)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, map[string]string{"message": "Created orders successfully"})
	}
}

func GetOrderByEmail(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.Email

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Error(w, err)
			return
		}

		var res []model.OrderResponse
		userLogic := logic.NewOrderLogic(r.Context(), ctx)
		err := userLogic.GetOrderByEmail(&req, &res)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, res)
	}
}

func DeleteOrder(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.OrderId

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Error(w, err)
			return
		}

		var res string
		userLogic := logic.NewOrderLogic(r.Context(), ctx)
		err := userLogic.DeleteOrderByOrderId(&req, &res)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, map[string]string{"message": "Deleted orders successfully"})
	}
}
