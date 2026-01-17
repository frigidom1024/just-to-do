package auth

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type TokenTool interface {
	GenerateToken(userID int64, username, role string) (string, error)

	ParseToken(token string) (*CustomClaims, error)

	RefreshToken(token string) (string, error)
}

type jwtToken struct {
	secretKey      []byte
	expireDuration time.Duration
}

var jwtTokenInstance TokenTool
var jwtTokenOnce sync.Once

func GetTokenTool() TokenTool {

	jwtTokenOnce.Do(func() {

		jwtTokenInstance = &jwtToken{

			secretKey: []byte("secret"),

			expireDuration: time.Hour * 24,
		}

	})

	return jwtTokenInstance

}

func (j *jwtToken) GenerateToken(userID int64, username, role string) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:   userID,
		Username: username,
		Role:     role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *jwtToken) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func (j *jwtToken) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return j.GenerateToken(claims.UserID, claims.Username, claims.Role)
}
