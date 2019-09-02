package backend

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func ResponseError(w http.ResponseWriter, statusCode int, responseBody string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(responseBody))
	return
}

func VaildToken(rawToken string) bool {
	token, _ := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("test"), nil
	})
	if claim, ok := token.Claims.(jwt.MapClaims); ok {
		if err := claim.Valid(); err != nil {
			return false
		}
	}
	return true
}