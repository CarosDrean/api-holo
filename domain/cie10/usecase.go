package cie10

import (
	"fmt"

	"api-holo/model"
)

type Cie10 struct {
	storage Storage
	logger  model.Logger
}

func New(storage Storage, logger model.Logger) *Cie10 {
	return &Cie10{storage: storage, logger: logger}
}

func (c Cie10) Create(m *model.Cie10) error {
	if err := c.storage.Create(m); err != nil {
		return fmt.Errorf("cie10: %v", err)
	}

	return nil
}

func (c Cie10) Update(m *model.Cie10) error {
	if err := c.storage.Update(m); err != nil {
		return fmt.Errorf("cie10: %v", err)
	}

	return nil
}

func (c Cie10) GetWhere(filter model.Fields, sort model.SortFields) (model.Cie10, error) {
	cie10, err := c.storage.GetWhere(filter, sort)
	if err != nil {
		return model.Cie10{}, fmt.Errorf("cie10: %v", err)
	}

	return cie10, nil
}

func (c Cie10) GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Cie10s, error) {
	cie10s, err := c.storage.GetAllWhere(filter, sort, pag)
	if err != nil {
		return nil, fmt.Errorf("cie10: %v", err)
	}

	return cie10s, nil
}
