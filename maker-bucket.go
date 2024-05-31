// File maker minio
// TODO параметры подключения s3 читать из конфига

package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MakeParams struct {
	nameBucket   string
	regionBucket string
}

func (c *MakeParams) MakerBucket(chid chan string) error {

	log.SetPrefix("main event: ")
	log.SetFlags(log.Lshortfile)

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
		log.Fatalf("Error init minio: %v\n", err)
	}

	newbk := <-chid

	// Create a bucket at region 'eu-central-2' with name like uuid.
	err = minioClient.MakeBucket(context.Background(), newbk, minio.MakeBucketOptions{Region: "eu-central-2", ObjectLocking: true})
	if err != nil {
		log.Printf("Error: %v", err)
		errs := err.Error()
		chid <- errs
		return err
	} else {
		log.Println("Successfully created uuid-bucket")
		return nil
	}
	return nil
}
