package api

import (
	"dcomicServer/database"
	"dcomicServer/model"
	"dcomicServer/utils"
	"github.com/gin-gonic/gin"
)

func addUserApi(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/login", login)
		user.POST("/logout", TokenAuth(logout))
		user.GET("/", TokenAuth(getAllUser))
		user.GET("/:username", TokenAuth(getUserById))
		user.POST("/add")
		user.PUT("/:username")
		user.DELETE("/:username")
	}
}

func TokenAuth(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get("token")
		if token == "" {
			context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "token required"})
			context.Abort()
			return
		}
		var user model.User
		err := database.Databases.C("user").Find(map[string]string{"token": token}).One(&user)
		if err != nil || user.Username == "" {
			context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
			context.Abort()
			return
		}
		context.Set("username", user.Username)
		handlerFunc(context)
	}
}

func getAllUser(context *gin.Context) {
	var users []model.User
	err := database.Databases.C("user").Find(map[string]string{}).All(&users)
	if err == nil {
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: users})
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}

func getUserById(context *gin.Context) {
	username := context.Param("username")
	var user model.User
	err := database.Databases.C("user").Find(map[string]string{"username": username}).One(&user)
	if err == nil {
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: user})
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}

// 登录
// @Summary 登录
// @Description 通过用户名密码登录
// @Tags user
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Accept mpfd
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct
// @failure 401 {object} model.StandJsonStruct
// @Router /user/login [post]
func login(context *gin.Context) {
	username, _ := context.GetPostForm("username")
	password, _ := context.GetPostForm("password")
	var user model.User
	err := database.Databases.C("user").Find(map[string]string{"username": username, "password": utils.GetPassword(password)}).One(&user)
	if err == nil {
		token := utils.GetNewToken(24)
		user.Token = token
		updateErr := database.Databases.C("user").Update(map[string]string{"username": username, "password": utils.GetPassword(password)}, user)
		if updateErr == nil {
			context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: map[string]string{"token": token}})
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: updateErr.Error()})
		}
	} else {
		context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "username or password error"})
	}
}

// 登出
// @Summary 登出
// @Description 通过token进行登出
// @security token
// @Tags user
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /user/logout [post]
func logout(context *gin.Context) {
	username, exist := context.Get("username")
	if exist {
		var user model.User
		err := database.Databases.C("user").Find(map[string]interface{}{"username": username}).One(&user)
		if err==nil{
			user.Token=""
			updateErr := database.Databases.C("user").Update(map[string]interface{}{"username": username}, user)
			if updateErr == nil {
				context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: updateErr.Error()})
			}
		}else{
			context.JSON(401,model.StandJsonStruct{Code:401,Msg: "login required"})
		}
	} else {
		context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
	}
}
