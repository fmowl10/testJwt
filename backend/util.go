package backend

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// ResponseError write statusCode to header and responseBody to body
func ResponseError(w http.ResponseWriter, statusCode int, responseBody string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(responseBody))
	return
}

// VaildToken check is rawToken vaild
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
