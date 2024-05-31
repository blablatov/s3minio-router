// File uploader minio
// TODO параметры подключения s3 читать из конфига

package main

import (
	"context"
	"log"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UpParams struct {
	mu              sync.Mutex
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
}

func (p *UpParams) Uploader(chup, chid chan string) (string, int64) {

	log.SetPrefix("main event: ")
	log.SetFlags(log.Lshortfile)

	ctx := context.Background()

	endpoint := "storage.yandexcloud.net"
	accessKeyID := "accessKeyID"
	secretAccessKey := "secretAccessKey"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Upload the test file
	// Change the value of filePath if the file is in another location
	bucketName := <-chid
	objectName := <-chup
	filePath := "./upload/" + objectName
	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Printf("Post error: %v", err)
		return "", 0
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return objectName, info.Size
}
