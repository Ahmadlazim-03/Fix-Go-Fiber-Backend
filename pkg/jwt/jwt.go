package jwt

import (
	"errors"
	"time"

	"Fix-Go-Fiber-Backend/pkg/config"
	"Fix-Go-Fiber-Backend/internal/domain/service"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct {
	secretKey     string
	expireMinutes int
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`     // "mahasiswa", "alumni", or "admin"
	Username string `json:"username"` // for admin
	jwt.RegisteredClaims
}

func NewJWTUtil(cfg *config.Config) *JWTUtil {
	return &JWTUtil{
		secretKey:     cfg.JWT.SecretKey,
		expireMinutes: cfg.JWT.ExpireMinutes,
	}
}

func (j *JWTUtil) GenerateToken(claims *service.JWTClaims) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(j.expireMinutes) * time.Minute)
	
	tokenClaims := &Claims{
		UserID:   claims.UserID,
		Email:    claims.Email,
		Role:     claims.Role,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	return signedToken, expiresAt, err
}

func (j *JWTUtil) ValidateToken(tokenString string) (*service.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &service.JWTClaims{
			UserID:   claims.UserID,
			Email:    claims.Email,
			Role:     claims.Role,
			Username: claims.Username,
		}, nil
	}

	return nil, errors.New("invalid token")
}