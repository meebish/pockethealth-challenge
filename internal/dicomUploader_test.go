package dicomFile_test

import (
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	dicomFile "github.com/meebish/pocket-health/internal"
	"github.com/stretchr/testify/assert"
)

const testFilePath = "../test/IM000001"

// Upload tests
func TestUpload_Success(t *testing.T) {
	// Successful upload test
	// Use the file for upload
	file, err := os.Open(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	fileHeader := &multipart.FileHeader{
		Filename: "IM000001",
		Size:     info.Size(),
	}

	uploader := &dicomFile.LocalUploader{
		UploadPath: "../files",
	}

	// Upload the file
	newFilename, err := uploader.Upload(file, fileHeader)
	assert.NoError(t, err)

	// Validate the file creatin
	newFileLocation := fmt.Sprintf("../files/%s", newFilename)
	_, err = os.Stat(newFileLocation)
	assert.NoError(t, err)

	// Cleanup
	os.Remove(newFileLocation)
}

func TestUpload_CannotCreateDirectory(t *testing.T) {
	// Try to upload file to a directory that can't be created
	// Use the file for upload
	file, err := os.Open(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	fileHeader := &multipart.FileHeader{
		Filename: "IM000001",
		Size:     info.Size(),
	}

	// Use root directory which doesn't have permissions
	uploader := &dicomFile.LocalUploader{
		UploadPath: "/files",
	}

	// Attempt to upload the file
	_, err = uploader.Upload(file, fileHeader)

	// Should get an error saying can't write to root directory
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not create upload folder")
}

func TestUpload_CannotUploadToReadOnlyDir(t *testing.T) {
	// Try to upload file to a directory that can't be uploaded to
	// Create a read only dir
	os.Mkdir("test", 0444)
	defer os.Remove("./test")

	// Use the file for upload
	file, err := os.Open(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	fileHeader := &multipart.FileHeader{
		Filename: "IM000001",
		Size:     info.Size(),
	}

	// Use root directory which doesn't have permissions
	uploader := &dicomFile.LocalUploader{
		UploadPath: "./test",
	}

	// Attempt to upload the file
	_, err = uploader.Upload(file, fileHeader)

	// Should get an error saying can't write to root directory
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not create file destination")
}

func TestUpload_InvalidFileSizeUpload(t *testing.T) {
	// Try to upload file but the size is incorrect
	// Use the file for upload
	file, err := os.Open(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileHeader := &multipart.FileHeader{
		Filename: "IM000001",
		Size:     0,
	}

	// Use root directory which doesn't have permissions
	uploader := &dicomFile.LocalUploader{
		UploadPath: "../files",
	}

	// Attempt to upload the file
	_, err = uploader.Upload(file, fileHeader)

	// Should get an error saying can't write to root directory
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the file was not uploaded properly")
}

// Helper Tests
func TestGenerateDicomFileName(t *testing.T) {
	expectedFilename := "exampleDICOMFile"

	actualFilename := dicomFile.GenerateDicomFileName("|" + expectedFilename)

	// Verify that the generated filename has the correct format (uuid + original filename)
	filenameParts := strings.Split(actualFilename, "|")
	err := uuid.Validate(strings.TrimSuffix(filenameParts[0], "-"))
	assert.NoError(t, err)
	assert.Equal(t, expectedFilename+".dcm", filenameParts[1])
}

func TestGenerateDicomFileName_AlreadyHasSuffix(t *testing.T) {
	expectedFilename := "exampleDICOMFile.dcm"

	actualFilename := dicomFile.GenerateDicomFileName("|" + expectedFilename)

	// Verify that the generated filename has the correct format (uuid + original filename)
	filenameParts := strings.Split(actualFilename, "|")
	err := uuid.Validate(strings.TrimSuffix(filenameParts[0], "-"))
	assert.NoError(t, err)
	assert.Equal(t, expectedFilename, filenameParts[1])
}

func TestGenerateLocalFilePath(t *testing.T) {
	filePath := "./directory"
	filename := "test.dcm"
	expectedPath := fmt.Sprintf("%s/%s", filePath, filename)

	actualPath := dicomFile.GenerateLocalFilePath(filePath, filename)

	assert.Equal(t, expectedPath, actualPath)
}
