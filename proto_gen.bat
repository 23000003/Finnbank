@echo off

start protoc --proto_path=.protobuf ".protobuf/account.proto" ^
    --go_out=common/grpc/account ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=common/grpc/account ^
    --go-grpc_opt=paths=source_relative