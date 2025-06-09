package client

import (
	"context"
	"errors"
	"net/http"

	"github.com/fazamuttaqien/ms-grpc/bff/pb/advice"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Advice struct {
	Advice     string `json:"advice"`
	Created_At string `json:"created_at"`
	UserId     string `json:"user_id,omitempty"`
}

type AdviceClient struct {
}

var (
	_                       = loadLocalEnv()
	adviceGrpcService       = GetEnv("ADVICE_GRPC_SERVICE")
	adviceGrpcServiceClient pb.AdviceServiceClient
)

func prepareAdviceGrpcClient(c *context.Context) error {
	conn, err := grpc.NewClient(adviceGrpcService, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}...)

	if err != nil {
		adviceGrpcServiceClient = nil
		return errors.New("connection to advice gRPC service failed")
	}

	if adviceGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	adviceGrpcServiceClient = pb.NewAdviceServiceClient(conn)
	return nil
}

func (ac *AdviceClient) CreateUpdateAdvice(advice Advice, c *context.Context, method string) error {
	if err := prepareAdviceGrpcClient(c); err != nil {
		return err
	}

	op := pb.Operation_CREATE
	if method == http.MethodPut {
		op = pb.Operation_UPDATE
	}

	if _, err := adviceGrpcServiceClient.CreateUpdateAdvice(*c, &pb.CreateUpdateAdviceRequest{
		Operation: op,
		UserId:    advice.UserId,
		Advice:    advice.Advice,
	}); err != nil {
		return errors.New(status.Convert(err).Message())
	}
	return nil
}

func (ac *AdviceClient) GetAdvice(id string, c *context.Context) (*Advice, error) {
	if err := prepareAdviceGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := adviceGrpcServiceClient.GetAdvice(*c, &pb.GetUserAdviceRequest{UserId: id})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}
	return &Advice{Advice: res.Advice, Created_At: res.CreatedAt.AsTime().String()}, nil
}
