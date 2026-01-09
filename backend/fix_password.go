package main

import (
"database/sql"
"fmt"
"log"

_ "github.com/go-sql-driver/mysql"
"golang.org/x/crypto/bcrypt"
)

func main() {
// 生成admin123的bcrypt hash
password := "admin123"
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
log.Fatal("生成密码失败:", err)
}

fmt.Println("生成的密码hash:", string(hash))

// 连接数据库
db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/app_platform?charset=utf8mb4&parseTime=True&loc=Local")
if err != nil {
log.Fatal("连接数据库失败:", err)
}
defer db.Close()

// 更新密码
_, err = db.Exec("UPDATE admins SET password = ? WHERE username = 'admin'", string(hash))
if err != nil {
log.Fatal("更新密码失败:", err)
}

fmt.Println("密码更新成功！")

// 验证密码
err = bcrypt.CompareHashAndPassword(hash, []byte(password))
if err != nil {
log.Fatal("密码验证失败:", err)
}

fmt.Println("密码验证成功！可以使用 admin/admin123 登录")
}
