package request

import (
	"encoding/json"
	"fmt"

	"api-holo/model"

	"github.com/labstack/echo/v4"
)

func GetCustomFields(c echo.Context) (model.CustomFields, error) {
	filters := c.QueryParam("filters")
	fields := model.Fields{}
	if filters != "" {
		err := json.Unmarshal([]byte(filters), &fields)
		if err != nil {
			return model.CustomFields{}, fmt.Errorf("invalid filter parameter")
		}
	}

	sorts := c.QueryParam("sorts")
	sortsFields := model.SortFields{}
	if sorts != "" {
		err := json.Unmarshal([]byte(sorts), &sortsFields)
		if err != nil {
			return model.CustomFields{}, fmt.Errorf("invalid filter parameter")
		}
	}

	pagination := c.QueryParam("pagination")
	pag := model.Pagination{}
	if pagination != "" {
		err := json.Unmarshal([]byte(pagination), &pag)
		if err != nil {
			return model.CustomFields{}, fmt.Errorf("invalid filter parameter")
		}
	}

	return model.CustomFields{
		Filters:    fields,
		Sorts:      sortsFields,
		Pagination: pag,
	}, nil
}
