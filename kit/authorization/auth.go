package authorization

import (
	"crypto/rsa"
	"sync"

	"github.com/labstack/echo/v4"
)

type Logger interface {
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type AuthMiddleware interface {
	ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc
}

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
	once      sync.Once
)
