package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suyashkumar/dicom"

	dicomFile "github.com/meebish/pocket-health/internal"
)

// POST /upload
// UploadDICOMFile validates and uploads the dicom file
func UploadDICOMFile(c *gin.Context, uploader dicomFile.FileUploader) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Could not read file - %s", err.Error())})
		return
	}
	defer file.Close()

	// If the file cannot be parsed, it is likely to be an invalid DICOM file
	if _, err := dicom.Parse(file, header.Size, nil); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File is invalid DICOM file - %s", err.Error())})
		return
	}

	savedFileName, err := uploader.Upload(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not save file - %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File saved as: %s", savedFileName)})
}
