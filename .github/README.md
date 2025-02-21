# Online Banking System
IT 3101A - Microservices / QA Final Project

## Golang Commands

   ```bash
   go run ./services/<ur-assigned-service> . 
   go get -u ./<ur-assigned-service> # Import external packages
   ./run_all_services.bat # run all services
   go mod tidy # on ur specific service directory for package cleanup/update
   ```

## Note:
Utillize ``logger`` from finnbank/common/utils for terminal logs

## Service URL's
   ```bash
   # http ports (for data response test)
   http://localhost:8080/api # api-gateway 
   http://localhost:8081/api/bankcard # bankcard-api 
   http://localhost:8082/api/account # account-api
   http://localhost:8083/api/graphql # graphql-api (not http)
   http://localhost:8084/api/statement # statement-api
   http://localhost:8085/api/transaction # transaction-api
   http://localhost:8086/api/notification # notification-api


   # gRPC ports (main communication)
   http://localhost:9001/api/bankcard # bankcard-grpc
   http://localhost:9002/api/account # account-grpc
   http://localhost:9004/api/statement # statement-grpc 
   http://localhost:9005/api/transaction # transaction-grpc 
   http://localhost:9006/api/notification # notification-grpc
   ```

## Generate proto for gRPC communication
   - Only do this if you added a proto file in ``Protobuf Directory``
   - This is to create gRPC services/communication of your respective service
   ```bash
   # install first
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   $env:PATH="$PATH:$(go env GOPATH)/bin"
   protoc --version # check installation

   # if its unrecognized then try:
   winget install protobuf

   # update proto_gen.bat
   start protoc \
        --proto_path=protobuf "protobuf/<change-lines-with-this-format>.proto" \

   # then generate
   ./proto_gen.bat
   ```

## Microservice Architecture [Not Final]
![alt text](PROJECT.drawio%20(1).png)