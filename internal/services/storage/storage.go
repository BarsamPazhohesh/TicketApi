package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"ticket-api/internal/config"
	"ticket-api/internal/errx"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type StorageService struct {
	Client *minio.Client
}

const (
	TmpPath    = "tickets/temp/"
	TicketPath = "tickets/files/"
)

// NewStorageService creates a new MinIO client and ensures the bucket exists
func NewStorageService(minioClient *minio.Client) *StorageService {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bucket := config.Get().Minio.Bucket
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		if err := minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			log.Fatal(err)
		}
		log.Printf("✅ Bucket %s created successfully", bucket)
	} else {
		log.Printf("✅ Bucket %s already exists", bucket)
	}

	return &StorageService{Client: minioClient}
}

// UploadFileFromReader uploads a file to MinIO from an io.Reader
func (m *StorageService) UploadFileFromReader(ctx context.Context, objectName string, fileReader io.Reader, fileSize int64, contentType string) (string, *errx.APIError) {
	bucket := config.Get().Minio.Bucket

	_, err := m.Client.PutObject(ctx, bucket, objectName, fileReader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", errx.Respond(errx.ErrServiceUnavailable, err)
	}

	return objectName, nil
}

// UploadTicketFileToTemp uploads a ticket file from an HTTP request to the temp path
func (m *StorageService) UploadTicketFileToTemp(ctx context.Context, filename string, fileReader io.Reader, fileSize int64, contentType string) (string, *errx.APIError) {
	objectName := fmt.Sprintf("%s%s", TmpPath, filename)
	return m.UploadFileFromReader(ctx, objectName, fileReader, fileSize, contentType)
}

// GetFile streams a file from MinIO
func (m *StorageService) GetFile(ctx context.Context, objectName string) (*minio.Object, *errx.APIError) {
	bucket := config.Get().Minio.Bucket
	obj, err := m.Client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return nil, errx.Respond(errx.ErrFileNotFound, err)
		}
		return nil, errx.Respond(errx.ErrServiceUnavailable, err)
	}
	return obj, nil
}

// GetPresignedURL generates a presigned URL for downloading
func (m *StorageService) GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, *errx.APIError) {
	reqParams := url.Values{}
	bucket := config.Get().Minio.Bucket

	urlObj, err := m.Client.PresignedGetObject(ctx, bucket, objectName, expires, reqParams)
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return "", errx.Respond(errx.ErrFileNotFound, err)
		}
		return "", errx.Respond(errx.ErrServiceUnavailable, err)
	}

	return urlObj.String(), nil
}

func (m *StorageService) GetPresignedTicketFileURL(ctx context.Context, ticketID string, filename string) (string, *errx.APIError) {
	uid, err := uuid.Parse(ticketID)
	if err != nil {
		return "", errx.Respond(errx.ErrBadRequest, err)
	}
	objectName := fmt.Sprintf("%s%s/%s", TicketPath, uid, filename)
	return m.GetPresignedURL(ctx, objectName, 15*time.Minute)
}

// DeleteFile removes a file from MinIO
func (m *StorageService) DeleteFile(ctx context.Context, objectName string) *errx.APIError {
	bucket := config.Get().Minio.Bucket
	err := m.Client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return errx.Respond(errx.ErrFileNotFound, err)
		}
		return errx.Respond(errx.ErrServiceUnavailable, err)
	}
	return nil
}

// MoveTempsFileToTickets moves specific files from temp to ticket folder
// Returns the list of successfully moved object names
func (m *StorageService) MoveTempsFileToTickets(ctx context.Context, ticketID string, objectNames []string) ([]string, *errx.APIError) {

	// Validate UUID
	uid, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, errx.Respond(errx.ErrBadRequest, err)
	}

	bucket := config.Get().Minio.Bucket
	successful := []string{}

	for _, name := range objectNames {
		tmpName := fmt.Sprintf("%s%s", TmpPath, name)

		// Check if object exists in the temporary path
		_, err := m.Client.StatObject(ctx, bucket, tmpName, minio.StatObjectOptions{})
		if err != nil {
			errResp := minio.ToErrorResponse(err)

			if errResp.Code == "NoSuchKey" {
				// Object not in tmp, check ticket path
				ticketName := fmt.Sprintf("%s%s/%s", TicketPath, uid, name)
				_, ticketErr := m.Client.StatObject(ctx, bucket, ticketName, minio.StatObjectOptions{})
				if ticketErr == nil {
					successful = append(successful, name)
					continue
				}

				// Object not found anywhere
				return successful, errx.Respond(errx.ErrFileNotFound, err)
			}

			// Other internal MinIO error
			return successful, errx.Respond(errx.ErrInternalServerError, err)
		}

		destKey := fmt.Sprintf("%s%s/%s", TicketPath, uid, name)

		// Copy object to ticket folder
		src := minio.CopySrcOptions{Bucket: bucket, Object: tmpName}
		dst := minio.CopyDestOptions{Bucket: bucket, Object: destKey}

		_, err = m.Client.CopyObject(ctx, dst, src)
		if err != nil {
			return successful, errx.Respond(errx.ErrInternalServerError, err)
		}

		// Delete the temp object (best-effort)
		err = m.Client.RemoveObject(ctx, bucket, tmpName, minio.RemoveObjectOptions{})
		if err != nil {
			fmt.Printf("⚠️ failed to delete temp file %s: %v\n", name, err)
		}

		successful = append(successful, name)
	}

	return successful, nil
}
