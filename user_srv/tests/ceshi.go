package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"strings"
)

//
func Password(code string) string {
	// Using custom options
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("admin123", options)
	fmt.Println(salt,encodedPwd)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)
	passwordInfo := strings.Split(newPassword,"$")
	fmt.Println(passwordInfo[2],passwordInfo[3])
	check := password.Verify("admin123", salt,encodedPwd, options)
	fmt.Println(check) // true
	return code
}
func main()  {
	Password("admin")
}
