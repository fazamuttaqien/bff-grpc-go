package client

import (
	"context"
	"errors"
	"net/http"

	"github.com/fazamuttaqien/ms-grpc/bff/pb/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Age      int32  `json:"age"`
	Greeting string `json:"greeting"`
}

type UserDetails struct {
	Salary int32  `json:"salary"`
	Power  string `json:"power"`
	Advice
}

type UserWithDetails struct {
	User
	UserDetails
}

type UserClient struct {
}

var (
	_                     = loadLocalEnv()
	userGrpcService       = GetEnv("USER_GRPC_SERVICE")
	userGrpcServiceClient pb.UserServiceClient
)

func prepareUserGrpcClient(c *context.Context) error {
	conn, err := grpc.NewClient(userGrpcService, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}...)

	if err != nil {
		userGrpcServiceClient = nil
		return errors.New("connection to user gRPC service failed")
	}

	if userGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	userGrpcServiceClient = pb.NewUserServiceClient(conn)
	return nil
}

func (uc *UserClient) CreateUpdateUser(user UserWithDetails, c *context.Context, method string) (string, error) {
	if err := prepareUserGrpcClient(c); err != nil {
		return "", err
	}

	op := pb.Operation_CREATE
	if method == http.MethodPut {
		op = pb.Operation_UPDATE
	}

	res, err := userGrpcServiceClient.CreateUpdateUser(*c, &pb.CreateUpdateUserRequest{Operation: op,
		Id: user.Id, Name: user.Name, Age: user.Age,
		Greeting: user.Greeting, Salary: user.Salary, Power: user.Power,
	})
	if err != nil {
		return "", errors.New(status.Convert(err).Message())
	}
	return res.Id, nil
}

func (uc *UserClient) GetUserDetails(id string, c *context.Context) (*UserDetails, error) {
	if err := prepareUserGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := userGrpcServiceClient.GetUserDetails(*c, &pb.GetUserDetailsRequest{Id: id})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}
	return &UserDetails{Salary: res.Salary, Power: res.Power}, nil
}

func (uc *UserClient) GetUsers(c *context.Context) (*[]User, error) {
	if err := prepareUserGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := userGrpcServiceClient.GetUsers(*c, &pb.GetUsersRequest{})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}

	var users []User
	for _, u := range res.GetUsers() {
		users = append(users, User{Id: u.Id, Name: u.Name, Age: u.Age, Greeting: u.Greeting})
	}
	return &users, nil
}
