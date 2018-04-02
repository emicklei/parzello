gen:
	rm -f v1/parcello.pb.go
	protoc -I. --go_out=plugins=grpc:${GOPATH}/src v1/parcello.proto

fmt:
	protofmt -w v1/parcello.proto

run:
	go run *.go -v