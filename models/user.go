package models

import "OnlineDoc/database"

type User struct {
	UserId   int    `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"user_name" gorm:"size:255"`
	Password string `json:"-" gorm:"size:255"`
	Created  string `json:"created" gorm:"size:255"`
	Role     int    `json:"role" gorm:"size:1"`
	Status   int    `json:"status" gorm:"size:1"`
}

//	func InitializeUser() {
//		db := database.GetDB()
//		err := db.AutoMigrate(&User{})
//		if err != nil {
//			return
//		}
//	}
func (user *User) TableName() string {
	return "t_user"
}

func (user *User) Add() error {
	db := database.GetDB()
	return db.Create(user).Error
}

func (user *User) Update() error {
	db := database.GetDB()
	return db.Save(user).Error
}

func (user *User) Delete() error {
	db := database.GetDB()
	return db.Delete(user).Error
}
func GetUserByUserId(userId int) (User, error) {
	var user User
	db := database.GetDB()
	err := db.Where("user_id = ?", userId).First(&user).Error
	return user, err
}

func GetUserByUsername(username string) (User, error) {
	var user User
	db := database.GetDB()
	err := db.Where("username = ?", username).First(&user).Error
	return user, err
}
