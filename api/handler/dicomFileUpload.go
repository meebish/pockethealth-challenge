package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func uploadDICOMFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Could not read file - %s", err.Error())})
		return
	}

	id := uuid.New().String()

	savedFileName := fmt.Sprintf("%s-%s", id, file.Filename)
	log.Println("saving file as: ", savedFileName)
	if err := c.SaveUploadedFile(file, savedFileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not save file - %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File saved as: %s", savedFileName)})
}
