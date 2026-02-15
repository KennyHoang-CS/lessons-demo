package content

import (
	"context"

	"github.com/kenny/lessons/demo/pb"
)

type Service struct {
	pb.UnimplementedContentServiceServer
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateContent(ctx context.Context, req *pb.CreateContentRequest) (*pb.CreateContentResponse, error) {
	id, err := s.store.SaveContent(ctx, req.Body)
	if err != nil {
		return nil, err
	}
	return &pb.CreateContentResponse{ContentId: id}, nil
}

func (s *Service) GetContent(ctx context.Context, req *pb.GetContentRequest) (*pb.GetContentResponse, error) {
	body, err := s.store.GetContent(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}
	return &pb.GetContentResponse{Body: body}, nil
}
