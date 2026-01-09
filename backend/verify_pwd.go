package main
import (
"fmt"
"golang.org/x/crypto/bcrypt"
)
func main() {
hash := "$2a$10$rT1dlqW6k5jR3ldljTLcwe38rrdmgdP5lmdWAwolbMoY8bbZiL4Iu"
err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("admin123"))
if err == nil {
fmt.Println("✓ 密码验证成功！")
} else {
fmt.Println("✗ 密码验证失败:", err)
}
}
