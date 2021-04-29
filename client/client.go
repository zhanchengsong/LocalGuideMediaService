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
	"google.golang.org/grpc/credentials"
)

func uploadImage(client media_pb.ImageClient, filePath string) {
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

}

func main() {
	serverAddr := "media.zhancheng.dev:443"
	log.Print(serverAddr)
	creds, _ := credentials.NewClientTLSFromFile("tls.crt", "")
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	log.Print("Connected to rpc server")
	defer conn.Close()
	client := media_pb.NewImageClient(conn)

	//Test
	pwd, _ := os.Getwd()
	uploadImage(client, pwd+"/test.png")

}
