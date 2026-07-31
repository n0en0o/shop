package grpc

import (
	"context"

	"github.com/n0en0o/marketplace/internal/promotion/applications/queries"
	"github.com/n0en0o/marketplace/internal/promotion/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PromotionService struct {
	pb.UnimplementedPromotionServiceServer
	queryHandler *queries.GetByCatalogItemHandler
}

func NewPromotionService(
	queryHandler *queries.GetByCatalogItemHandler,
) *PromotionService {
	return &PromotionService{
		queryHandler: queryHandler,
	}
}

func (s *PromotionService) GetPromoByCatalogItem(
	ctx context.Context,
	req *pb.GetPromoByCatalogItemRequest,
) (*pb.GetPromoByCatalogItemResponse, error) {
	query := queries.GetByCatalogItemQuery{
		CatalogItemID: req.CatalogItemId,
	}

	p, err := s.queryHandler.Handle(ctx, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	if p == nil {
		return nil, status.Errorf(
			codes.NotFound,
			"ничего для %s не найдено", req.CatalogItemId,
		)
	}

	return &pb.GetPromoByCatalogItemResponse{
		Promotion: &pb.Promotion{
			Id:            p.ID,
			CatalogItemId: p.CatalogItemID,
			Title:         p.Title,
			Value:         p.Value,
		},
	}, nil
}
