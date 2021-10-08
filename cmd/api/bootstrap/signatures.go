package bootstrap

import (
	"api-holo/kit/authorization"
	"fmt"
	"io/ioutil"
	"log"
)

type Logger interface {
	Fatalf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

func loadSignatures(conf Configuration, logger Logger) {
	priv := conf.PrivateFileSign
	publ := conf.PublicFileSign

	fpriv, err := ioutil.ReadFile(priv)
	checkErr(err, fmt.Sprintf("no se pudo leer el archivo de firma privado %s", priv))

	fpubl, err := ioutil.ReadFile(publ)
	checkErr(err, fmt.Sprintf("no se pudo leer el archivo de firma público %s", publ))

	authorization.LoadSignatures(fpriv, fpubl, logger)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
