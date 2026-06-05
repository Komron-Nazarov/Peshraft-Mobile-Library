// package pkg

// import (
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// type Claims struct {
// 	UserID uint   `json:"user_id"`
// 	Email  string `json:"email"`
// 	Name   string `json:"name"`
// 	jwt.RegisteredClaims
// }

// func GenerateToken(userID uint, email, name, secret string, expiry time.Duration) (string, error) {
// 	claims := &Claims{
// 		UserID: userID,
// 		Email:  email,
// 		Name:   name,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 			Issuer:    "mobile-library",
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(secret))
// }

// func ParseToken(tokenString, secret string) (*Claims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, jwt.ErrSignatureInvalid
// 		}
// 		return []byte(secret), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
// 		return claims, nil
// 	}

// 	return nil, jwt.ErrInvalidKey
// }

package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"` // <--- 1. ДОБАВИЛИ ПОЛЕ В СТРУКТУРУ КЛЕЙМОВ
	jwt.RegisteredClaims
}

// 2. ДОБАВИЛИ ПАРАМЕТР role В СИГНАТУРУ ФУНКЦИИ
func GenerateToken(userID uint, email, name, role, secret string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Name:   name,
		Role:   role, // <--- 3. ИНИЦИАЛИЗИРУЕМ РОЛЬ В ТОКЕНЕ
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "mobile-library",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}