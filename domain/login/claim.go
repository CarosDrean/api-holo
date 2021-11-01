package login

import (
	"api-holo/model"

	"github.com/dgrijalva/jwt-go"
)

const issuer = "api"

type Claim struct {
	ID   string         `json:"id"`
	Type model.UserType `json:"type"`
	jwt.StandardClaims
}
