package jwt

import (
    "time"
    "auth-service/internal/models"
    "github.com/golang-jwt/jwt/v4"
    "os"
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

func (j *JWTService) GenerateToken(userID string, email string) (string, error) {
    claims := &models.Claims{
        UserID: userID,
        Email:  email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
            IssuedAt:  time.Now().Unix(),
            Issuer:    "auth-service",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateToken(tokenString string) (*models.Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
        return j.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrInvalidKey
}