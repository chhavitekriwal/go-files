package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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

        next.ServeHTTP(w, r)
    })
}
