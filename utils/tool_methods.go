package utils

import (
	"crypto/md5"
	"dcomicServer/database"
	"dcomicServer/model"
	"encoding/hex"
	"math/rand"
	"time"
)

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetNewToken(l int) string {
	token := GetRandomString(l)
	var user model.User
	err := database.Databases.C("user").Find(map[string]string{"token": token}).One(&user)
	if err == nil {
		return GetNewToken(l)
	}
	return token
}

func GetPassword(password string) string  {
	hash:=md5.New()
	hash.Write([]byte("dcomic_salt"))
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
