package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func RespondWithJSON(w http.ResponseWriter, r *http.Request,statuscode int,response interface{}) {
    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(statuscode)
    json.NewEncoder(w).Encode(response)
}

func RespondWithError(w http.ResponseWriter,r *http.Request,statuscode int, message string) {
    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(statuscode)
    json.NewEncoder(w).Encode(ErrorResponse{Statuscode: statuscode,Error: message})
}

func GenerateJWTToken(username string) (string,error){
    expirationTime := time.Now().Add(30 * time.Minute) 
    claims := &Claims{
            Username: username,
            RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)

    return tokenString,err
}

