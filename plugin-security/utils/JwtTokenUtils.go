package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var SecretKey = []byte("SecretKey012345678901234567890123456789012345678901234567890123456789")
// seconds
const TokenValidityInMilliseconds int64 = 60 * 30

type Token struct {
	Token string `json:"token"`
}

type MyCustomClaims struct {
	ID int64 `json:"id,omitempty"`
	Auth string `json:"auth,omitempty"`
	jwt.StandardClaims
}

func CreateToken(id int64, name string) (t *Token, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		id,
		"",
		jwt.StandardClaims{
			Subject: name,
			ExpiresAt: time.Now().UnixNano() / 1e9 + TokenValidityInMilliseconds,
		},
	})

	tk, err := token.SignedString(SecretKey)
	if err == nil {
		return &Token{Token: tk}, nil
	} else {
		return nil, err
	}
}

func GetSubject(tokenStr string) (id int64, auth string, subject string, err error)  {
	token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return SecretKey, nil
	})


	if err != nil {
		return 0, "", "", err
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return 0, "", "", errors.New("not MyCustomClaims")
	}

	return claims.ID, claims.Auth, claims.Subject, nil
}

func ValidToken(tokenStr string) error  {

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return SecretKey, nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return errors.New("not MyCustomClaims")
	}
	return claims.StandardClaims.Valid()
}