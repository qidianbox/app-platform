package main
import (
"database/sql"
"fmt"
_ "github.com/go-sql-driver/mysql"
"golang.org/x/crypto/bcrypt"
)
func main() {
hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
fmt.Println("新密码hash:", string(hash))
db, _ := sql.Open("mysql", "root:root@tcp(localhost:3306)/app_platform")
defer db.Close()
db.Exec("DELETE FROM admins WHERE username='admin'")
db.Exec("INSERT INTO admins (username, password, nickname) VALUES ('admin', ?, 'Administrator')", string(hash))
fmt.Println("密码已更新，请使用 admin/admin123 登录")
}
