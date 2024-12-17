package config

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

var CloudinaryConfig *cloudinary.Cloudinary

func InitCloudinary() {
	var err error

	// Initialize Cloudinary with credentials from environment variables
	CloudinaryConfig, err = cloudinary.NewFromParams(
		os.Getenv("CLOUD_NAME"),
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
	)

	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
}

func UploadFile(file multipart.File, fileName string, folder string) (string, string, error) {
	// Ensure the CloudinaryConfig is initialized
	if CloudinaryConfig == nil {
		return "", "", errors.New("CloudinaryConfig is not initialized")
	}
	publicID := uuid.New().String()
	// Perform the upload
	uploadParams := uploader.UploadParams{
		PublicID: publicID,
		Folder:   folder,
	}
	uploadResult, err := CloudinaryConfig.Upload.Upload(context.Background(), file, uploadParams)
	if err != nil {
		log.Printf("Error uploading file to Cloudinary: %v", err)
		return "", "", err
	}

	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

func DeleteFile(publicID string) error {
	// parts := strings.Split(publicID, "/")

	// Ensure the CloudinaryConfig is initialized
	if CloudinaryConfig == nil {
		return errors.New("CloudinaryConfig is not initialized")
	}

	// Perform the deletion
	response, err := CloudinaryConfig.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "video",
	})

	if err != nil {
		log.Printf("Error deleting file from Cloudinary: %v", err)
		return err
	}

	log.Printf("Cloudinary Destroy Response: %+v", response)

	if response.Result != "ok" {
		log.Printf("File not deleted. Result: %s", response.Result)
		return errors.New("file not deleted from Cloudinary")
	}

	log.Printf("File with Public ID '%s' successfully deleted from Cloudinary.", publicID)
	return nil
}
