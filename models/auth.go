package models

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Tus1688/kim-hackathon-2023-api/authutil"
	"github.com/Tus1688/kim-hackathon-2023-api/database"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

type compareUser struct {
	id             string
	hashedPassword string
	isAdmin        bool
}

type internalRefresh struct {
	Uid   string   `json:"uid"`
	Jti   string   `json:"jti"`
	Roles []string `json:"roles"`
}

type CreateUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserResponse struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	IsAdmin   bool   `json:"is_admin"`
	UpdatedOn string `json:"updated_on"`
}

type ModifyUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func (c *compareUser) SerializeRoles() []string {
	var roles []string
	if c.isAdmin {
		roles = append(roles, "admin")
	}
	return roles
}

func (l *LoginRequest) Login() (string, string, error, LoginResponse) {
	var row compareUser
	err := database.MysqlInstance.QueryRow(
		`SELECT BIN_TO_UUID(id), hashed_password, is_admin FROM users WHERE username = ?`, l.Username,
	).Scan(&row.id, &row.hashedPassword, &row.isAdmin)
	if err != nil {
		time.Sleep(55 * time.Millisecond)
		return "", "", fmt.Errorf("invalid username or password"), LoginResponse{}
	}

	err = bcrypt.CompareHashAndPassword([]byte(row.hashedPassword), []byte(l.Password))
	if err != nil {
		time.Sleep(55 * time.Millisecond)
		return "", "", fmt.Errorf("invalid username or password"), LoginResponse{}
	}

	roles := row.SerializeRoles()
	jti := authutil.GenerateRandomString(5)
	accessToken, err := authutil.GenerateJWTAccessUser(row.id, jti, roles)
	if err != nil {
		return "", "", err, LoginResponse{}
	}
	refreshValue := internalRefresh{
		Uid:   row.id,
		Jti:   jti,
		Roles: roles,
	}
	jsonString, err := json.Marshal(refreshValue)
	if err != nil {
		return "", "", err, LoginResponse{}
	}
	refreshToken := authutil.GenerateRandomString(32)
	err = database.RedisInstance[0].Set(context.Background(), refreshToken, jsonString, 24*time.Hour).Err()
	if err != nil {
		return "", "", err, LoginResponse{}
	}

	return accessToken, refreshToken, nil, LoginResponse{
		Username: l.Username,
		IsAdmin:  row.isAdmin,
	}
}

func GetRefreshToken(refreshToken string) (string, error) {
	res, err := database.RedisInstance[0].Get(context.Background(), refreshToken).Result()
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	var redisValue internalRefresh
	err = json.Unmarshal([]byte(res), &redisValue)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	// generate access token
	accessToken, err := authutil.GenerateJWTAccessUser(
		redisValue.Uid, redisValue.Jti, redisValue.Roles,
	)
	if err != nil {
		return "", fmt.Errorf("cannot generate access token")
	}

	return accessToken, nil
}

func (c *CreateUser) Create() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = database.MysqlInstance.Exec(
		"INSERT INTO users (username, hashed_password, is_admin) VALUES (?, ?, ?)",
		c.Username, string(bytes), c.IsAdmin,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]UserResponse, error) {
	rows, err := database.MysqlInstance.Query("SELECT BIN_TO_UUID(id), username, is_admin, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []UserResponse
	for rows.Next() {
		var temp UserResponse
		err := rows.Scan(&temp.Id, &temp.Username, &temp.IsAdmin, &temp.UpdatedOn)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}
	return res, nil
}

func DeleteUser(id string) error {
	res, err := database.MysqlInstance.Exec(
		`DELETE FROM users WHERE id = UUID_TO_BIN(?)`, id,
	)
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); err != nil || affected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (m *ModifyUser) Modify() error {
	if m.Username == os.Getenv("ADMIN_USERNAME") || m.Username == "admin" {
		return fmt.Errorf("cannot modify admin account")
	}
	query := "UPDATE users SET updated_at = NOW()"
	var args []interface{}
	if m.Password != "" {
		bytes, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		query += ", hashed_password = ?"
		args = append(args, string(bytes))
	}
	if m.IsAdmin {
		query += ", is_admin = TRUE"
	} else {
		query += ", is_admin = FALSE"
	}
	query += " WHERE username = ?"
	args = append(args, m.Username)

	res, err := database.MysqlInstance.Exec(query, args...)
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); err != nil || affected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
