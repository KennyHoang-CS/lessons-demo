package main

import (
	"log"
	"net"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver

	"google.golang.org/grpc"

	"github.com/kenny/lessons/demo/pb"
	metadata "github.com/kenny/lessons/demo/service/metadata/internal"
)

func main() {
	// 1. Connect to Postgres using sqlx. Prefer POSTGRES_DSN env var.
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://lessons:lessons@localhost:5432/lessons?sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("cannot connect to Postgres:", err)
	}

	// 2. Create store
	store := metadata.NewStore(db)

	// 3. Start gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	svc := metadata.NewService(store)

	pb.RegisterMetadataServiceServer(grpcServer, svc)

	log.Println("MetadataService listening on :50052")
	log.Fatal(grpcServer.Serve(lis))
}
