package helpers

import (
	"errors"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/yosa12978/halo/internal/pkg/dto"
	"github.com/yosa12978/halo/internal/pkg/models"
)

func GetJwtSecret() string {
	return os.Getenv("JWT_SECRET")
}

func GetJwtToken(user *models.BaseAccount) (string, error) {
	signingKey := []byte(GetJwtSecret())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"roles":   user.Roles,
	})
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func VerifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte(GetJwtSecret())
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, nil
}

func ParseToken(token_header string) (jwt.Claims, error) {
	if len(token_header) == 0 {
		return nil, errors.New("authorization header is empty")
	}
	token_string := strings.Replace(token_header, "Bearer ", "", 1)
	claims, err := VerifyToken(token_string)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func GetUser(token_string string) (dto.UserInfo, error) {
	claims, err := ParseToken(token_string)
	if err != nil {
		return dto.UserInfo{}, err
	}
	user_id := claims.(jwt.MapClaims)["user_id"].(string)
	roles := claims.(jwt.MapClaims)["roles"].([]string)
	user_info := dto.UserInfo{User_id: user_id, Roles: roles}
	return user_info, nil
}
