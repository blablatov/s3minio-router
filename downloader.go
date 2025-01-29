// File downloader minio

package main

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func downloader(chs, chid chan string) string {

	log.SetPrefix("main event: ")
	log.SetFlags(log.Lshortfile)

	pm := parseConfig()

	ctx := context.Background()

	// Initialize minio client object.
	minioClient, err := minio.New(pm.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(pm.AccessKeyID, pm.SecretAccessKey, ""),
		Secure: pm.UseSSL,
	})
	if err != nil {
		log.Fatalf("Error minio client: %v", err)
	}

	// Download the test file
	// Change the value of filePath if the file is in another location
	bucketName := <-chid
	objectName := <-chs
	filePath := `./` + pm.DownloaDir + `/` + objectName

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
