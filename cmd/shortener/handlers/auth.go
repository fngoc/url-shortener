package handlers

import (
	"errors"
	"fmt"
	"github.com/fngoc/url-shortener/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"math/rand"
	"net/http"
	"time"
)

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString() (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		// собственное утверждение
		UserID: rand.Int(),
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUserID(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		return -1
	}
	return claims.UserID
}

// AuthMiddleware — middleware для аунтификации HTTP-запросов.
func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if errors.Is(err, http.ErrNoCookie) {
			token, err := BuildJWTString()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    token,
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: true,
				Secure:   true,
			})
			logger.Log.Info(fmt.Sprintf("Create new cookie with token: %s", token))
		} else {
			userId := GetUserID(cookie.Value)
			if userId == -1 {
				w.WriteHeader(http.StatusUnauthorized)
				logger.Log.Info(fmt.Sprintf("Token is not valid"))
				return
			}
			logger.Log.Info(fmt.Sprintf("Auth success with UserId: %d", userId))
		}

		h.ServeHTTP(w, r)
	}
}
