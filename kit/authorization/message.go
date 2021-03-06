package authorization

// Codes of message
const (
	ForbiddenCode     = "FORBIDDEN"
	InvalidToken      = "invalid_token"
	UnauthorizedRoles = "unauthorized_roles"
)

// Types of message

const (
	ForbiddenType = "Forbidden"
)

// ResponseData Message dof response
type ResponseData struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Href    string `json:"href"`
}

// ResponseMessage contained messages and errors
type ResponseMessage struct {
	Data     interface{}    `json:"data"`
	Messages []ResponseData `json:"messages"`
	Errors   []ResponseData `json:"errors"`
}

// AddMessage permite agregar un mensaje a ResponseMessage
func (rm *ResponseMessage) AddMessage(t, c, m, h string) {
	rd := ResponseData{
		Type:    t,
		Code:    c,
		Message: m,
		Href:    h,
	}
	rm.Messages = append(rm.Messages, rd)
}

// AddError permite agregar un error a ResponseMessage
func (rm *ResponseMessage) AddError(t, c, m, h string) {
	rd := ResponseData{
		Type:    t,
		Code:    c,
		Message: m,
		Href:    h,
	}
	rm.Errors = append(rm.Errors, rd)
}
