package dicomFile

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
)

const LocalPath = "./files"

// Uploader interface so we can store files other than local
type FileUploader interface {
	Upload(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

// Implementation of FileUploader locally
type LocalUploader struct {
	UploadPath string
}

func (lu *LocalUploader) Upload(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// In case folder doesn't exist
	if err := os.MkdirAll(lu.UploadPath, os.ModePerm); err != nil {
		return "", err
	}

	// Generate file names and paths for local storage
	newFilename := GenerateDicomFileName(fileHeader.Filename)
	newFilepath := GenerateLocalFilePath(lu.UploadPath, newFilename)

	// Create and move file there
	dst, err := os.Create(newFilepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// If something goes wrong, remove the file that was created
	if bytes, err := io.Copy(dst, file); err != nil {
		os.Remove(dst.Name())
		return "", err
	} else if bytes != fileHeader.Size {
		os.Remove(dst.Name())
		return "", fmt.Errorf("the file was not uploaded properly, please try again")
	}

	return newFilename, nil
}

// Helpers

// Generate unique dicom file name in case of repeat file names
func GenerateDicomFileName(filename string) string {
	id := uuid.New().String()

	return fmt.Sprintf("%s-%s.dcm", id, filename)
}

// Generate full loacl dicom file path
func GenerateLocalFilePath(filePath, filename string) string {
	return fmt.Sprintf("%s/%s", filePath, filename)
}
