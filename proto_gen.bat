@echo off

start protoc --proto_path=.protobuf ".protobuf/statement.proto" ^
    --go_out=common/grpc/statement ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=common/grpc/statement ^
    --go-grpc_opt=paths=source_relative

REM Check if the previous command failed
IF %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to generate gRPC code from bankcard.proto
    exit /b %ERRORLEVEL%
)

echo Success: gRPC code generated successfully.