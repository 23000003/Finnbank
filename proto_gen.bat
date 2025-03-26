@echo off

<<<<<<< HEAD
start protoc --proto_path=.protobuf ".protobuf/statement.proto" ^
    --go_out=common/grpc/statement ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=common/grpc/statement ^
=======
protoc --proto_path=protobuf "protobuf/bankcard.proto" ^
    --go_out=common/grpc/bankcard ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=common/grpc/bankcard ^
>>>>>>> 61b4e36081ee33e3cb41db22673ac0e9d18e1b40
    --go-grpc_opt=paths=source_relative

REM Check if the previous command failed
IF %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to generate gRPC code from bankcard.proto
    exit /b %ERRORLEVEL%
)

<<<<<<< HEAD
echo Success: gRPC code generated successfully.
=======
echo Success: gRPC code generated successfully.
>>>>>>> 61b4e36081ee33e3cb41db22673ac0e9d18e1b40
