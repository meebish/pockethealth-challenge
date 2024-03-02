package dicomFile

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

type DICOMHeaderAttributes struct {
	TagValues            string `json:"tag-values"`
	HeaderAttributeName  string `json:"header-attribute-name"`
	HeaderAttributeValue string `json:"header-attribute-value"`
}

// Core logic to strip and find header attributes from tags
func GetDICOMAttribute(tagParam string, dicomData dicom.Dataset) (*DICOMHeaderAttributes, error) {
	// Trims the brackets and separates the tag values from the tagParam input
	tagValues := strings.Split(strings.Trim(tagParam, "()"), ",")

	if len(tagValues) != 2 {
		return nil, fmt.Errorf("invalid number of values in tag query param, tag=(xxxx,yyyy) query param is expected, got: %s", tagParam)
	}

	tagGroup, err := strconv.ParseUint(tagValues[0], 16, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid tag group found, expected hexadecimal values, got: %s", tagValues[0])
	}

	tagElement, err := strconv.ParseUint(tagValues[1], 16, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid tag element found, expected hexadecimal values, got: %s", tagValues[1])
	}

	newTag := tag.Tag{
		Group:   uint16(tagGroup),
		Element: uint16(tagElement),
	}

	tagInfo, err := tag.Find(newTag)
	if err != nil {
		return nil, fmt.Errorf("invalid tag %s was requested: %s", tagParam, err.Error())
	}

	dicomElement, err := dicomData.FindElementByTag(newTag)
	if err != nil {
		return nil, fmt.Errorf("no DICOM element was found for the %s tag %s: %s", tagInfo.Name, tagParam, err.Error())
	}

	return &DICOMHeaderAttributes{
		TagValues:            tagParam,
		HeaderAttributeName:  tagInfo.Name,
		HeaderAttributeValue: dicomElement.Value.String(),
	}, nil
}

// Core logic to only retrieve the pixel data as an image from the DICOM file
func GetDICOMImage(dicomData dicom.Dataset) (*image.Image, error) {
	pixelDataElement, err := dicomData.FindElementByTag(tag.PixelData)
	if err != nil {
		return nil, fmt.Errorf("no DICOM Pixel Data was found: %s", err.Error())
	}

	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	images := []image.Image{}
	for _, fr := range pixelDataInfo.Frames {
		img, err := fr.GetImage()
		if err != nil {
			return nil, err
		}

		images = append(images, img)
	}

	return &images[0], nil
}
