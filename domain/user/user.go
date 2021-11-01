package user

import "api-holo/model"

type UseCase interface {
	Create(m *model.User) error
	Update(m *model.User) error
	GetWhere(filter model.Fields, sort model.SortFields) (model.User, error)
	GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Users, error)
}

type Storage interface {
	Create(m *model.User) error
	Update(m *model.User) error
	GetWhere(filter model.Fields, sort model.SortFields) (model.User, error)
	GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Users, error)
}
