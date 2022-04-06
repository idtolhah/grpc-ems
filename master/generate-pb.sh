source ~/.bash_profile
protoc masterpb/master.proto --go_out=. --go-grpc_out=. 