// File downloader minio
// TODO параметры подключения s3 читать из конфига

package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type DownParams struct {
	mu              sync.Mutex
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	Params          *context.Context
}

func (p *DownParams) Downloader(chs, chid chan string) string {

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

	// Download the test file
	// Change the value of filePath if the file is in another location
	bucketName := <-chid
	objectName := <-chs
	filePath := "./download/" + objectName

	start := time.Now()
	// Download the test file with FPutObject
	if err = minioClient.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{}); err != nil {
		log.Printf("Error: %v", err)
		return ""
	}
	secs := time.Since(start).Seconds()
	log.Printf("%.2fs Time of request\n", secs)

	log.Printf("Successfully downloaded %s\n", objectName)
	return objectName
}
