package service

import "api-holo/model"

type UseCase interface {
	Create(m *model.Service) error
	Update(m *model.Service) error
	GetWhere(filter model.Fields, sort model.SortFields) (model.Service, error)
	GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Services, error)
}

type Storage interface {
	Create(m *model.Service) error
	Update(m *model.Service) error
	GetWhere(filter model.Fields, sort model.SortFields) (model.Service, error)
	GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Services, error)
}
