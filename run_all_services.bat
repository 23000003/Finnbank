@echo off
echo Starting all microservices...

start /D account-service go run main.go
start /D transaction-service go run main.go
start /D bankcard-service go run main.go
start /D statement-service go run main.go
start /D graphql-db-service go run main.go
start /D api-gateway go run main.go

echo All microservices started.
pause