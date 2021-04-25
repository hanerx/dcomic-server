package api

import (
	"bytes"
	"dcomicServer/ipfs"
	"dcomicServer/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strings"
)

func addUploadApi(r *gin.Engine) {
	upload := r.Group("/upload")
	{
		upload.POST("/image", uploadImage)
		upload.POST("/manga")
		upload.GET("/ipfs/:cid", catCid)
	}
}

// 上传图片
// @Summary 上传图片
// @Description 上传图片并通过ipfs客户端上传至网络，返回cid
// @Tags upload
// @Param image formData file true "图片"
// @Accept mpfd
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 服务器内部错误
// @failure 400 {object} model.StandJsonStruct 文件格式不正确
// @Router /upload/image [post]
func uploadImage(context *gin.Context) {
	file, err := context.FormFile("image")
	if err != nil {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		return
	} else {
		fileExt := strings.ToLower(path.Ext(file.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" && fileExt != ".webp" {
			context.JSON(400, model.StandJsonStruct{Code: 400, Msg: "wrong file ext"})
			return
		}
		fileContent, err := file.Open()
		cid, err1 := ipfs.Api.Add(fileContent)
		if err == nil && err1 == nil {
			context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: map[string]string{"cid": cid}})
		} else if err != nil {
			context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		} else if err1 != nil {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err1.Error()})
		}
	}
}

// 获取ipfs网络内容
// @Summary 获取ipfs网络内容
// @Description 输入cid，通过cid获取ipfs网络内容
// @Tags upload
// @Param cid path string true "cid"
// @Success 200
// @failure 500 {object} model.StandJsonStruct 服务器内部错误
// @failure 400 {object} model.StandJsonStruct 文件格式不正确
// @Router /upload/ipfs/{cid} [get]
func catCid(context *gin.Context) {
	cid := context.Param("cid")
	reader, err := ipfs.Api.Cat(cid)
	if err == nil {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(reader)
		if err == nil {
			context.Data(200, http.DetectContentType(buf.Bytes()), buf.Bytes())
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
		}
	} else {
		context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
	}

}
