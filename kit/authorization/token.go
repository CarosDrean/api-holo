package authorization

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Claim estructura de respuesta de token
type Claim struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	SessionID uint   `json:"session_id"`
	IPClient  string `json:"ip_client"`
	UserType  uint   `json:"user_type"`
	jwt.StandardClaims
}

// LoadSignatures Carga la información de los certificados de firma y confirmación
func LoadSignatures(private, public []byte, logger Logger) {
	once.Do(func() {
		var err error
		signKey, err = jwt.ParseRSAPrivateKeyFromPEM(private)
		if err != nil {
			logger.Fatalf("authorization.LoadSignatures: realizando el parse en jwt RSA private: %s", err)
		}

		verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(public)
		if err != nil {
			logger.Fatalf("authorization.LoadSignatures: realizando el parse en jwt RSA public: %s", err)
		}
	})
}

type AuthValidator struct {
	logger Logger
}

func NewAuthServiceValidator(logger Logger) *AuthValidator {
	return &AuthValidator{logger: logger}
}

// ValidateJWT Middleware para validar los JWT token
func (s *AuthValidator) ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var tokenString string
		tokenString, err := getTokenFromAuthorizationHeader(c.Request())
		if err != nil {
			tokenString, err = getTokenFromURLParams(c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Se encontró un error al tratar de leer el token")
			}
		}

		verifyFunction := func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claim{}, verifyFunction)
		if err != nil {
			status := http.StatusUnauthorized
			var msg string
			switch err.(type) {
			case *jwt.ValidationError:
				vErr := err.(*jwt.ValidationError)

				switch vErr.Errors {
				case jwt.ValidationErrorExpired:
					msg = "Su token ha expirado, por favor vuelva a ingresar"
				default:
					msg = "Error de validación del token"
				}
			default:
				status = http.StatusInternalServerError
				msg = "Error al procesar el token"
			}

			return echo.NewHTTPError(status, msg)
		}
		if !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token de acceso no válido")
		}

		userID := token.Claims.(*Claim).UserID
		email := token.Claims.(*Claim).Email
		userType := token.Claims.(*Claim).UserType

		c.Set("userID", userID)
		c.Set("email", email)
		c.Set("userType", userType)

		return next(c)
	}
}

// getTokenFromAuthorizationHeader busca el token del header Authorization
func getTokenFromAuthorizationHeader(r *http.Request) (string, error) {
	ah := r.Header.Get("Authorization")
	if ah == "" {
		return "", errors.New("el encabezado no contiene la autorización")
	}

	// Should be a bearer token
	if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
		return ah[7:], nil
	} else {
		return "", errors.New("el header no contiene la palabra Bearer")
	}
}

// getTokenFromURLParams busca el token de la URL
func getTokenFromURLParams(r *http.Request) (string, error) {
	ah := r.URL.Query().Get("authorization")
	if ah == "" {
		return "", errors.New("la URL no contiene la autorización")
	}

	return ah, nil
}
