package response

import (
	"net/http"

	"api-holo/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// HTTPErrorHandler handler the error response of echo
func HTTPErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok {
		var msgErr string
		msgErr, ok := he.Message.(string)
		if !ok {
			msgErr = "¡Upps! algo inesperado ocurrió"
		}

		registryLogError(c, he.Code, err)
		err = c.JSON(he.Code, model.Response{Message: msgErr})

		return
	}

	// custom error
	e, ok := err.(*model.Error)
	if ok {
		err = c.JSON(getResponseError(e))
		return
	}

	// if the handler not returns a "model.Error" then it returns a generic error JSON response
	err = c.JSON(http.StatusInternalServerError, model.MessageResponse{
		Errors: model.Responses{
			{Code: UnexpectedError, Message: "¡Uy! metimos la pata, disculpanos lo solucionaremos"},
		},
	})
}

// getResponseError returns the status code and a Response
func getResponseError(err *model.Error) (outputStatus int, outputResponse model.MessageResponse) {
	if !err.HasCode() {
		err.SetCode(UnexpectedError)
	}

	if !err.HasAPIMessage() {
		err.SetErrorAsAPIMessage()
	}

	if err.HasData() {
		outputResponse.Data = err.Data()
	}

	if !err.HasStatusHTTP() {
		err.SetStatusHTTP(http.StatusInternalServerError)
	}

	outputStatus = err.StatusHTTP()
	outputResponse.Errors = model.Responses{model.Response{
		Code:    err.Code(),
		Message: err.APIMessage(),
	}}

	return
}

func registryLogError(c echo.Context, status int, err error) {
	fields := logrus.Fields{
		"status":      status,
		"uri":         c.Path(),
		"query_param": c.QueryParams(),
		"remote_ip":   c.RealIP(),
		"method":      c.Request().Method,
	}

	if hasToken(c.Request()) {
		fields["user"] = GetUserID(c)
	}

	if e, ok := err.(*model.Error); ok {
		fields["where"] = e.Where()
		fields["who"] = e.Who()
	}

	if status >= 500 {
		logrus.WithFields(fields).Error(err)
		return
	}

	if status >= 400 {
		logrus.WithFields(fields).Warn(err)
		return
	}
}

func GetUserID(c echo.Context) uint {
	userID := c.Get("userID")
	if userID != nil {
		return userID.(uint)
	}
	return 0
}

func hasToken(r *http.Request) bool {
	return r.Header.Get("Authorization") != ""
}
