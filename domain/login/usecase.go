package login

import (
	"crypto/md5"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"api-holo/model"

	"github.com/dgrijalva/jwt-go"
)

type Login struct {
	useCaseUser UseCaseUser
	privateKey  *rsa.PrivateKey
}

func New(useCaseUser UseCaseUser, privateKey *rsa.PrivateKey) *Login {
	return &Login{useCaseUser: useCaseUser, privateKey: privateKey}
}

func (l Login) ValidateLogin(userName, password string) (model.UserLogin, error) {
	user, err := l.useCaseUser.GetWhere(model.Fields{model.Field{Name: "v_UserName", Value: userName}}, model.SortFields{})
	if errors.Is(err, sql.ErrNoRows) {
		customErr := model.Error{}
		customErr.SetAPIMessage("Usuario no encontrado.")
		customErr.SetError(err)

		return model.UserLogin{}, &customErr
	}
	if err != nil {
		return model.UserLogin{}, fmt.Errorf("login: get where user, err: %v", err)
	}

	if user.Password != encryptMD5(password) {
		customErr := model.Error{}
		customErr.SetAPIMessage("Contraseña invalida.")
		customErr.SetError(err)

		return model.UserLogin{}, &customErr
	}

	token, err := generateToken(user, l.privateKey)
	if err != nil {
		return model.UserLogin{}, fmt.Errorf("login: generate token, err: %v", err)
	}

	res := model.UserLogin{
		User:  user,
		Token: token,
	}

	return res, nil
}

func generateToken(user model.User, privateKey *rsa.PrivateKey) (string, error) {
	claims := Claim{
		ID:   user.ID,
		Type: user.Type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenSigned, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenSigned, err
}

func encryptMD5(text string) string {
	data := []byte(text)
	hash := md5.Sum(forHash(data))

	return base64.StdEncoding.EncodeToString(hash[:])
}

func forHash(data []byte) []byte {
	res := make([]byte, 0)
	for _, e := range data {
		res = append(res, e)
		res = append(res, 0)
	}

	return res
}
