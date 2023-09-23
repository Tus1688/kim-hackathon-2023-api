package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const (
	maxOpenConns    = 10
	maxIdleConns    = 5
	connMaxLifetime = 5 * time.Minute
)

var MysqlInstance *sql.DB

func InitMysql() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	MysqlInstance, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	//	validate connection
	if err := MysqlInstance.Ping(); err != nil {
		return err
	}

	//	set connection pool
	MysqlInstance.SetMaxOpenConns(maxOpenConns)
	MysqlInstance.SetMaxIdleConns(maxIdleConns)
	MysqlInstance.SetConnMaxLifetime(connMaxLifetime)

	return nil
}

func InitAdmin() error {
	username := os.Getenv("ADMIN_USERNAME")
	if username == "" {
		username = "admin"
	}
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		return fmt.Errorf("admin password is empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = MysqlInstance.Exec(
		`INSERT INTO users(username, hashed_password, is_admin, kim.users.is_user) VALUES (?, ?, ?, TRUE) ON DUPLICATE KEY UPDATE hashed_password = ?, is_admin = TRUE, is_user = TRUE`,
		username, string(bytes), true, string(bytes),
	)
	if err != nil {
		return err
	}
	return nil
}
