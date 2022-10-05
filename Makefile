proto:
	protoc -I . api/v1/*.proto --go-grpc_out=. --go_out=.
server:
	go run cmd/main.go

vet:
	go vet cmd/main.go
	shadow cmd/main.go
.PHONY:vet