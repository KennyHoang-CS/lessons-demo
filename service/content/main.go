package main

import (
	"log"
	"net"
	"os"

	"github.com/kenny/lessons/demo/pb"
	content "github.com/kenny/lessons/demo/service/content/internal"

	"google.golang.org/grpc"
)

func main() {
	endpoint := os.Getenv("MINIO_ENDPOINT") // "minio:9000"
	if endpoint == "" {
		endpoint = "minio:9000"
	}
	accessKey := os.Getenv("MINIO_ACCESS_KEY") // "minioadmin"
	if accessKey == "" {
		accessKey = "minioadmin"
	}
	secretKey := os.Getenv("MINIO_SECRET_KEY") // "minioadmin"
	if secretKey == "" {
		secretKey = "minioadmin"
	}
	bucket := os.Getenv("MINIO_BUCKET") // "lessons"
	if bucket == "" {
		bucket = "lessons"
	}
	useSSL := false

	store, err := content.NewStore(endpoint, accessKey, secretKey, bucket, useSSL)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	svc := content.NewService(store)
	pb.RegisterContentServiceServer(grpcServer, svc)

	log.Println("ContentService listening on :50051")
	log.Fatal(grpcServer.Serve(lis))
}
