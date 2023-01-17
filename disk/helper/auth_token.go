package helper

import (
	"cloud-disk/disk/internal/config"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// token 解析
func AuthToken(token string) (*UserClaim, error) {
	user := new(UserClaim)
	auth, err := jwt.ParseWithClaims(token, user, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !auth.Valid {
		return user, errors.New("token is invalid")
	}
	return user, nil
}
