@echo off


start protoc --proto_path=.protobuf ".protobuf/auth.proto" ^
    --go_out=common/grpc/auth ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=common/grpc/auth ^
    --go-grpc_opt=paths=source_relative ^