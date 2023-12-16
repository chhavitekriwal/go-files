package main

import (
    "log"
    "os"
)

func GetJWTSecretKey() []byte {
    secret := os.Getenv("JWT_SECRET_KEY")
    if secret == "" {
        log.Fatal("JWT secret key is not set")
    }
    return []byte(secret)
}
