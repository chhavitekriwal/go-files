package main

import (
	"errors"
    "log"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func GetDB() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("file_management.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database open error: %v",err)
		return nil, err
	}
	return db, nil
}

func MigrateModels(db *gorm.DB) error {
	user_mig_err := db.AutoMigrate(&User{})
	if user_mig_err != nil {
		log.Fatalf("Users table automigrate error: %v",user_mig_err)
		return user_mig_err
	}

	transaction_mig_err := db.AutoMigrate(&FileTransaction{})
	if transaction_mig_err != nil {
		log.Fatalf("File transactions table automigrate error: %v",transaction_mig_err)
		return transaction_mig_err
	}

	return nil
}

func CreateTransaction(db *gorm.DB,transaction FileTransaction) (int64,error) {
    result := db.Create(&transaction)
    if(result.RowsAffected == 0) {
        return 0,errors.New("transaction not added")
    }
    return result.RowsAffected,nil
}