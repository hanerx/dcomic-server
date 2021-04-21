package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
	"time"
)

func addUploadApi(r *gin.Engine) {
	upload := r.Group("/upload")
	{
		upload.POST("/image", uploadImage)
		upload.POST("/manga")
	}
}

func uploadImage(context *gin.Context) {
	file,err:=context.FormFile("image")
	if err!=nil{
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		return
	}else{
		fileExt:=strings.ToLower(path.Ext(file.Filename))
		if fileExt!=".png"&&fileExt!=".jpg"&&fileExt!=".gif"&&fileExt!=".jpeg"{
			context.JSON(400, gin.H{
				"code": 400,
				"msg":  "上传失败!只允许png,jpg,gif,jpeg文件",
			})
			return
		}
		fileName:=fmt.Sprintf("%s%s",time.Now().Format("2006-01-02"),file.Filename)
		filepath:=fmt.Sprintf("./uploads/%s%s",fileName,fileExt)
		err:=context.SaveUploadedFile(file, filepath)
		if err==nil{
			context.JSON(200, gin.H{
				"code": 200,
				"msg":  "上传成功!",
				"result":gin.H{
					"path":filepath,
				},
			})
		}else{
			context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		}
	}
}
