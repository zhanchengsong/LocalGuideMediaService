package main

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	media_minio "github.com/zhanchengsong/LocalGuideMediaService/minio"
	media_pb "github.com/zhanchengsong/LocalGuideMediaService/proto"
	"google.golang.org/grpc"
)

const (
	port = ":5005"
)

type mediaServer struct {
	media_pb.UnimplementedImageServer
}

func getLogger(handler string) *log.Entry {
	return log.WithFields(log.Fields{
		"source": handler,
	})
}

func (s *mediaServer) ImageUpload(ctx context.Context, req *media_pb.ImageUploadRequest) (*media_pb.ImageUploadResponse, error) {
	getLogger("ImageUpload").Info("Processing Image Upload")
	byteData := bytes.NewReader(req.GetChunk())
	objectName := uuid.NewString()
	objectSize := req.GetImageSize()
	contentType := req.GetImageType()
	minioClient, err := media_minio.ConnectMinio()
	if err != nil {
		getLogger("ImageUpload").Error(err.Error())
		return &media_pb.ImageUploadResponse{ImageId: "", Url: "", Size: int32(0)}, err
	}
	info, err := minioClient.PutObject(context.Background(), media_minio.BUCKET, objectName, byteData, objectSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		getLogger("ImageUpload").Error(err.Error())
		return &media_pb.ImageUploadResponse{ImageId: "", Url: "", Size: int32(0)}, err
	}
	getLogger("ImageUpload").Info(fmt.Sprintf("upload successful with id = %s", objectName))
	return &media_pb.ImageUploadResponse{ImageId: objectName, Url: "", Size: int32(info.Size)}, err
}

func main() {
	// setup minio
	minioClient, minioError := media_minio.ConnectMinio()
	if minioError != nil {
		getLogger("Server").Fatal(minioError.Error())
	}
	bucketName := "user-avator"
	location := "us-east-1"
	ctx := context.Background()
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	// Create the user avatar bucket
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// set up go rpc server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		getLogger("Server").Error(err.Error())
	}
	s := grpc.NewServer()
	media_pb.RegisterImageServer(s, &mediaServer{})
	getLogger("Server").Info(fmt.Sprintf("Running rpc server at port %s", port))
	if err := s.Serve(lis); err != nil {
		getLogger("Server").Fatal("failed to server : %v", err)
	}
}
