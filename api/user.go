package api

import (
	"dcomicServer/database"
	"dcomicServer/model"
	"dcomicServer/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func addUserApi(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/login", login)
		user.POST("/logout", TokenAuth(logout, 0))
		user.GET("/my", TokenAuth(getMyInfo, 0))
		user.GET("/", TokenAuth(getAllUser, 1))
		user.GET("/:username", TokenAuth(getUserById, 1))
		user.POST("/:username", TokenAuth(addUser, 1))
		user.PUT("/:username", TokenAuth(updateUser, 0))
		user.DELETE("/:username", TokenAuth(deleteUser, 1))
		subscribe := user.Group("/subscribe")
		subscribe.GET("/", TokenAuth(getSubscribe, 0))
		subscribe.GET("/:comic_id", TokenAuth(getIfSubscribe, 0))
		subscribe.POST("/:comic_id", TokenAuth(addSubscribe, 0))
		subscribe.DELETE("/:comic_id", TokenAuth(cancelSubscribe, 0))
	}
}

func TokenAuth(handlerFunc gin.HandlerFunc, rightNum int) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get("token")
		if token == "" {
			context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "token required"})
			context.Abort()
			return
		}
		var user model.User
		var err error
		if rightNum == 0 {
			err = database.Databases.C("user").Find(map[string]string{"token": token}).One(&user)
		} else {
			err = database.Databases.C("user").Find(map[string]interface{}{"token": token, "rights.right_num": rightNum}).One(&user)
		}
		if err != nil || user.Username == "" {
			context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
			context.Abort()
			return
		}
		context.Set("username", user.Username)
		handlerFunc(context)
	}
}

// 获取所有用户
// @Summary 获取所有用户
// @Description 获取所有用户详情
// @security token
// @Tags user
// @Produce json
// @Success 200 {object} model.StandJsonStruct{data=model.User} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /user/ [get]
func getAllUser(context *gin.Context) {
	var users []model.User
	err := database.Databases.C("user").Find(map[string]string{}).All(&users)
	if err == nil {
		var noPasswordUsers = make([]model.User, len(users))
		for i := 0; i < len(users); i++ {
			noPasswordUsers[i] = users[i].GetUserWithoutPassword()
		}
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: noPasswordUsers})
	} else {
		context.JSON(404, gin.H{"code": 500, "msg": err.Error()})
	}
}

// 获取特定用户
// @Summary 获取特定用户
// @Description 通过username获取用户详情
// @security token
// @Param username path string true "用户名"
// @Tags user
// @Produce json
// @Success 200 {object} model.StandJsonStruct{data=model.User} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /user/{username} [get]
func getUserById(context *gin.Context) {
	username := context.Param("username")
	var user model.User
	err := database.Databases.C("user").Find(map[string]string{"username": username}).One(&user)
	if err == nil {
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: user.GetUserWithoutPassword()})
	} else {
		context.JSON(404, gin.H{"code": 500, "msg": err.Error()})
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
		if err == nil {
			user.Token = ""
			updateErr := database.Databases.C("user").Update(map[string]interface{}{"username": username}, user)
			if updateErr == nil {
				context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: updateErr.Error()})
			}
		} else {
			context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
		}
	} else {
		context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
	}
}

// 添加用户
// @Summary 添加用户
// @Description 添加用户
// @Tags user
// @security token
// @Param username path string true "用户名"
// @Param user body model.User true "用户详情"
// @Accept json
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct
// @failure 400 {object} model.StandJsonStruct
// @Router /user/{username} [post]
func addUser(context *gin.Context) {
	username := context.Param("username")
	err := database.Databases.C("user").Find(map[string]string{"username": username}).One(nil)
	if err != nil {
		var user model.User
		err = context.BindJSON(&user)
		if err == nil {
			user.Password = utils.GetPassword(user.Password)
			user.Username = username
			user.Token = ""
			err = database.Databases.C("user").Insert(user)
			if err == nil {
				context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(400, model.StandJsonStruct{Code: 400, Msg: err.Error()})
		}
	} else {
		context.JSON(500, model.StandJsonStruct{Code: 500, Msg: "user already exist"})
	}
}

// 更新用户
// @Summary 更新用户
// @Description 更新用户
// @Tags user
// @security token
// @Param username path string true "用户名"
// @Param user body model.User true "用户详情"
// @Accept json
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 服务器内部错误
// @failure 400 {object} model.StandJsonStruct body解析错误
// @failure 404 {object} model.StandJsonStruct 用户不存在
// @Router /user/{username} [put]
func updateUser(context *gin.Context) {
	username := context.Param("username")
	var user model.User
	err := database.Databases.C("user").Find(map[string]string{"username": username}).One(&user)
	if err == nil {
		password := user.Password
		err = context.BindJSON(&user)
		if err == nil {
			if user.Password != password && user.Password != "" {
				user.Password = utils.GetPassword(user.Password)
			} else {
				user.Password = password
			}
			user.Username = username
			user.Token = ""
			err = database.Databases.C("user").Update(map[string]string{"username": username}, user)
			if err == nil {
				context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(400, model.StandJsonStruct{Code: 400, Msg: err.Error()})
		}
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "user not found"})
	}
}

// 删除用户
// @Summary 删除用户
// @Description 删除用户
// @Tags user
// @security token
// @Param username path string true "用户名"
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 服务器内部错误
// @failure 404 {object} model.StandJsonStruct 用户不存在
// @Router /user/{username} [delete]
func deleteUser(context *gin.Context) {
	username := context.Param("username")
	err := database.Databases.C("user").Remove(map[string]string{"username": username})
	if err == nil {
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "user not found"})
	}
}

// 获取特定用户
// @Summary 获取自己的用户信息
// @Description 获取登录token的用户信息
// @security token
// @Tags user
// @Produce json
// @Success 200 {object} model.StandJsonStruct{data=model.User} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /user/my [get]
func getMyInfo(context *gin.Context) {
	username, exist := context.Get("username")
	if exist {
		var user model.User
		err := database.Databases.C("user").Find(map[string]interface{}{"username": username}).One(&user)
		if err == nil {
			context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: user.GetUserWithoutPassword()})
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
		}
	} else {
		context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
	}
}

func getSubscribe(context *gin.Context) {
	username, exist := context.Get("username")
	if exist {
		var user model.User
		err := database.Databases.C("user").Find(map[string]interface{}{"username": username}).One(&user)
		if err == nil {
			subscribes, subErr := user.GetSubscribe()
			if subErr == nil {
				context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: subscribes})
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: subErr.Error()})
			}
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
		}
	} else {
		context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
	}
}

func addSubscribe(context *gin.Context) {
	username, exist := context.Get("username")
	if exist {
		var user model.User
		err := database.Databases.C("user").Find(map[string]interface{}{"username": username}).One(&user)
		if err == nil {
			comicId := context.Param("comic_id")
			err = user.AddSubscribe(comicId)
			if err == nil {
				err = database.Databases.C("user").Update(map[string]interface{}{"username": username}, user)
				if err == nil {
					context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
				} else {
					context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
				}
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
		}
	} else {
		context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
	}
}

func cancelSubscribe(context *gin.Context) {
	username, exist := context.Get("username")
	if exist {
		var user model.User
		err := database.Databases.C("user").Find(map[string]interface{}{"username": username}).One(&user)
		if err == nil {
			comicId := context.Param("comic_id")
			err = user.CancelSubscribe(comicId)
			if err == nil {
				err = database.Databases.C("user").Update(map[string]interface{}{"username": username}, user)
				if err == nil {
					context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
				} else {
					context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
				}
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
		}
	} else {
		context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
	}
}

func getIfSubscribe(context *gin.Context) {
	username, exist := context.Get("username")
	if exist {
		comicId := context.Param("comic_id")
		err := database.Databases.C("user").Find(bson.M{"username": username, "subscribes.comic_id": comicId}).One(nil)
		if err == nil {
			context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
		} else {
			context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "not found"})
		}
	} else {
		context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
	}
}
