@echo off

start protoc --proto_path=.protobuf ".protobuf/<Change>.proto" ^
    --go_out=common/grpc/<change> ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=common/grpc/<chage> ^
    --go-grpc_opt=paths=source_relative