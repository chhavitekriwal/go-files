package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string
const usernameKey contextKey = "username"

func AuthenticationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return GetJWTSecretKey(), nil
        })

        if err != nil || !token.Valid {
			if err!=nil {
				log.Println(err)
				} else { 
					log.Println("Invalid token")
				}
            RespondWithError(w,r,http.StatusUnauthorized,"Unauthorized")
            return
        }

        ctx := context.WithValue(r.Context(), usernameKey, claims.Username)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
