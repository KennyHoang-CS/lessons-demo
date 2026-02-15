package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"google.golang.org/grpc"

	"github.com/kenny/lessons/demo/pb"
)

type Server struct {
	contentClient  pb.ContentServiceClient
	metadataClient pb.MetadataServiceClient
	diffClient     pb.DiffServiceClient
}

func main() {
	// Connect to gRPC services (use env vars with sensible defaults)
	contentAddr := os.Getenv("CONTENT_GRPC_ADDR")
	if contentAddr == "" {
		contentAddr = "content-service:50051"
	}
	metadataAddr := os.Getenv("METADATA_GRPC_ADDR")
	if metadataAddr == "" {
		metadataAddr = "metadata-service:50052"
	}
	diffAddr := os.Getenv("DIFF_GRPC_ADDR")
	if diffAddr == "" {
		diffAddr = "diff-service:50053"
	}

	contentConn, err := grpc.Dial(contentAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	metadataConn, err := grpc.Dial(metadataAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	diffConn, err := grpc.Dial(diffAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	s := &Server{
		contentClient:  pb.NewContentServiceClient(contentConn),
		metadataClient: pb.NewMetadataServiceClient(metadataConn),
		diffClient:     pb.NewDiffServiceClient(diffConn),
	}

	app := fiber.New()

	// API routes
	app.Post("/lessons", s.handleCreateLesson)
	app.Post("/lessons/clone", s.handleCloneLesson)
	app.Post("/lessons/:rootID/versions", s.handleCreateVersion)
	app.Post("/lessons/:rootID/versions/:version/publish", s.handlePublish)
	app.Get("/lessons/:rootID/versions/:version", s.handleGetVersion)
	app.Get("/diff", s.handleDiff)

	// Serve GUI
	app.Get("/*", static.New("./static"))

	log.Println("API Gateway listening on :8080")
	log.Fatal(app.Listen(":8080"))
}

func (s *Server) handleCreateLesson(c fiber.Ctx) error {
	body := string(c.Body())

	// 1) Store content
	contentResp, err := s.contentClient.CreateContent(context.Background(), &pb.CreateContentRequest{
		Body: body,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	// 2) Create lesson root + version 1
	metaResp, err := s.metadataClient.CreateLesson(context.Background(), &pb.CreateLessonRequest{
		ContentId: contentResp.ContentId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	return c.JSON(fiber.Map{
		"root_id":        metaResp.RootId,
		"version_number": metaResp.VersionNumber,
		"content_id":     contentResp.ContentId,
	})
}

func (s *Server) handleCreateVersion(c fiber.Ctx) error {
	rootID := c.Params("rootID")
	body := string(c.Body())

	// 1) Store new content
	contentResp, err := s.contentClient.CreateContent(context.Background(), &pb.CreateContentRequest{
		Body: body,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	// 2) Create new version
	metaResp, err := s.metadataClient.CreateVersion(context.Background(), &pb.CreateVersionRequest{
		RootId:    rootID,
		ContentId: contentResp.ContentId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	// Return metadata along with the newly created content id so clients can diff
	return c.JSON(fiber.Map{
		"root_id":        metaResp.RootId,
		"version_number": metaResp.VersionNumber,
		"content_id":     contentResp.ContentId,
	})
}

func (s *Server) handlePublish(c fiber.Ctx) error {
	rootID := c.Params("rootID")
	versionStr := c.Params("version")

	var version int32
	_, err := fmt.Sscanf(versionStr, "%d", &version)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid version")
	}

	_, err = s.metadataClient.PublishVersion(context.Background(), &pb.PublishVersionRequest{
		RootId:        rootID,
		VersionNumber: version,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	return c.JSON(fiber.Map{"status": "published"})
}

func (s *Server) handleGetVersion(c fiber.Ctx) error {
	rootID := c.Params("rootID")
	versionStr := c.Params("version")

	var version int32
	_, err := fmt.Sscanf(versionStr, "%d", &version)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid version")
	}

	resp, err := s.metadataClient.GetVersion(context.Background(), &pb.GetVersionRequest{
		RootId:        rootID,
		VersionNumber: version,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	return c.JSON(resp)
}

func (s *Server) handleDiff(c fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing from/to")
	}

	resp, err := s.diffClient.DiffContent(context.Background(), &pb.DiffContentRequest{
		FromContentId: from,
		ToContentId:   to,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	return c.JSON(resp)
}

type cloneRequest struct {
	FromRootID  string `json:"from_root_id"`
	FromVersion int32  `json:"from_version"`
}

type cloneResponse struct {
	RootId    string `json:"root_id"`
	Version   int32  `json:"version_number"`
	ContentId string `json:"content_id"`
}

// POST /lessons/clone
func (s *Server) handleCloneLesson(c fiber.Ctx) error {
	var req cloneRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	// 1) Get source version metadata
	src, err := s.metadataClient.GetVersion(context.Background(), &pb.GetVersionRequest{
		RootId:        req.FromRootID,
		VersionNumber: req.FromVersion,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	// 2) Fetch content body from content service
	cont, err := s.contentClient.GetContent(context.Background(), &pb.GetContentRequest{ContentId: src.ContentId})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	// 3) Store as a new content blob
	newCont, err := s.contentClient.CreateContent(context.Background(), &pb.CreateContentRequest{Body: cont.Body})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	// 4) Create new lesson referencing the new content
	newLesson, err := s.metadataClient.CreateLesson(context.Background(), &pb.CreateLessonRequest{ContentId: newCont.ContentId})
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	return c.JSON(cloneResponse{RootId: newLesson.RootId, Version: newLesson.VersionNumber, ContentId: newCont.ContentId})
}
