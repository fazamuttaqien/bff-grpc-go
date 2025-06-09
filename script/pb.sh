#!/bin/bash
cd ../backend/services/user/pb/
protoc --go_out=. --go-grpc_out=. user.proto 

cd ../../../.. 

cd ../backend/services/advice/pb/
protoc --go_out=. --go-grpc_out=. advice.proto