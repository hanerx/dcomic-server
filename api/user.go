package api

import (
	"dcomicServer/database"
	"github.com/gin-gonic/gin"
)

type UserRight struct {
	RightNum         int
	RightTarget      interface{}
	RightDescription string
}

type User struct {
	Nickname string `json:"nickname" bson:"nickname"`
	Username string `json:"username" bson:"username"`
	Avatar   string `json:"avatar" bson:"avatar"`
	Password string `json:"-" bson:"password"`
	Token    string `json:"token" bson:"token"`
	Rights   []UserRight
}

func addUserApi(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/login")
		user.POST("/logout")
		admin := user.Group("/admin")
		admin.Use(TokenAuth())
		admin.GET("/", getAllUser)
		admin.POST("/:username", getUserById)
	}
}

func TokenAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get("token")
		if token == "" {
			context.JSON(401, StandJsonStruct{Code: 401, Msg: "token required"})
			context.Abort()
			return
		}
		var user User
		err := database.Databases.C("user").Find(map[string]string{"token": token}).One(&user)
		if err != nil || user.Username == "" {
			context.JSON(401, StandJsonStruct{Code: 401, Msg: "login required"})
			context.Abort()
			return
		}
	}
}

func getAllUser(context *gin.Context) {
	var users []User
	err := database.Databases.C("user").Find(map[string]string{}).All(&users)
	if err == nil {
		context.JSON(200, StandJsonStruct{Code: 200, Msg: "success", Data: users})
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}

func getUserById(context *gin.Context) {
	username := context.Param("username")
	var user User
	err := database.Databases.C("user").Find(map[string]string{"username": username}).One(&user)
	if err == nil {
		context.JSON(200, StandJsonStruct{Code: 200, Msg: "success", Data: user})
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}
