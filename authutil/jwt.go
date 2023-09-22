package authutil

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKeyUser []byte

func InitializeJWTKey() error {
	jwtKeyUser = []byte(os.Getenv("JWT_KEY_USER"))
	if jwtKeyUser == nil {
		return fmt.Errorf("JWT_KEY_USER is empty")
	}
	return nil
}

type JWTClaimAccessUser struct {
	Uid   string   `json:"uid"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateJWTAccessUser(uid string, jti string, perm []string) (string, error) {
	claims := &JWTClaimAccessUser{
		Uid:   uid,
		Roles: perm,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID:       jti,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKeyUser)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ExtractClaimAccessUser(signedToken string) (*JWTClaimAccessUser, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaimAccessUser{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKeyUser, nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaimAccessUser)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}
