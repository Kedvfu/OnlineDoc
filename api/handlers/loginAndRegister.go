package handlers

import (
	"OnlineDoc/api/authenticate"
	"OnlineDoc/models"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func HandleLogin(context *gin.Context) {
	username := context.PostForm("username")
	password := context.PostForm("password")

	trueUser, err := models.GetUserByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(200, gin.H{
			"message": "No user found",
		})
		return
	}
	truePasswordBytes := []byte(trueUser.Password)
	err = bcrypt.CompareHashAndPassword(truePasswordBytes, []byte(password))
	if err == nil {
		token, err := authenticate.GenerateSessionToken()
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Unable to generate session token",
			})
			return
		}
		RedisClient.Set(token, strconv.Itoa(trueUser.UserId), time.Second*3600)
		//log.Println("Token - UserId generated: ", RedisClient.Get(token).Val())
		context.SetCookie("session_token", token, 3600, "/", "127.0.0.1:8080", false, true)
		context.JSON(200, gin.H{
			"message":       "Login successful",
			"session_token": token,
			"user_id":       trueUser.UserId,
		})
	} else {
		context.JSON(200, gin.H{
			"message": "Invalid password",
		})
	}
}

func HandleRegister(context *gin.Context) {
	username := context.PostForm("username")
	password := context.PostForm("password")
	_, err := models.GetUserByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Invalid password",
			})
			return
		}
		user := models.User{
			Username: username,
			Password: string(bytes),
		}
		err = user.Add()
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Unable to create user",
			})
			return
		} else {
			context.JSON(200, gin.H{
				"message": "User created successfully",
				"success": 1,
			})
		}
	} else {
		context.JSON(200, gin.H{
			"message": "Username already exists",
		})
	}
}
