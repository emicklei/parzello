gen:
	protoc -I. --go_out=plugins=grpc:. delivery.proto
	# protoc -I. --go_out=plugins=grpc:. parcel.proto
	cp delivery.pb.go example

fmt:
	protofmt -w delivery.proto

run:
	go run *.go