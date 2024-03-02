package dicomFile

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

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
		return "", fmt.Errorf("could not create upload folder: %s", err.Error())
	}

	// Generate file names and paths for local storage
	newFilename := GenerateDicomFileName(fileHeader.Filename)
	newFilepath := GenerateLocalFilePath(lu.UploadPath, newFilename)

	// Create and move file there
	dst, err := os.Create(newFilepath)
	if err != nil {
		return "", fmt.Errorf("could not create file destination: %s", err.Error())
	}
	defer dst.Close()

	// If something goes wrong, remove the file that was created
	if bytes, err := io.Copy(dst, file); err != nil {
		if remErr := os.Remove(dst.Name()); remErr != nil {
			return "", fmt.Errorf("error cleaning the file after something went wrong during upload: %s. Orig err: %s", remErr.Error(), err.Error())
		}
		return "", fmt.Errorf("the file could not be uploaded, please try again: %s", err.Error())
	} else if bytes != fileHeader.Size {
		if remErr := os.Remove(dst.Name()); remErr != nil {
			return "", fmt.Errorf("error cleaning the file after uploading incorrect size: %s. Orig err: %s", remErr.Error(), err.Error())
		}
		return "", fmt.Errorf("the file was not uploaded properly, please try again")
	}

	return newFilename, nil
}

// Helpers
// Generate unique DICOM file name in case of repeat file names
func GenerateDicomFileName(filename string) string {
	id := uuid.New().String()

	// Check if the filename already ends with ".dcm"
	if !strings.HasSuffix(filename, ".dcm") {
		filename += ".dcm"
	}

	return fmt.Sprintf("%s-%s", id, filename)
}

// Generate full loacl DICOM file path
func GenerateLocalFilePath(filePath, filename string) string {
	return fmt.Sprintf("%s/%s", filePath, filename)
}
