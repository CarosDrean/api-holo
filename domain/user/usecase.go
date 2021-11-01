package user

import (
	"fmt"

	"api-holo/model"
)

type User struct {
	storage Storage
}

func New(storage Storage) *User {
	return &User{storage: storage}
}

func (u User) Create(m *model.User) error {
	if err := u.storage.Create(m); err != nil {
		return fmt.Errorf("user: %v", err)
	}

	return nil
}

func (u User) Update(m *model.User) error {
	if err := u.storage.Update(m); err != nil {
		return fmt.Errorf("user: %v", err)
	}

	return nil
}

func (u User) GetWhere(filter model.Fields, sort model.SortFields) (model.User, error) {
	m, err := u.storage.GetWhere(filter, sort)
	if err != nil {
		return model.User{}, fmt.Errorf("user: %w", err)
	}

	return m, nil
}

func (u User) GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Users, error) {
	ms, err := u.storage.GetAllWhere(filter, sort, pag)
	if err != nil {
		return nil, fmt.Errorf("user: %v", err)
	}

	return ms, nil
}
