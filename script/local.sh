#!/bin/bash
cd ../services/bff
go run main.go & P1=$!

cd ../..
cd ../services/user
go run main.go & P2=$!

cd ../..
cd ../services/advice
go run main.go & P3=$!

wait $P1 $P2 $P3