#!/bin/bash
cd ../services/user/pb/
protoc --go_out=. --go-grpc_out=. user.proto 

cd ../../..

cd ../services/advice/pb/
protoc --go_out=. --go-grpc_out=. advice.proto