package bootstrap

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"

	"api-holo/kit/authorization"
)

type Logger interface {
	Fatalf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

func loadSignatures(conf Configuration, logger Logger) *rsa.PrivateKey {
	priv := conf.PrivateFileSign
	publ := conf.PublicFileSign

	fpriv, err := ioutil.ReadFile(priv)
	checkErr(err, fmt.Sprintf("no se pudo leer el archivo de firma privado %s", priv))

	fpubl, err := ioutil.ReadFile(publ)
	checkErr(err, fmt.Sprintf("no se pudo leer el archivo de firma público %s", publ))

	authorization.LoadSignatures(fpriv, fpubl, logger)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(fpriv)
	checkErr(err, fmt.Sprintf("no se pudo convertir el archivo de firma privado %s", priv))

	return privateKey
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
