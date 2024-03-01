package handler

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	dicomFile "github.com/meebish/pocket-health/internal"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func GetDICOMData(c *gin.Context) {
	filename := c.Param("fileName")
	filepath := dicomFile.GenerateLocalFilePath(dicomFile.LocalPath, filename)
	dicomData, err := dicom.ParseFile(filepath, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Could not find the file: %s", filename)})
		return
	}

	tagParam, tagParamExists := c.GetQuery("tag")
	_, getPngParamExists := c.GetQuery("png")

	if !tagParamExists && !getPngParamExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search params, at least one of 'tag' or 'png' query param is required"})
		return
	}

	if tagParamExists {
		if tagParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty tag query param, tag=(xxxx,yyyy) query param is expected"})
			return
		}

		tagValues := strings.Split(strings.Trim(tagParam, "()"), ",")

		if len(tagValues) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid number of values in tag query param, tag=(xxxx,yyyy) query param is expected, got: %s", tagParam)})
			return
		}

		tagGroup, err := strconv.ParseUint(tagValues[0], 16, 16)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag group found, a hexadecimal value is expected, got: %s", tagValues[0])})
			return
		}

		tagElement, err := strconv.ParseUint(tagValues[1], 16, 16)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag element found, a hexadecimal value is expected, got: %s", tagValues[1])})
			return
		}

		newTag := tag.Tag{
			Group:   uint16(tagGroup),
			Element: uint16(tagElement),
		}

		tagInfo, err := tag.Find(newTag)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag %s was requested: %s", tagParam, err.Error())})
			return
		}

		dicomElement, err := dicomData.FindElementByTag(newTag)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No DICOM element was found for the %s tag %s", tagInfo.Name, tagParam)})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"tag-values":             tagParam,
				"header-attribute-name":  tagInfo.Name,
				"header-attribute-value": dicomElement.Value.String(),
			},
		})

	} else if getPngParamExists {
		pixelDataElement, _ := dicomData.FindElementByTag(tag.PixelData)
		pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
		images := []image.Image{}
		for _, fr := range pixelDataInfo.Frames {
			img, _ := fr.GetImage() // The Go image.Image for this frame

			images = append(images, img)
		}

		var buffer bytes.Buffer

		_ = png.Encode(&buffer, images[0])

		c.Data(http.StatusOK, "image/png", buffer.Bytes())
	}
}
