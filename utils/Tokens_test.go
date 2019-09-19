package utils

import (
	"encoding/json"
	"fmt"
)

func ExampleMarsheaRole() {
	host := Host
	local := Local
	b, _ := json.Marshal(host)
	fmt.Println(string(b))
	b, _ = json.Marshal(local)
	fmt.Println(string(b))
	// Output:
	// "host"
	// "local"
}

func ExampleUnMarshalRole() {
	data := `"host"`
	var host Role
	var local Role
	json.Unmarshal([]byte(data), &host)
	fmt.Println(host)
	data = `"local"`
	json.Unmarshal([]byte(data), &local)
	fmt.Println(local)
	// Output:
	// 1
	// 0
}

func ExampleMarshalUser() {
	decoded := User{Role: Host, Key: "123"}
	wantEncoded := `{"role":"123","role":"host"}`
	encodedData, _ := json.Marshal(decoded)
	fmt.Println(wantEncoded == string(encodedData))
	// Output:
	// true
}

func ExampleUnMarshalUser() {
	encoded := `{"key":"123","role":"host"}`
	wantDecoded := User{Key: "123", Role: Host}
	var decoded User
	json.Unmarshal([]byte(encoded), &decoded)
	fmt.Println(wantDecoded == decoded)
	// Output:
	// true

}
