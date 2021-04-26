package media_minio

import (
	"fmt"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

var (
	HOST              = "localhost"
	PORT              = "9000"
	BUCKET            = "user-avator"
	ACCESS_KEY_ID     = "localTestKey"
	ACCESS_SECRET_KEY = "localSecretKey"
	USE_SSL           = false
	LOCATION          = "us-east-1"
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
	})
	return minioClientInstance, minioClientError
}
