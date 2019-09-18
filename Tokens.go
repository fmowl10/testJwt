package testJwt

import (
	"encoding/json"
	"fmt"
	"log"
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

type jwtClaim struct {
	User
	jwt.StandardClaims
}

func (r Role) MashalJSON() ([]byte, error) {
	switch r {
	case Local:
		return []byte("local"), nil
	case Host:
		return []byte("host"), nil
	}
	return nil, fmt.Errorf("its type isn't Role")
}

func (r *Role) UnmashalJSON(data []byte) error {
	switch string(data) {
	case "local":
		*r = Local
		return nil
	case "host":
		*r = Host
		return nil
	}
	return fmt.Errorf("its type isn't Role")
}

func (u User) MashalJSON() ([]byte, error) {
	rawData, err := u.Role.MashalJSON()
	if err != nil {
		return nil, fmt.Errorf("Role Error : " + err.Error())
	}
	rawData = append([]byte(`"key":`+u.Key), rawData...)
	return rawData, nil
}

func (u *User) UnmashalJSON(data []byte) error {
	var mapData map[string]string
	err := json.Unmarshal(data, &mapData)
	if err != nil {
		return err
	}
	u.Key = mapData["key"]
	log.Println("dd")
	err = json.Unmarshal([]byte(mapData["role"]), &u.Role)
	if err != nil {
		return err
	}
	return nil
}