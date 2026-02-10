package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	pass := "password123"

	if len(os.Args) > 1 {
		pass = os.Args[1]
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(hash))
}
