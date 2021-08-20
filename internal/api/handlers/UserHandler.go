package handlers

type IUserHandler interface{}

type UserHandlers struct{}

func NewUserHandlers() IUserHandler {
	return &UserHandlers{}
}
