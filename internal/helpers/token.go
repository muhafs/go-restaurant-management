package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/model"
)

func CreateAccessToken(user *entity.User, secret string, expiry int) (accessToken string, err error) {
	claims := model.JwtCustomClaims{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		ID:        user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expiry) * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString([]byte(secret))

	return
}

func CreateRefreshToken(user *entity.User, secret string, expiry int) (refreshToken string, err error) {
	claimsRefresh := &model.JwtCustomRefreshClaims{
		ID: user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(expiry)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	refreshToken, err = token.SignedString([]byte(secret))

	return
}

func IsAuthorized(requestToken string, secret string) (isValid bool, err error) {
	kf := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	}

	if _, err = jwt.Parse(requestToken, kf); err != nil {
		return
	}

	isValid = true
	return
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	kf := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	}

	token, err := jwt.Parse(requestToken, kf)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("Invalid Token")
	}

	return claims["id"].(string), nil
}
