package handler

import (
	"encoding/json"
	"micro-project/user_service/internal/logic"
	"micro-project/user_service/internal/model"
	"micro-project/user_service/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			httpx.Error(w, err)
			return
		}

		userLogic := logic.NewUserLogic(r.Context(), ctx)
		err := userLogic.Register(user)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, map[string]string{"message": "User registered successfully"})
	}
}

func LoginHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Error(w, err)
			return
		}

		userLogic := logic.NewUserLogic(r.Context(), ctx)
		user, err := userLogic.Login(req.Email, req.Password)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, user)
	}
}
