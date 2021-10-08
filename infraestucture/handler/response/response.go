package response

import (
	"net/http"

	"api-holo/model"

	"github.com/labstack/echo/v4"
)

// HTTPErrorHandler handler the error response of echo
func HTTPErrorHandler(err error, c echo.Context) {
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
