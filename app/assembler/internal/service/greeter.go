package service

import (
	v1 "rag/api/assembler/v1"
	"rag/app/assembler/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedAssemblerServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
// func (s *GreeterService) SayHello(ctx context.Context, in *v1.AssembleContextRequest) (*v1.HelloReply, error) {
// 	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
// }
