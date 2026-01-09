package main

import (
"fmt"
"golang.org/x/crypto/bcrypt"
)

func main() {
password := "admin123"
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
tln("Error:", err)

}
fmt.Println(string(hash))
}
