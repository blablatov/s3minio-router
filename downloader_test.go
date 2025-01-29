// File downloader minio testing

package main

import (
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestDown(t *testing.T) {

	var tests = []struct {
		endpoint        string
		accessKeyID     string
		secretAccessKey string
		useSSL          bool
		bucketName      string
		objectName      string
		filePath        string
	}{
		{"storage.yandexcloud.net", "accessKeyID", "secretAccessKey", false,
			"bucketName", "go.gif", "./download/go.gif"},
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

		ctx := context.Background()

		// Initialize minio client object.
		minioClient, err := minio.New(prev_endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(prev_accessKeyID, prev_secretAccessKey, ""),
			Secure: prev_useSSL,
		})
		if err != nil {
			t.Fatal(err)
		}

		// Download the test file with FPutObject
		if err = minioClient.FGetObject(ctx, prev_bucketName, prev_objectName, prev_filePath, minio.GetObjectOptions{}); err != nil {
			t.Fatal(err)
		}

		t.Logf("Successfully downloaded %s\n", prev_objectName)
	}
}

func BenchmarkDown(b *testing.B) {
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
		}{
			{"othercode.ddns.net:9000", "accessKeyID", "secretAccessKey", false,
				"test", "go.gif", "./download/go.gif"},
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

			ctx := context.Background()

			// Initialize minio client object.
			minioClient, err := minio.New(prev_endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(prev_accessKeyID, prev_secretAccessKey, ""),
				Secure: prev_useSSL,
			})
			if err != nil {
				b.Fatal(err)
			}

			// Download the test file with FPutObject
			if err = minioClient.FGetObject(ctx, prev_bucketName, prev_objectName, prev_filePath, minio.GetObjectOptions{}); err != nil {
				b.Fatal(err)
			}

			b.Logf("Successfully downloaded %s\n", prev_objectName)
		}
	}
}
