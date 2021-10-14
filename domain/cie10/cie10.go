package cie10

import "api-holo/model"

type UseCase interface {
	Create(m *model.Cie10) error
	Update(m *model.Cie10) error
	GetWhere(filter model.Fields, sort model.SortFields) (model.Cie10, error)
	GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Cie10s, error)
}

type Storage interface {
	Create(m *model.Cie10) error
	Update(m *model.Cie10) error
	GetWhere(filter model.Fields, sort model.SortFields) (model.Cie10, error)
	GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Cie10s, error)
}
