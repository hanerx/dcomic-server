package api

import (
	"dcomicServer/database"
	"dcomicServer/model"
	"dcomicServer/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func addServerApi(r *gin.Engine) {
	server := r.Group("/server")
	{
		server.GET("/", TokenAuth(getAllServer, 1))
		server.POST("/add", TokenAuth(addServer, 1))
		server.DELETE("/delete", TokenAuth(deleteServer, 1))
		server.GET("/state", getServerState)
		node := server.Group("/node")
		node.POST("/:servername/comic", ServerAuth(nodeUpdateComic))
		node.GET("/:servername/comic", nodeGetComic)
		node.POST("/:servername/user", nodeUpdateComic)
		node.GET("/:servername/user", nodeGetComic)
		node.POST("/:servername/server", nodeUpdateComic)
		node.GET("/:servername/server", nodeGetComic)
	}
}

// 获取server
// @Summary 获取所有已有服务器
// @Description 获取已有服务器
// @security token
// @Tags server
// @Produce json
// @Success 200 {object} model.StandJsonStruct{data=model.Node} 正确返回
// @failure 500 {object} model.StandJsonStruct 内部处理错误或已存在
// @Router /server/ [get]
func getAllServer(context *gin.Context) {
	var nodes []model.Node
	err := database.Databases.C("server").Find(nil).All(&nodes)
	if err == nil {
		if nodes == nil {
			context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: []model.Node{}})
		} else {
			context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: nodes})
		}

	} else {
		context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
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

	serverMode, envErr := strconv.Atoi(os.Getenv("server_mode"))
	if envErr != nil {
		serverMode = 0
	}
	address := os.Getenv("hostname")
	if address == "" {
		ip, err := utils.GetExternalIP()
		if err == nil {
			address = ip.String() + ":" + os.Getenv("port")
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			return
		}
	}
	context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: model.Node{
		Address:     address,
		Timestamp:   time.Now().Unix(),
		Token:       "",
		Trust:       0,
		Type:        serverMode,
		Version:     "1.0.0",
		Name:        os.Getenv("server_name"),
		Description: os.Getenv("server_description"),
		Title:       os.Getenv("server_title"),
	}})
}

// 添加server
// @Summary 添加分布式服务器
// @Description 添加一个分布式服务器
// @security token
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
			response, httpErr := http.Get(fmt.Sprintf("http://%s/server/state", node.Address))
			if httpErr == nil && response.StatusCode == 200 {
				body, readErr := ioutil.ReadAll(response.Body)
				if readErr == nil {
					type Response struct {
						Code int        `json:"code"`
						Msg  string     `json:"msg"`
						Data model.Node `json:"data"`
					}
					var jsonData Response
					err = json.Unmarshal(body, &jsonData)
					if err == nil {
						node.Title = jsonData.Data.Title
						node.Description = jsonData.Data.Description
						node.Version = jsonData.Data.Version
						node.Name = jsonData.Data.Name
						node.Type = jsonData.Data.Type
						err = database.Databases.C("server").Insert(node)
						if err == nil {
							context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
						} else {
							context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
						}
					} else {
						context.JSON(500, model.StandJsonStruct{Code: 500, Msg: "server check error", Data: err})
					}
				} else {
					context.JSON(500, model.StandJsonStruct{Code: 500, Msg: "server check error", Data: readErr})
				}
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: "server check error", Data: httpErr})
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
// @security token
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

func ServerAuth(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get("token")
		if token == "" {
			context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "token required"})
			context.Abort()
			return
		}
		var server model.Node
		err := database.Databases.C("server").Find(map[string]string{"token": token}).One(&server)
		if err != nil || server.Address == "" {
			context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "login required"})
			context.Abort()
			return
		}
		handlerFunc(context)
		//reqIP := context.ClientIP()
		//addr, err := net.ResolveIPAddr("ip", server.Address)
		//if reqIP == server.Address {
		//	handlerFunc(context)
		//} else if err == nil && addr.IP.String() == server.Address {
		//	handlerFunc(context)
		//} else if err != nil {
		//	context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
		//} else {
		//	context.JSON(401, model.StandJsonStruct{Code: 401, Msg: "ip not match"})
		//}
	}
}

// 节点同步
// @Summary 节点同步接口
// @Description 向目标服务器同步内容
// @security server-token
// @Tags server
// @Param address path string true "服务器地址"
// @Param data body []model.ComicDetail true "数据详情"
// @Accept json
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 内部处理错误或不存在
// @failure 400 {object} model.StandJsonStruct 缺少参数
// @Router /server/node/{address}/comic [post]
func nodeUpdateComic(context *gin.Context) {
	var comics []model.ComicDetail
	err := context.BindJSON(&comics)
	if err == nil {
		failed := 0
		success := 0
		ignore := 0
		for i := 0; i < len(comics); i++ {
			var comic model.ComicDetail
			err = database.Databases.C("comic").Find(map[string]string{"comic_id": comics[i].ComicId}).One(&comic)
			if err == nil && comic.Timestamp < comics[i].Timestamp {
				err = database.Databases.C("comic").Update(map[string]string{"comic_id": comics[i].ComicId}, comics[i])
				if err == nil {
					success++
				} else {
					failed++
				}
			} else if err != nil {
				err = database.Databases.C("comic").Insert(comics[i])
				if err == nil {
					success++
				} else {
					failed++
				}
			} else {
				ignore++
			}
		}
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: map[string]int{"success": success, "failed": failed, "ignore": ignore}})
	} else {
		context.JSON(400, model.StandJsonStruct{Code: 400, Msg: err.Error()})
	}
}

// 节点同步
// @Summary 节点同步接口
// @Description 从目标服务器同步内容
// @security server-token
// @Tags server
// @Param address path string true "服务器地址"
// @Produce json
// @Success 200 {object} model.StandJsonStruct{data=model.ComicDetail} 正确返回
// @failure 500 {object} model.StandJsonStruct 内部处理错误或不存在
// @Router /server/node/{address}/comic [get]
func nodeGetComic(context *gin.Context) {
	var comics []model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{}).All(&comics)
	if err == nil {
		data := make([]model.ComicDetail, len(comics))
		for i := 0; i < len(comics); i++ {
			if comics[i].Redirect {
				data[i] = comics[i]
			} else {
				if os.Getenv("hostname") == "" {
					ip, ipErr := utils.GetExternalIP()
					if ipErr == nil {
						data[i] = model.ComicDetail{
							Title:       comics[i].Title,
							Timestamp:   comics[i].Timestamp,
							Description: comics[i].Description,
							Cover:       comics[i].Cover,
							Tags:        comics[i].Tags,
							Authors:     comics[i].Authors,
							ComicId:     comics[i].ComicId,
							HotNum:      comics[i].HotNum,
							Redirect:    true,
							RedirectUrl: ip.String(),
						}
					} else {
						data[i] = comics[i]
					}
				} else {
					data[i] = model.ComicDetail{
						Title:       comics[i].Title,
						Timestamp:   comics[i].Timestamp,
						Description: comics[i].Description,
						Cover:       comics[i].Cover,
						Tags:        comics[i].Tags,
						Authors:     comics[i].Authors,
						ComicId:     comics[i].ComicId,
						HotNum:      comics[i].HotNum,
						Redirect:    true,
						RedirectUrl: os.Getenv("hostname"),
					}
				}
			}
		}
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: data})
	} else {
		context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
	}
}

func nodePutUser(context *gin.Context) {

}

func nodeGetUser(context *gin.Context) {

}

func nodePutServer(context *gin.Context) {

}

func nodeGetServer(context *gin.Context) {

}
