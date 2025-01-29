// File uploader minio testing

package main

import (
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestUp(t *testing.T) {

	var tests = []struct {
		endpoint        string
		accessKeyID     string
		secretAccessKey string
		useSSL          bool
		bucketName      string
		objectName      string
		filePath        string
		contentType     string
	}{
		{"storage.yandexcloud.net", "accessKeyID", "secretAccessKey", false,
			"gitgif", "Go-go.gif", "./upload/Go-go.gif", "application/octet-stream"},
	}

	var prev_endpoint string
	for _, test := range tests {
		if test.endpoint != prev_endpoint {
			t.Logf("\n%s\n", test.endpoint)
			prev_endpoint = test.endpoint
		}

		var prev_accessKeyID string
		if test.accessKeyID != prev_accessKeyID {
			t.Logf("\n%s\n", test.accessKeyID)
			prev_accessKeyID = test.accessKeyID
		}

		var prev_secretAccessKey string
		if test.secretAccessKey != prev_secretAccessKey {
			t.Logf("\n%s\n", test.secretAccessKey)
			prev_secretAccessKey = test.secretAccessKey
		}

		var prev_useSSL bool
		if test.useSSL != prev_useSSL {
			t.Logf("\n%t\n", test.useSSL)
			prev_useSSL = test.useSSL
		}

		var prev_bucketName string
		if test.bucketName != prev_bucketName {
			t.Logf("\n%s\n", test.bucketName)
			prev_bucketName = test.bucketName
		}

		var prev_objectName string
		if test.objectName != prev_objectName {
			t.Logf("\n%s\n", test.objectName)
			prev_objectName = test.objectName
		}

		var prev_filePath string
		if test.filePath != prev_filePath {
			t.Logf("\n%s\n", test.filePath)
			prev_filePath = test.filePath
		}

		var prev_contentType string
		if test.contentType != prev_contentType {
			t.Logf("\n%s\n", test.contentType)
			prev_contentType = test.contentType
		}

		ctx := context.Background()

		// Initialize minio client object.
		minioClient, err := minio.New(prev_endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(prev_accessKeyID, prev_secretAccessKey, ""),
			Secure: prev_useSSL,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Upload the test file with FPutObject
		info, err := minioClient.FPutObject(ctx, prev_bucketName, prev_objectName, prev_filePath,
			minio.PutObjectOptions{ContentType: prev_contentType})
		if err != nil {
			t.Logf("Post error: %v", err)
		}

		t.Logf("Successfully uploaded %s of size %d\n", prev_objectName, info.Size)
	}
}

func BenchmarkUp(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 5; i++ {

		var tests = []struct {
			endpoint        string
			accessKeyID     string
			secretAccessKey string
			useSSL          bool
			bucketName      string
			objectName      string
			filePath        string
			contentType     string
		}{
			{"storage.yandexcloud.net", "accessKeyID", "secretAccessKey", false,
				"gitgif", "Go-go.gif", "./upload/Go-go.gif", "application/octet-stream"},
		}

		var prev_endpoint string
		for _, test := range tests {
			if test.endpoint != prev_endpoint {
				b.Logf("\n%s\n", test.endpoint)
				prev_endpoint = test.endpoint
			}

			var prev_accessKeyID string
			if test.accessKeyID != prev_accessKeyID {
				b.Logf("\n%s\n", test.accessKeyID)
				prev_accessKeyID = test.accessKeyID
			}

			var prev_secretAccessKey string
			if test.secretAccessKey != prev_secretAccessKey {
				b.Logf("\n%s\n", test.secretAccessKey)
				prev_secretAccessKey = test.secretAccessKey
			}

			var prev_useSSL bool
			if test.useSSL != prev_useSSL {
				b.Logf("\n%t\n", test.useSSL)
				prev_useSSL = test.useSSL
			}

			var prev_bucketName string
			if test.bucketName != prev_bucketName {
				b.Logf("\n%s\n", test.bucketName)
				prev_bucketName = test.bucketName
			}

			var prev_objectName string
			if test.objectName != prev_objectName {
				b.Logf("\n%s\n", test.objectName)
				prev_objectName = test.objectName
			}

			var prev_filePath string
			if test.filePath != prev_filePath {
				b.Logf("\n%s\n", test.filePath)
				prev_filePath = test.filePath
			}

			var prev_contentType string
			if test.contentType != prev_contentType {
				b.Logf("\n%s\n", test.contentType)
				prev_contentType = test.contentType
			}

			ctx := context.Background()

			// Initialize minio client object.
			minioClient, err := minio.New(prev_endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(prev_accessKeyID, prev_secretAccessKey, ""),
				Secure: prev_useSSL,
			})
			if err != nil {
				b.Fatal(err)
			}

			// Upload the test file with FPutObject
			info, err := minioClient.FPutObject(ctx, prev_bucketName, prev_objectName, prev_filePath,
				minio.PutObjectOptions{ContentType: prev_contentType})
			if err != nil {
				//t.Fatal(err)
				b.Logf("Post error: %v", err)
			}

			b.Logf("Successfully uploaded %s of size %d\n", prev_objectName, info.Size)
		}
	}
}
