#!/bin/bash
cd ../backend/services/bff
go run main.go & P1=$!

cd ../../..
cd ../backend/services/user
go run main.go & P2=$!

cd ../../..
cd ../backend/services/advice
go run main.go & P3=$!

wait $P1 $P2 $P3