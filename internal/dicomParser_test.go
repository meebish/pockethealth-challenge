package dicomFile_test

import (
	"fmt"
	"testing"

	dicomFile "github.com/meebish/pocket-health/internal"
	"github.com/stretchr/testify/assert"
	"github.com/suyashkumar/dicom"
)

// GetDICOMAttribute Tests
func TestGetDICOMAttribute_Success(t *testing.T) {
	// Actual valid DICOM file parsed, header attribute successfully retrieved
	dataset, err := dicom.ParseFile(testFilePath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dataset)

	expectedDICOMAttr := &dicomFile.DICOMHeaderAttributes{
		TagValues:            "(0008,0080)",
		HeaderAttributeName:  "InstitutionName",
		HeaderAttributeValue: "[SUNNYVALE IMAGING CENTER]",
	}

	actualDICOMAttr, err := dicomFile.GetDICOMAttribute("(0008,0080)", dataset)
	assert.NoError(t, err)
	assert.NotNil(t, actualDICOMAttr)
	assert.Equal(t, expectedDICOMAttr, actualDICOMAttr)
}

func TestGetDICOMAttribute_NotEnoughTagValues(t *testing.T) {
	// Invalid number of tag values received in tags, 1 too few
	dataset, err := dicom.ParseFile(testFilePath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dataset)

	actualDICOMAttr, err := dicomFile.GetDICOMAttribute("(0008)", dataset)
	assert.Error(t, err)
	assert.Nil(t, actualDICOMAttr)
	assert.Contains(t, err.Error(), "invalid number of values in tag query param")
}

func TestGetDICOMAttribute_TooManyTagValues(t *testing.T) {
	// Invalid number of tag values received in tags, 1 too many
	dataset, err := dicom.ParseFile(testFilePath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dataset)

	actualDICOMAttr, err := dicomFile.GetDICOMAttribute("(0008,0080,0800)", dataset)
	assert.Error(t, err)
	assert.Nil(t, actualDICOMAttr)
	assert.Contains(t, err.Error(), "invalid number of values in tag query param")
}

func TestGetDICOMAttribute_InvalidTagGroup(t *testing.T) {
	// Invalid Tag Group in tags
	dataset, err := dicom.ParseFile(testFilePath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dataset)

	actualDICOMAttr, err := dicomFile.GetDICOMAttribute("(000G,0080)", dataset)
	assert.Error(t, err)
	assert.Nil(t, actualDICOMAttr)
	assert.Contains(t, err.Error(), "invalid tag group found")
}

func TestGetDICOMAttribute_InvalidTagElement(t *testing.T) {
	// Invalid Tag Element in tags
	dataset, err := dicom.ParseFile(testFilePath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dataset)

	actualDICOMAttr, err := dicomFile.GetDICOMAttribute("(0008,00G0)", dataset)
	assert.Error(t, err)
	assert.Nil(t, actualDICOMAttr)
	assert.Contains(t, err.Error(), "invalid tag element found")
}

func TestGetDICOMAttribute_NonexistentTag(t *testing.T) {
	// Provided Tag isn't a real tag in the DICOM standard
	dataset, err := dicom.ParseFile(testFilePath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dataset)

	invalidTag := "(FFFF,FFFF)"
	actualDICOMAttr, err := dicomFile.GetDICOMAttribute(invalidTag, dataset)
	assert.Error(t, err)
	assert.Nil(t, actualDICOMAttr)
	assert.Contains(t, err.Error(), fmt.Sprintf("invalid tag %s was requested", invalidTag))
}

func TestGetDICOMAttribute_EmptyElement(t *testing.T) {
	// DICOM file doesn't have an element for the tag
	dataset, err := dicom.ParseFile(testFilePath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dataset)

	actualDICOMAttr, err := dicomFile.GetDICOMAttribute("(0008,1111)", dataset)
	assert.Error(t, err)
	assert.Nil(t, actualDICOMAttr)
	assert.Contains(t, err.Error(), "no DICOM element was found")
}

// GetDICOMImage Tests
func TestGetDICOMImage_Success(t *testing.T) {
	// Actual valid DICOM file parsed, image successfully retrieved
	dataset, err := dicom.ParseFile(testFilePath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dataset)

	img, err := dicomFile.GetDICOMImage(dataset)
	assert.NoError(t, err)
	assert.NotNil(t, img)
}

func TestGetDICOMImage_NoPixelDataElement(t *testing.T) {
	// Empty DICOM dataset without any data in it, no pixel data found
	dataset := dicom.Dataset{}

	// Call the function being tested
	img, err := dicomFile.GetDICOMImage(dataset)
	assert.Error(t, err)
	assert.Nil(t, img)
	assert.Contains(t, err.Error(), "no DICOM Pixel Data was found")
}
