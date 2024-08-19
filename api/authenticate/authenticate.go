package authenticate

import (
	"crypto/rand"
	"encoding/base64"
)

type Session struct {
	Token string
	//Context context.Context
}

type SessionList struct {
	value map[int]Session
}

//var UserSessionList SessionList

func Init() {
	//UserSessionList.value = make(map[int]Session)

}
func GenerateSessionToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

func (userSessionList SessionList) AddSessionToken(userId int, token string) (Session, error) {
	userSessionList.value[userId] = Session{Token: token}
	return userSessionList.value[userId], nil
}
