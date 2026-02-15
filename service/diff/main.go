package main

import (
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/kenny/lessons/demo/pb"
	diff "github.com/kenny/lessons/demo/service/diff/internal"
)

func main() {
	contentAddr := os.Getenv("CONTENT_GRPC_ADDR") // "content-service:50051"
	if contentAddr == "" {
		contentAddr = "content-service:50051"
	}
	redisAddr := os.Getenv("REDIS_ADDR") // "redis:6379"
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	contentConn, err := grpc.Dial(contentAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	contentClient := pb.NewContentServiceClient(contentConn)

	cache := diff.NewCache(redisAddr, 10*time.Minute)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	svc := diff.NewService(contentClient, cache)
	pb.RegisterDiffServiceServer(grpcServer, svc)

	log.Println("DiffService listening on :50053")
	log.Fatal(grpcServer.Serve(lis))
}
