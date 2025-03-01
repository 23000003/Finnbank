@echo off

     protoc --proto_path=protobuf "protobuf/account.proto" ^
    --go_out=services/common/grpc/account ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=services/common/grpc/account ^
    --go-grpc_opt=paths=source_relative