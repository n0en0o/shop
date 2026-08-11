package promotiongrpc

import (
	"context"
	"fmt"

	"github.com/n0en0o/marketplace/internal/promotion/grpc/pb"
)

type GreeterService struct {
	pb.UnimplementedGreeterServer
}

func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

func (s *GreeterService) SayHello(
	_ context.Context,
	req *pb.HelloRequest,
) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Msg: fmt.Sprintf("Hello %s!", req.GetName()),
	}, nil
}

func (s *GreeterService) Add(
	_ context.Context,
	req *pb.AddRequest,
) (*pb.AddReply, error) {
	return &pb.AddReply{
		Result: req.GetA() + req.GetB(),
	}, nil
}
