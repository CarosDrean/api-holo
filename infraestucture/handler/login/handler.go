package login

import (
	"api-holo/domain/login"
	"api-holo/infraestucture/handler/response"
	"api-holo/model"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	useCase  login.UseCase
	response response.Responser
}

func NewHandler(useCase login.UseCase, response response.Responser) Handler {
	return Handler{useCase: useCase, response: response}
}

func (h Handler) ValidateLogin(c echo.Context) error {
	m := model.User{}
	if err := c.Bind(&m); err != nil {
		return h.response.BindFailed(c, err)
	}

	userLogin, err := h.useCase.ValidateLogin(m.UserName, m.Password)
	if err != nil {
		return h.response.Error(c, "useCase.ValidateLogin()", err)
	}

	return c.JSON(h.response.OK(userLogin))
}
