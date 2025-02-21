@echo off
echo Starting all microservices...

start /D "services\account" go run .
start /D "services\transaction" go run .
start /D "services\bankcard" go run .
start /D "services\statement" go run .
start /D "services\api-gateway" go run .
start /D "services\graphql" go run .
start /D "services\notification" go run .

echo All microservices started.
pause