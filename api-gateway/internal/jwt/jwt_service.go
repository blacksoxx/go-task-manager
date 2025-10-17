package jwt

import (
    "os"
    "github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
    secretKey []byte
}

func NewJWTService() *JWTService {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        panic("JWT_SECRET environment variable is required")
    }
    return &JWTService{
        secretKey: []byte(secret),
    }
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return j.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrInvalidKey
}

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.StandardClaims
}