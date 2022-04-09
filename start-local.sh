#!/bin/bash
cd redis-cache
go run main.go & P1=$!

cd ..
cd user
go run main.go & P2=$!

cd ..
cd master
go run main.go & P3=$!

cd ..
cd bff
go run main.go & P4=$!

wait $P1 $P2 $P3 $P4