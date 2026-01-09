package main

import (
"fmt"
"golang.org/x/crypto/bcrypt"
)

func main() {
hash := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"

passwords := []string{"password", "admin123", "admin", "123456"}

for _, pwd := range passwords {
err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
if err == nil {
fmt.Printf("✓ 密码 '%s' 匹配成功！\n", pwd)
} else {
fmt.Printf("✗ 密码 '%s' 不匹配\n", pwd)
}
}
}
