package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JwtService struct {
	SecretKey string
}

func NewJwtService() *JwtService {
	secret := viper.GetString("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_change_me"
	}
	return &JwtService{
		SecretKey: secret,
	}
}

func (s *JwtService) Sign(claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Copy claims
	tokenClaims := token.Claims.(jwt.MapClaims)
	for k, v := range claims {
		tokenClaims[k] = v
	}

	// Add expiry if not present
	if _, ok := tokenClaims["exp"]; !ok {
		tokenClaims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	}

	return token.SignedString([]byte(s.SecretKey))
}

func (s *JwtService) Verify(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
