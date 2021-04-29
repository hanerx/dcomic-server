package api

import (
	_ "dcomicServer/docs"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	addUtilApi(r)

	addUserApi(r)

	addComicApi(r)

	addUploadApi(r)

	addServerApi(r)

	addTagApi(r)

	return r
}


func Run() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		print(err)
	}
}
