#!/bin/bash
cd user
go run main.go & P1=$!

cd ..
cd master
go run main.go & P2=$!

cd ..
cd bff
go run main.go & P3=$!

wait $P1 $P2 $P3