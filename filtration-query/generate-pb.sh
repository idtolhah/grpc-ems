source ~/.bash_profile
protoc pb/filtrationquerypb/filtrationquery.proto --go_out=. --go-grpc_out=. 