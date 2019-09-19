package utils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// Role have two state
// Local and Host
type Role int

const (
	Local Role = iota
	Host
)

type User struct {
	Key  string `json:"key"`
	Role Role   `json:"role"`
}

type JwtClaim struct {
	User
	jwt.StandardClaims
}

func (r Role) MarshalJSON() ([]byte, error) {
	switch r {
	case Local:
		return []byte(`"local"`), nil
	case Host:
		return []byte(`"host"`), nil
	}
	return nil, fmt.Errorf("its type isn't Role")
}

func (r *Role) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"local"`:
		*r = Local
		return nil
	case `"host"`:
		*r = Host
		return nil
	}
	return fmt.Errorf("its type isn't Role")
}
