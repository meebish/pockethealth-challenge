package handler

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
	dicomFile "github.com/meebish/pocket-health/internal"
	"github.com/suyashkumar/dicom"
)

type DicomHeaderAttrResp struct {
	DicomHeaderAttributes dicomFile.DICOMHeaderAttributes `json:"data"`
}

// GET /dicomFile/:filename
// GetDICOMData will retrieve either the header attribute or the png of the DICOM file depending on the query parameter
// Query Params:
// tag=(xxxx,yyyy) - Tag to look for in the dicom data, where xxxx is the tag group, and yyyy is the tag element
// png - If it exists, will retrieve png
func GetDICOMData(c *gin.Context) {
	filename := c.Param("fileName")
	filepath := dicomFile.GenerateLocalFilePath(dicomFile.LocalPath, filename)

	dicomData, err := dicom.ParseFile(filepath, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Could not find file: %s. Error: %s", filename, err.Error())})
		return
	}

	tagParam, tagParamExists := c.GetQuery("tag")
	_, getPngParamExists := c.GetQuery("png")

	if !tagParamExists && !getPngParamExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search params, at least one of 'tag' or 'png' query param is required"})
		return
	}

	// Choose what kind of DICOM data to retrieve, will prioritze header attributes if both exists
	if tagParamExists {
		// Retrieve header attributes
		if tagParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty tag query param, tag=(xxxx,yyyy) query param is expected"})
			return
		}

		dicomHeaderAttr, err := dicomFile.GetDICOMAttribute(tagParam, dicomData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Could not retrieve DICOM header attributes: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, DicomHeaderAttrResp{
			DicomHeaderAttributes: *dicomHeaderAttr,
		})

	} else if getPngParamExists {
		// Retrieve png
		dicomImage, err := dicomFile.GetDICOMImage(dicomData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not retrieve image from DICOM file: %s", err.Error())})
			return
		}

		var buffer bytes.Buffer
		if err := png.Encode(&buffer, *dicomImage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not properly encode DICOM as a png: %s", err.Error())})
			return
		}

		c.Data(http.StatusOK, "image/png", buffer.Bytes())
	}
}
