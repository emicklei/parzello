gen:
	rm -f v1/queue/parcel.pb.go
	rm -f v1/parcello.pb.go
	protoc -I. --go_out=plugins=grpc:${GOPATH}/src v1/parcello.proto
	protoc -I. --go_out=plugins=grpc:${GOPATH}/src v1/queue/parcel.proto

fmt:
	protofmt -w v1/parcello.proto
	protofmt -w v1/queue/parcel.proto

run:
	go run *.go -v