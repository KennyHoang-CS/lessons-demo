package diff

import (
	"context"
	"strings"

	"github.com/kenny/lessons/demo/pb"
)

type Service struct {
	pb.UnimplementedDiffServiceServer
	contentClient pb.ContentServiceClient
	cache         *Cache
}

func NewService(contentClient pb.ContentServiceClient, cache *Cache) *Service {
	return &Service{contentClient: contentClient, cache: cache}
}

func (s *Service) DiffContent(ctx context.Context, req *pb.DiffContentRequest) (*pb.DiffContentResponse, error) {
	if s.cache != nil {
		if cached, ok, err := s.cache.Get(ctx, req.FromContentId, req.ToContentId); err == nil && ok {
			return &pb.DiffContentResponse{Diff: cached}, nil
		}
	}

	from, err := s.contentClient.GetContent(ctx, &pb.GetContentRequest{
		ContentId: req.FromContentId,
	})
	if err != nil {
		return nil, err
	}

	to, err := s.contentClient.GetContent(ctx, &pb.GetContentRequest{
		ContentId: req.ToContentId,
	})
	if err != nil {
		return nil, err
	}

	diff := computeDiff(from.Body, to.Body)

	if s.cache != nil {
		_ = s.cache.Set(ctx, req.FromContentId, req.ToContentId, diff)
	}

	return &pb.DiffContentResponse{Diff: diff}, nil
}

// computeDiff returns a unified-style diff between two text blobs.
// Example output:
// - old line
// + new line
func computeDiff(a, b string) string {
	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")

	var out strings.Builder

	i, j := 0, 0
	for i < len(aLines) || j < len(bLines) {
		if i < len(aLines) && j < len(bLines) {
			if aLines[i] == bLines[j] {
				i++
				j++
				continue
			}
			out.WriteString("- " + aLines[i] + "\n")
			out.WriteString("+ " + bLines[j] + "\n")
			i++
			j++
			continue
		}
		if i < len(aLines) {
			out.WriteString("- " + aLines[i] + "\n")
			i++
			continue
		}
		if j < len(bLines) {
			out.WriteString("+ " + bLines[j] + "\n")
			j++
			continue
		}
	}

	return out.String()
}
