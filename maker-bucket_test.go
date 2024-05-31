//File maker minio testing

package main

import (
	"context"
	"log"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestMaker(t *testing.T) {

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
		t.Fatalf("Error init minio: %v\n", err)
	}

	// Remove a bucket to setted name
	err = minioClient.RemoveBucket(context.Background(), "newbucket")
	if err != nil {
		t.Logf("Error remove bucket: %v\n", err)
	} else {
		t.Log("Successfully remove bucket to setted name")
	}

	// Create a bucket at region 'us-east-1' with object locking enabled.
	err = minioClient.MakeBucket(context.Background(), "newbucket", minio.MakeBucketOptions{Region: "eu-central-2", ObjectLocking: true})
	if err != nil {
		t.Logf("Error: %v", err)
	}
	t.Log("Successfully created mybucket")
}
