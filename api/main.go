package api

import (
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

	return r
}

// @title DComic API
// @version 1.0.0
// @description  DComic API Doc
// @termsOfService http://github.com/hanerx

// @contact.name GITHUB ISSUE
// @contact.url http://www.github.com/hanerx/dcomic-server/issue

//@host 127.0.0.1:8081
func Run() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		print(err)
	}
}
