package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/yosa12978/halo/internal/pkg/dto"
	"github.com/yosa12978/halo/internal/pkg/repositories"
	"github.com/yosa12978/halo/pkg/helpers"
)

type IAuthHandler interface {
	LogIn(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct{}

func NewAuthHandler() IAuthHandler {
	return &AuthHandler{}
}

func (ah *AuthHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var userd dto.LoginUser
	json.NewDecoder(r.Body).Decode(&userd)
	if err := validator.New().Struct(userd); err != nil {
		helpers.RespondStatusCode(w, 400, "bad request")
		return
	}

	ur := repositories.NewUserRepository()
	token, err := ur.Login(userd)
	if err != nil {
		helpers.RespondStatusCode(w, 404, err.Error())
		return
	}
	helpers.RespondJson(w, 200, map[string]string{"token": token})
}

func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var userd dto.CreateUser
	json.NewDecoder(r.Body).Decode(&userd)
	if err := validator.New().Struct(userd); err != nil {
		helpers.RespondStatusCode(w, 400, "bad request")
		return
	}
	ur := repositories.NewUserRepository()
	if err := ur.CreateUser(userd); err != nil {
		helpers.RespondStatusCode(w, 400, err.Error())
		return
	}
	helpers.RespondStatusCode(w, 201, "created")
}
