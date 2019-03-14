package utils

import (
	"fmt"
	"testing"
)

func TestCreateToken(t *testing.T) {
	name := "jsen"
	token, err := CreateToken(name)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(token.Token)
	}
}

func TestGetSubject(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTA1NzE1MzgsInN1YiI6ImpzZW4ifQ.OaXNf14WDI7-QWvPEqpoZpxleCnOHvjGjUDsmMV4Jh4"

	subject, err := GetSubject(token)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(subject)
	}
}
func TestVaildToken(t *testing.T) {
	// signature is invalid 错误的token
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTA1Njk0NTQsInN1YiI6ImpzZW4ifQ.-MypPNfT2f2-mXxcNVjjtxtpNSBQDnjXtHOWOeAtyzM"
	err := ValidToken(token)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ok")
	}
}