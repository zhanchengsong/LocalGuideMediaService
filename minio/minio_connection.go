package media_minio

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

var (
	HOST              = os.Getenv("MINIO_HOST")
	PORT              = os.Getenv("MINIO_PORT")
	BUCKET            = os.Getenv("MINIO_BUCKET_NAME")
	ACCESS_KEY_ID     = os.Getenv("MINIO_ACCESS_KEY")
	ACCESS_SECRET_KEY = os.Getenv("MINIO_SECRET_KEY")
	USE_SSL           = false
	LOCATION          = os.Getenv("MINIO_LOCATION")
)

var minioOnce sync.Once

var minioClientInstance *minio.Client
var minioClientError error

func getLogger() *log.Entry {
	return log.WithFields(log.Fields{
		"source": "monio connection",
	})
}

func ConnectMinio() (*minio.Client, error) {
	minioOnce.Do(func() {
		var endpoint = fmt.Sprintf("%s:%s", HOST, PORT)
		getLogger().Info(fmt.Sprintf("Connecting to minio at %s", endpoint))
		minioClientInstance, minioClientError = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(ACCESS_KEY_ID, ACCESS_SECRET_KEY, ""),
			Secure: USE_SSL,
		})

		if minioClientError != nil {
			getLogger().Error(minioClientError.Error())
			return
		}
		bucketName := BUCKET
		location := LOCATION
		ctx := context.Background()
		err := minioClientInstance.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		// Create the user avatar bucket
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			exists, errBucketExists := minioClientInstance.BucketExists(ctx, bucketName)
			if errBucketExists == nil && exists {
				getLogger().Info("Already owned " + bucketName)
			} else {
				getLogger().Fatal(err.Error())
			}
		} else {
			getLogger().Info("Successfully created " + bucketName)
		}

	})
	return minioClientInstance, minioClientError
}
