package grpc

import "context"

const Address = "127.0.0.1:50052"

type helloService struct {
}

var HelloService = helloService{}

func (h helloService) SayHello(ctx context.Context, in *pb.H) {

}
