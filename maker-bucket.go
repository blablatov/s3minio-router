// File maker minio

package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func makerBucket(chid chan string) error {

	log.SetPrefix("main event: ")
	log.SetFlags(log.Lshortfile)

	pm := parseConfig()

	// Initialize minio client object.
	minioClient, err := minio.New(pm.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(pm.AccessKeyID, pm.SecretAccessKey, ""),
		Secure: pm.UseSSL,
	})
	if err != nil {
		log.Fatalf("Error init minio: %v\n", err)
	}

	pm.Mu.Lock()
	defer pm.Mu.Unlock()
	newBucket := <-chid // uuid of user

	// Create a bucket at region with name like uuid.
	err = minioClient.MakeBucket(context.Background(), newBucket, minio.MakeBucketOptions{Region: pm.Region, ObjectLocking: true})
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
