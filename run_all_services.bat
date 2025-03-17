@echo off
echo Starting all microservices...

start /D "internal-services\account" go run .
start /D "internal-services\transaction" go run .
start /D "internal-services\bankcard" go run .
start /D "internal-services\statement" go run .
start /D "api-gateway" go run .
start /D "internal-services\graphql" go run .
start /D "internal-services\notification" go run .

echo All microservices started.
pause