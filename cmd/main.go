package main

import (
	"github.com/gin-gonic/gin"

	"github.com/meebish/pocket-health/api/handler"
	dicomFile "github.com/meebish/pocket-health/internal"
)

func main() {
	router := gin.Default()

	localUploader := &dicomFile.LocalUploader{
		UploadPath: dicomFile.LocalPath,
	}
	router.POST("/upload", func(c *gin.Context) {
		handler.UploadDICOMFile(c, localUploader)
	})

	router.GET("/dicomFile/:fileName", handler.GetDICOMData)

	router.Run("localhost:8080")
}
