package service

import (
	"fmt"

	"api-holo/model"
)

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{storage: storage}
}

func (c Service) Create(m *model.Service) error {
	if err := c.storage.Create(m); err != nil {
		return fmt.Errorf("service: %v", err)
	}

	return nil
}

func (c Service) Update(m *model.Service) error {
	if err := c.storage.Update(m); err != nil {
		return fmt.Errorf("service: %v", err)
	}

	return nil
}

func (c Service) GetWhere(filter model.Fields, sort model.SortFields) (model.Service, error) {
	service, err := c.storage.GetWhere(filter, sort)
	if err != nil {
		return model.Service{}, fmt.Errorf("service: %v", err)
	}

	return service, nil
}

func (c Service) GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Services, error) {
	services, err := c.storage.GetAllWhere(filter, sort, pag)
	if err != nil {
		return nil, fmt.Errorf("service: %v", err)
	}

	return services, nil
}
