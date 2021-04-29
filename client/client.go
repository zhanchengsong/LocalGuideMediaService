package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	media_pb "github.com/zhanchengsong/LocalGuideMediaService/proto"
	"google.golang.org/grpc"
)

func uploadImage(client media_pb.ImageClient, filePath string) string {
	log.Printf("Client upload image")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	imageData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	req := media_pb.ImageUploadRequest{ImageName: "testImage", ImageType: "image/png", ImageSize: int64(len(imageData)), Chunk: imageData}
	uploadResult, err := client.ImageUpload(ctx, &req)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print(fmt.Sprintf("Uploaded %d bytes of data", uploadResult.GetSize()))
	return uploadResult.ImageId

}

func downloadImage(client media_pb.ImageClient, imageId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := media_pb.ImageDownloadRequest{ImageId: imageId}
	downloadResult, err := client.ImageDownload(ctx, &req)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print(fmt.Sprintf("Downloaded %v bytes of data", len(downloadResult.GetChunk())))
}

func main() {
	serverAddr := "localhost:50051"
	log.Print(serverAddr)
	//creds, _ := credentials.NewClientTLSFromFile("tls.crt", "")
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	log.Print("Connected to rpc server")
	defer conn.Close()
	client := media_pb.NewImageClient(conn)

	//Test
	pwd, _ := os.Getwd()
	imageId := uploadImage(client, pwd+"/test.png")
	downloadImage(client, imageId)

}
