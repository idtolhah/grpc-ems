source ~/.bash_profile
protoc pb/packingquerypb/packingquery.proto --go_out=. --go-grpc_out=. 