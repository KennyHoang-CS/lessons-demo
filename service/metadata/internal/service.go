package metadata

import (
	"context"

	"github.com/kenny/lessons/demo/pb"
)

type Service struct {
	pb.UnimplementedMetadataServiceServer
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateLesson(ctx context.Context, req *pb.CreateLessonRequest) (*pb.CreateLessonResponse, error) {
	rootID, versionNumber, err := s.store.CreateLesson(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}
	return &pb.CreateLessonResponse{
		RootId:        rootID,
		VersionNumber: versionNumber,
	}, nil
}

func (s *Service) CreateVersion(ctx context.Context, req *pb.CreateVersionRequest) (*pb.CreateVersionResponse, error) {
	versionNumber, err := s.store.CreateVersion(ctx, req.RootId, req.ContentId)
	if err != nil {
		return nil, err
	}
	return &pb.CreateVersionResponse{
		RootId:        req.RootId,
		VersionNumber: versionNumber,
	}, nil
}

func (s *Service) PublishVersion(ctx context.Context, req *pb.PublishVersionRequest) (*pb.PublishVersionResponse, error) {
	if err := s.store.PublishVersion(ctx, req.RootId, req.VersionNumber); err != nil {
		return nil, err
	}
	return &pb.PublishVersionResponse{}, nil
}

func (s *Service) GetVersion(ctx context.Context, req *pb.GetVersionRequest) (*pb.GetVersionResponse, error) {
	return s.store.GetVersion(ctx, req.RootId, req.VersionNumber)
}
