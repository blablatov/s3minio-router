// File uploader minio

package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func uploader(chup, chid chan string) (string, int64) {

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

	// Upload the test file
	// Change the value of filePath if the file is in another location
	bucketName := <-chid
	objectName := <-chup
	filePath := `./` + pm.UploaDir + `/` + objectName
	// contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: pm.ContentType})
	if err != nil {
		log.Printf("Post error: %v", err)
		return "", 0
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return objectName, info.Size
}
