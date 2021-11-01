package login

import "api-holo/model"

type UseCase interface {
	ValidateLogin(userName, password string) (model.UserLogin, error)
}

type UseCaseUser interface {
	GetWhere(filter model.Fields, sort model.SortFields) (model.User, error)
}
