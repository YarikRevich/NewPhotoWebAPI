.PHONY: stub proto build

stub:
	@echo "It's a stub"

proto:
	export PATH=$$PATH:$$GOPATH/bin;\
	protoc -I $$GOPATH/src/NewPhotoWebAPI  --go_out=$$GOPATH/src/NewPhotoWebAPI  --go-grpc_out=$$GOPATH/src/NewPhotoWebAPI  $$GOPATH/src/NewPhotoWebAPI/logic/proto/api.proto

build: proto
	go build main.go
