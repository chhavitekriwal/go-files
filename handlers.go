package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


func RegisterHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var creds AuthRequest
        var user User

        err := json.NewDecoder(r.Body).Decode(&creds)
        if err != nil || creds.Username == "" || creds.Password == "" {
            RespondWithError(w,r,http.StatusBadRequest,"Invalid request body")
            return
        }

        result := db.Where("username = ?", creds.Username).First(&user)
        if result.Error != gorm.ErrRecordNotFound {
            RespondWithError(w,r,http.StatusBadRequest,"User already exists. Login instead")
            return
        }        
        hashedPass, err := bcrypt.GenerateFromPassword([]byte(creds.Password),bcrypt.DefaultCost)
        if(err!=nil) {
            RespondWithError(w,r,http.StatusInternalServerError,"Failed to hash password")
        }
        newUser := User{Username:creds.Username,PasswordHash: string(hashedPass)}

        if err := db.Create(&newUser).Error; err != nil {
            RespondWithError(w,r,http.StatusInternalServerError,"Failed to register user")
            return
        }

        tokenString, err:= GenerateJWTToken(creds.Username)
        if err != nil {
            log.Fatalf("Failed to generate token: %v", err)
            RespondWithError(w,r,http.StatusInternalServerError,"Failed to generate token")
            return
        } 
        response := AuthResponse{"User registered and logged in",tokenString}
        RespondWithJSON(w,r,http.StatusCreated,response)
    }
}

func LoginHandler(db *gorm.DB) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {
        var creds AuthRequest
        var user User
        
        err := json.NewDecoder(r.Body).Decode(&creds)
        if err != nil || creds.Username == "" || creds.Password == "" {
            RespondWithError(w,r,http.StatusBadRequest,"Invalid request body")
            return
        }

        result := db.Where("username = ?", creds.Username).First(&user)
        if result.Error != nil {
            RespondWithError(w,r,http.StatusUnauthorized,"Invalid username")
            return
        }
        
        err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(creds.Password))
        if err != nil {
            RespondWithError(w,r,http.StatusUnauthorized,"Incorrect password")
            return
        }
        jwtTokenString, err := GenerateJWTToken(creds.Username)
        if err != nil {
            log.Fatalf("Failed to generate token: %v", err)
            RespondWithError(w,r,http.StatusInternalServerError,"Failed to generate token")
            return
        }

        RespondWithJSON(w,r,http.StatusOK,AuthResponse{"Logged in",jwtTokenString})
    }
}
