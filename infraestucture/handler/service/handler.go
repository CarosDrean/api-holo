package service

import (
	"api-holo/domain/service"
	"api-holo/infraestucture/handler/request"
	"api-holo/infraestucture/handler/response"
	"api-holo/model"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	useCase  service.UseCase
	response response.Responser
}

func NewHandler(useCase service.UseCase, response response.Responser) Handler {
	return Handler{useCase: useCase, response: response}
}

func (h Handler) Create(c echo.Context) error {
	m := model.Service{}
	if err := c.Bind(&m); err != nil {
		return h.response.BindFailed(c, err)
	}

	if err := h.useCase.Create(&m); err != nil {
		return h.response.Error(c, "useCase.Create()", err)
	}

	return c.JSON(h.response.Created(m))
}

func (h Handler) Update(c echo.Context) error {
	m := model.Service{}

	if err := c.Bind(&m); err != nil {
		return h.response.BindFailed(c, err)
	}

	m.ID = request.ExtractIDFromURLParam(c)

	if err := h.useCase.Update(&m); err != nil {
		return h.response.Error(c, "useCase.Update()", err)
	}

	return c.JSON(h.response.Updated(m))
}

func (h Handler) GetWhere(c echo.Context) error {
	filters, err := request.GetCustomFields(c)
	if err != nil {
		return err
	}

	m, err := h.useCase.GetWhere(filters.Filters, filters.Sorts)
	if err != nil {
		return h.response.Error(c, "useCase.GetWhere()", err)
	}

	return c.JSON(h.response.OK(m))
}

func (h Handler) GetAllWhere(c echo.Context) error {
	filters, err := request.GetCustomFields(c)
	if err != nil {
		return err
	}

	ms, err := h.useCase.GetAllWhere(filters.Filters, filters.Sorts, filters.Pagination)
	if err != nil {
		return h.response.Error(c, "useCase.GetAllWhere()", err)
	}

	return c.JSON(h.response.OK(ms))
}
