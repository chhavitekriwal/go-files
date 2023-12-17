package main

import (
	"gorm.io/gorm"
	"time"
    "github.com/golang-jwt/jwt/v4"
)

type User struct {
    gorm.Model
    Username     string `gorm:"unique"`
    PasswordHash string
}

type FileTransaction struct {
    gorm.Model
    UserID   uint		`gorm:"column:user_id"`
    Filename string		`gorm:"column:filename"`
    Type     string 	`gorm:"column:transaction_type"`
    Timestamp time.Time	`gorm:"column:timestamp"`
}

type ErrorResponse struct {
    Statuscode int      `json:"status"`
    Error string        `json:"error"`
}

type AuthRequest struct {
	Username string
	Password string 
}

type AuthResponse struct {
    Message string
    Token string
}

type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

type UploadResponse struct {
    Message string `json:"message"`
    Filename string `json:"filename"`
    Size int64 `json:"size_in_bytes"`
    Type string `json:"type"`
}

type FileEntry struct {
    Name string `json:"filename"`
    Size int64 `json:"size_in_bytes"`
    Type string `json:"type"`
    Modified time.Time `json:"modified_at"`
}