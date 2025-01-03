package tokens

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/wetrade/config"
)

type SignedDetails struct {
	Name  string
	Email string
	Role  string
	Id    string
	jwt.StandardClaims
}

func GenerateToken(name string, email string, userId string, role string) (string, error) {
	claims := &SignedDetails{
		Name:  name,
		Email: email,
		Role:  role,
		Id:    userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.LoadConfig().JWT_SECRET))
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.LoadConfig().JWT_SECRET), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "Invalid token claims"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token has expired"
		return
	}

	return claims, msg
}
