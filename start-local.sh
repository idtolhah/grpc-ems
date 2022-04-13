#!/bin/bash
cd user-query
go run main.go & P1=$!
cd ..

cd master-query
go run main.go & P2=$!
cd ..

cd packing-query
go run main.go & P3=$!
cd ..

cd packing-cmd
go run main.go & P4=$!
cd ..

cd bff
go run main.go & P5=$!

wait $P1 $P2 $P3 $P4 $P5