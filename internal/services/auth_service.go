package services

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "rekber/internal/config"
)

type AuthService interface {
    GenerateToken(userID uint, role string) (string, error)
    ValidateToken(tokenString string) (uint, string, error)
}

type authService struct {
    jwtSecret string
    jwtExpire time.Duration
}

func NewAuthService() AuthService {
    // Safety check: ensure AppConfig is loaded
    if config.AppConfig == nil {
        // Return default values if config not loaded
        // This shouldn't happen in production
        return &authService{
            jwtSecret: "default-secret-for-testing",
            jwtExpire: 72 * time.Hour,
        }
    }
    
    hours := config.AppConfig.JWTExpire
    return &authService{
        jwtSecret: config.AppConfig.JWTSecret,
        jwtExpire: time.Duration(hours) * time.Hour,
    }
}

type Claims struct {
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func (s *authService) GenerateToken(userID uint, role string) (string, error) {
    expirationTime := time.Now().Add(s.jwtExpire)

    claims := &Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "rekber",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ValidateToken(tokenString string) (uint, string, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Validate signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(s.jwtSecret), nil
    })

    if err != nil {
        return 0, "", err
    }

    if !token.Valid {
        return 0, "", errors.New("invalid token")
    }

    return claims.UserID, claims.Role, nil
}