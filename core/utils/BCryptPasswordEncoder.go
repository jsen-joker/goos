package utils

import "golang.org/x/crypto/bcrypt"

func Encode(password string) (str string, err error) {
	if strBytes, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); e != nil {
		return "", e
	} else {
		return string(strBytes), nil
	}
}

func Matches(encodedPassword string, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(rawPassword))
}
