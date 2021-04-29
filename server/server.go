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
	port = ":50051"
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

func (s *mediaServer) ImageDownload(ctx context.Context, req *media_pb.ImageDownloadRequest) (*media_pb.ImageDownloadResponse, error) {
	getLogger("ImageDownload").Info("Processing Image Download")
	imageId := req.GetImageId()
	minioClient, err := media_minio.ConnectMinio()
	if err != nil {
		getLogger("ImageDownload").Error(err.Error())
		return nil, err
	}
	data, err := minioClient.GetObject(context.Background(), media_minio.BUCKET, imageId, minio.GetObjectOptions{})
	defer data.Close()
	if err != nil {
		getLogger("ImageDownload").Error(err.Error())
		return nil, err
	}
	stat, _ := data.Stat()
	dataSize := stat.Size
	imageBuffer := make([]byte, dataSize)
	data.Read(imageBuffer)
	return &media_pb.ImageDownloadResponse{Chunk: imageBuffer}, nil

}

func main() {
	// setup minio
	_, minioError := media_minio.ConnectMinio()
	if minioError != nil {
		getLogger("Server").Fatal(minioError.Error())
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
		getLogger("Server").Fatal(fmt.Sprintf("failed to server : %v", err))
	}
}
