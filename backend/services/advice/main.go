package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/fazamuttaqien/ms-grpc/advice/db"
	"github.com/fazamuttaqien/ms-grpc/advice/pb"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	timeout         = time.Second
	postgresClient *pgxpool.Pool
)

type server struct {
	pb.UnimplementedAdviceServiceServer
}

func (*server) CreateUpdateAdvice(ctx context.Context, req *pb.CreateUpdateAdviceRequest) (*pb.CreateUpdateAdviceResponse, error) {
	log.Println("Called CreateUpdateAdvice, Operation", req.Operation)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var err error
	if req.Operation == pb.Operation_CREATE {
		err = db.CreateOne(postgresClient, c, &db.Advice{
			UserId: req.UserId,
			Advice: req.Advice})
	} else {
		err = db.UpdateOne(postgresClient, c, &db.Advice{
			UserId: req.UserId,
			Advice: req.Advice})
	}
	if err != nil {
		return nil, errorResponse(err)
	}

	return &pb.CreateUpdateAdviceResponse{}, nil
}

func (*server) GetAdvice(ctx context.Context, req *pb.GetUserAdviceRequest) (*pb.GetUserAdviceResponse, error) {
	log.Println("Called GetAdvice for User Id", req.UserId)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result, err := db.FindOne(postgresClient, c, req.UserId)
	if err != nil {
		return nil, errorResponse(err)
	}

	return &pb.GetUserAdviceResponse{Advice: result.Advice, CreatedAt: timestamppb.New(result.CreatedAt)}, nil
}

func errorResponse(err error) error {
	log.Fatalln("Error:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

func main() {
	log.Println("Advice Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalln("Error:", err.Error())
	}

	postgresClient, err = db.NewClient(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer postgresClient.Close()

	s := grpc.NewServer()
	pb.RegisterAdviceServiceServer(s, &server{})

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Fatalln("Error:", err.Error())
	}
}
