package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/saurabh/Microservices/account"
	"github.com/saurabh/Microservices/catalog"
	"github.com/saurabh/Microservices/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	pb.UnimplementedOrderServiceServer
}

func ListenGRPC(s Service, accountURL string, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	pb.RegisterOrderServiceServer(srv, &grpcServer{
		UnimplementedOrderServiceServer: pb.UnimplementedOrderServiceServer{}, s, accountClient, catalogClient})

	reflection.Register(srv)
	return srv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("Error in getting account:", err)
		return nil, errors.New("account not found")
	}
	productIDs := []string{}

}

func (s *grpcServer) GetOrder(ctx context.Context, r *pb.GetOrderRequest) (*pb.GetOrderResposne, error) {

}
func (s *grpcServer) GetOrdersFromAccount(ctx context.Context, r *pb.GetOrdersFromAccountRequest) (*pb.GetOrdersFromAccountResponse, error) {

}
