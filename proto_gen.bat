@echo off


start protoc --proto_path=.protobuf ".protobuf/statement.proto" ^
    --go_out=common/grpc/statement ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=common/grpc/statement ^
    --go-grpc_opt=paths=source_relative ^