package promotiongrpc

import (
	"context"
	"errors"

	"github.com/n0en0o/marketplace/internal/promotion/applications/commands"
	"github.com/n0en0o/marketplace/internal/promotion/applications/queries"
	"github.com/n0en0o/marketplace/internal/promotion/domain"
	"github.com/n0en0o/marketplace/internal/promotion/domain/repositories"
	"github.com/n0en0o/marketplace/internal/promotion/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PromotionService struct {
	pb.UnimplementedPromotionServiceServer
	queryHandler  *queries.GetByCatalogItemHandler
	createHandler *commands.CreatePromoHandler
	updateHandler *commands.UpdatePromoHandler
	deleteHandler *commands.DeletePromoHandler
}

func NewPromotionService(
	queryHandler *queries.GetByCatalogItemHandler,
	createHandler *commands.CreatePromoHandler,
	updateHandler *commands.UpdatePromoHandler,
	deleteHandler *commands.DeletePromoHandler,
) *PromotionService {
	return &PromotionService{
		queryHandler:  queryHandler,
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
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
		return nil, promoErrorStatus(err)
	}

	if p == nil {
		return nil, status.Error(codes.NotFound, "акция не найдена")
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

func (s *PromotionService) CreatePromo(
	ctx context.Context,
	req *pb.CreatePromoRequest,
) (*pb.CreatePromoResponse, error) {
	cmd := commands.CreatePromoCommand{
		CatalogItemID: req.CatalogItemId,
		Title:         req.Title,
		Value:         req.Value,
	}

	result, err := s.createHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, promoErrorStatus(err)
	}

	return &pb.CreatePromoResponse{
		Id:          result.ID,
		Success:     result.Success,
		Description: result.Description,
	}, nil
}

func (s *PromotionService) UpdatePromo(
	ctx context.Context,
	req *pb.UpdatePromoRequest,
) (*pb.UpdatePromoResponse, error) {
	cmd := commands.UpdatePromoCommand{
		ID:    req.Id,
		Title: req.Title,
		Value: req.Value,
	}

	result, err := s.updateHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, promoErrorStatus(err)
	}

	return &pb.UpdatePromoResponse{
		Success:     result.Success,
		Description: result.Description,
	}, nil
}

func (s *PromotionService) DeletePromo(
	ctx context.Context,
	req *pb.DeletePromoRequest,
) (*pb.DeletePromoResponse, error) {
	cmd := commands.DeletePromoCommand{
		ID: req.Id,
	}

	result, err := s.deleteHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, promoErrorStatus(err)
	}

	return &pb.DeletePromoResponse{
		Success:     result.Success,
		Description: result.Description,
	}, nil
}

func promoErrorStatus(err error) error {
	switch {
	case errors.Is(err, commands.ErrInvalidPromoValue):
		return status.Error(codes.InvalidArgument, "некорректное значение скидки")
	case errors.Is(err, domain.ErrPromoAlreadyExists):
		return status.Error(codes.AlreadyExists, "акция уже существует")
	case errors.Is(err, repositories.ErrPromoNotFound):
		return status.Error(codes.NotFound, "акция не найдена")
	default:
		return status.Error(codes.Internal, "внутренняя ошибка сервера")
	}
}
