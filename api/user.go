package api

import "github.com/gin-gonic/gin"

type User struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	password string
	Token    string `json:"token"`
}

func addUserApi(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/login")
		user.GET("/")
	}
}
