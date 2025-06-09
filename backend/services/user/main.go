package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/fazamuttaqien/ms-grpc/user/db"
	"github.com/fazamuttaqien/ms-grpc/user/pb"

	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var timeout = time.Second

type server struct {
	pb.UnimplementedUserServiceServer
}

func (*server) CreateUpdateUser(ctx context.Context, req *pb.CreateUpdateUserRequest) (*pb.CreateUpdateUserResponse, error) {
	log.Println("Called CreateUpdateUser Operation", req.Operation)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var uid bson.ObjectID
	if req.Operation == pb.Operation_CREATE {
		uid = bson.NewObjectID()
	} else {
		var err error
		uid, err = bson.ObjectIDFromHex(req.Id)
		if err != nil {
			return nil, errorResponse(err)
		}
	}

	err := db.UpsertOne(c, &db.User{
		Id: uid, Name: req.Name,
		Age: req.Age, Greeting: req.Greeting,
		Salary: req.Salary, Power: req.Power})
	if err != nil {
		return nil, errorResponse(err)
	}

	return &pb.CreateUpdateUserResponse{Id: uid.Hex()}, nil
}

func (*server) GetUserDetails(ctx context.Context, req *pb.GetUserDetailsRequest) (*pb.GetUserDetailsResponse, error) {
	log.Println("Called GetUserDetails, Id", req.Id)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	uid, err := bson.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, errorResponse(err)
	}

	result, err := db.FindOne(c, uid)
	if err != nil {
		return nil, errorResponse(err)
	}

	return &pb.GetUserDetailsResponse{Salary: result.Salary, Power: result.Power}, nil
}

func (*server) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	log.Println("Called GetUsers")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := db.Find(c)

	if err != nil {
		return nil, errorResponse(err)
	}

	var res pb.GetUsersResponse
	for _, d := range *data {
		res.Users = append(res.Users, &pb.GetUserResponse{Id: d.Id.Hex(), Name: d.Name, Age: d.Age, Greeting: d.Greeting})
	}

	return &res, nil
}

func errorResponse(err error) error {
	log.Fatalln("Error:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

func main() {
	log.Println("User Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Error:", err.Error())
	}

	db.MongoClient, err = db.NewClient(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.MongoClient.Disconnect(context.Background())

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Fatalln("Error:", err.Error())
	}
}
