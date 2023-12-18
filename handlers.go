package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

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
            RespondWithError(w,r,http.StatusInternalServerError,"Failed to register")
            return
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
            log.Printf("Failed to generate token: %v", err)
            RespondWithError(w,r,http.StatusInternalServerError,"Failed to generate token")
            return
        }

        RespondWithJSON(w,r,http.StatusOK,AuthResponse{"Logged in",jwtTokenString})
    }
}

func UploadHandler (db *gorm.DB) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {
        r.ParseMultipartForm(10<<20)
        file, fileHeader, err := r.FormFile("file")
        if err!= nil {
            log.Println(err)
            RespondWithError(w,r,http.StatusInternalServerError,"Could not upload file")
            return
        }
        defer file.Close()

        cwd,_ := os.Getwd()
        err = os.Mkdir("uploads",0755)
        if err!= nil {
            log.Println(err)
        }
        newFilePath := filepath.Join(cwd,"uploads",fileHeader.Filename)
        newFile, err := os.Create(newFilePath)
        if err!= nil {
            log.Printf("Error creating file: %v",err)
            RespondWithError(w,r,http.StatusInternalServerError,"Could not upload file")
            return
        }
        defer newFile.Close()

        _, err = io.Copy(newFile,file)
        if err!= nil {
            RespondWithError(w,r,http.StatusInternalServerError,"Could not upload file")
            return
        }

        transaction := FileTransaction{
            Filename: fileHeader.Filename,
            Transaction: "UPLOAD",
            Username: r.Context().Value(usernameKey).(string),
            CreatedAt: time.Now(),
        }
        if _,err := CreateTransaction(db,transaction); err!= nil {
            log.Printf("Error updating transaction")
        }

        RespondWithJSON(w,r, http.StatusCreated,
            UploadResponse{
                "Successfully uploaded",
                fileHeader.Filename,
                fileHeader.Size,
                fileHeader.Header.Get("Content-Type") })
    }
}

func DownloadHandler (db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        filename := r.URL.Query().Get("filename")
        if filename == "" {
            RespondWithError(w,r,http.StatusBadRequest,"Missing filename")
            return
        }

        cwd,_ := os.Getwd()
        path := filepath.Join(cwd,"uploads",filepath.Base(filename))

        fileInfo,err := os.Stat(path)
        if err != nil {
            if errors.Is(err,os.ErrNotExist) {
                RespondWithError(w,r,http.StatusNotFound,"File doesn't exist")
            } else {
                RespondWithError(w,r,http.StatusInternalServerError,"Could not download file")
            }
            return
        }

        if fileInfo.IsDir() {
            RespondWithError(w,r,http.StatusNotFound,"File doesn't exist")
            return
        }

        file, err := os.Open(path)
        if err!= nil {
            RespondWithError(w,r,http.StatusInternalServerError,"Could not download file")
            return
        }
        defer file.Close()

        if _,err := io.Copy(w,file); err!=nil {
            RespondWithError(w,r,http.StatusInternalServerError,"Could not download file")
            log.Println(err)
            return
        }

        transaction := FileTransaction{
            Filename: filename,
            Transaction: "DOWNLOAD",
            Username: r.Context().Value(usernameKey).(string),
            CreatedAt: time.Now(),
        }

        if _,err := CreateTransaction(db,transaction); err!= nil {
            log.Printf("Error updating transaction")
        }
    }
}

func ListFilesHandler (w http.ResponseWriter, r *http.Request) {
    cwd,err := os.Getwd()
    if err != nil {
        RespondWithError(w,r,http.StatusInternalServerError,"Could not list files")
    }

    files,_ := os.ReadDir(filepath.Join(cwd,"uploads"))
    
    fileList := make([]FileEntry, 0, len(files))

    for _, file := range files {
        info, err := file.Info()
        if err != nil {
            continue
        }
        fi,_ := os.Open(filepath.Join(cwd,"uploads",file.Name()))
        buf := make([]byte,512)
        fi.Read(buf)
        fi.Close()
        fileType := http.DetectContentType(buf)
        fileList = append(fileList, FileEntry{
            Name: info.Name(),
            Size: info.Size(),
            Type: fileType,
            Modified: info.ModTime(),
        })
    }

    sort.Slice(fileList,func(i, j int) bool {
        return fileList[j].Size < fileList[i].Size
    })
    RespondWithJSON(w,r,http.StatusOK,fileList)
}

func DeleteHandler(db *gorm.DB) http.HandlerFunc {
    return func (w http.ResponseWriter,r *http.Request) {
        filename := r.URL.Query().Get("filename");
        if filename == "" {
            RespondWithError(w,r,http.StatusBadRequest,"No filename provided")
            return
        }
        cwd,err := os.Getwd();
        if err!= nil {
            log.Println(err)
            RespondWithError(w,r,http.StatusInternalServerError,"Failed to delete")
            return
        }
        path := filepath.Join(cwd,"uploads",filename)
        if err := os.Remove(path); err!= nil {
            log.Println(err)
            RespondWithError(w,r,http.StatusInternalServerError,"Failed to delete")
            return
        }
        RespondWithJSON(w,r,http.StatusOK,map[string]string{"message": "Successfully deleted"})
        
        transaction := FileTransaction{
            Filename: filename,
            Transaction: "DELETE",
            Username: r.Context().Value(usernameKey).(string),
            CreatedAt: time.Now(),
        }
        if _,err := CreateTransaction(db,transaction); err!= nil {
            log.Printf("Error updating transaction")
        }
    }
}

func GetTransactionsHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter,r *http.Request) {
        result,err := GetAllTransactions(db)
        if(err!=nil) {
            log.Println(err)
            RespondWithError(w,r,http.StatusInternalServerError,"Could not get file transactions")
            return
        }
        RespondWithJSON(w,r,http.StatusOK,result)
    }
}