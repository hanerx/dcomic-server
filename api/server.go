package api

import (
	"dcomicServer/database"
	"dcomicServer/model"
	"dcomicServer/utils"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
)

func addServerApi(r *gin.Engine) {
	server := r.Group("/server")
	{
		server.GET("/")
		server.POST("/add", addServer)
		server.DELETE("/delete", deleteServer)
		server.GET("/state", getServerState)
	}
}

// 获取server_state
// @Summary 获取服务器状态
// @Description 获取服务器状态
// @Tags server
// @Produce json
// @Success 200 {object} model.StandJsonStruct{data=model.Node} 正确返回
// @failure 500 {object} model.StandJsonStruct 内部处理错误或已存在
// @Router /server/state [get]
func getServerState(context *gin.Context) {
	ip, err := utils.GetExternalIP()
	if err == nil {
		serverMode, envErr := strconv.Atoi(os.Getenv("server_mode"))
		if envErr != nil {
			serverMode = 0
		}
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: model.Node{Address: ip.String(), Timestamp: time.Now().Unix(), Token: "", Trust: 0, Type: serverMode, Version: "1.0.0", Name: os.Getenv("server_name")}})
	} else {
		context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
	}
}

// 添加server
// @Summary 添加分布式服务器
// @Description 添加一个分布式服务器
// @Tags server
// @Accept json
// @Param server body model.Node true "服务器详情"
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 内部处理错误或已存在
// @failure 400 {object} model.StandJsonStruct body解析失败
// @Router /server/add [post]
func addServer(context *gin.Context) {
	var node model.Node
	err := context.BindJSON(&node)
	if err == nil {
		err = database.Databases.C("server").Find(map[string]string{"address": node.Address}).One(nil)
		if err != nil {
			err = database.Databases.C("server").Insert(node)
			if err == nil {
				context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: "server already exist"})
		}
	} else {
		context.JSON(400, model.StandJsonStruct{Code: 400, Msg: err.Error()})
	}
}

// 删除server
// @Summary 删除分布式服务器
// @Description 删除一个分布式服务器
// @Tags server
// @Param address query string true "服务器地址"
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 内部处理错误或不存在
// @failure 400 {object} model.StandJsonStruct 缺少参数
// @Router /server/delete [delete]
func deleteServer(context *gin.Context) {
	address, exist := context.GetQuery("address")
	if exist {
		err := database.Databases.C("server").Remove(map[string]string{"address": address})
		if err == nil {
			context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
		}
	} else {
		context.JSON(400, model.StandJsonStruct{Code: 400, Msg: "cannot find address"})
	}
}
