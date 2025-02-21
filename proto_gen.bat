@echo off

start protoc --proto_path=protobuf "protobuf/<ChangeFIleName>.proto" ^
    --go_out=services/common/grpc/<ChangeDirectoryName> ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=services/common/grpc/<ChangeDirectoryName> ^
    --go-grpc_opt=paths=source_relative