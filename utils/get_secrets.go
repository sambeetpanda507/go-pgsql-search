package utils

import "os"

type Secrets struct {
	PORT string
}

func GetSecrets() Secrets {
	secrets := Secrets{
		PORT: os.Getenv("PORT"),
	}
	return secrets
}
