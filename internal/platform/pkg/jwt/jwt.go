package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)

type Service interface {
	GenerateToken(email, userID string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type CustomClaims struct {
	Email  string
	UserID string
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey   string
	issuer      string
	expiredTime time.Duration
}

func NewJWTService(secretKey, issuer string, expiredTime time.Duration) Service {
	return &jwtServices{
		secretKey:   secretKey,
		issuer:      issuer,
		expiredTime: expiredTime,
	}
}

func (j jwtServices) GenerateToken(email, userID string) (string, error) {
	claims := &CustomClaims{
		Email:  email,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.expiredTime).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		encodedToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return []byte(j.secretKey), nil
		})
	if err != nil {
		return nil, err
	}
	if _, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return token, nil
	}
	return nil, ErrInvalidToken
}
