.PHONY: gen fmt run dock drun

gen:
	rm -f v1/parcello.pb.go
	protoc -I. --go_out=plugins=grpc:${GOPATH}/src v1/parcello.proto

fmt:
	protofmt -w v1/parcello.proto

run:
	go run *.go -v

dock:
	docker build -t parcello .

drun:
	docker run -it \
		-v ~/.config/gcloud/:/gcloud \
		-e GOOGLE_APPLICATION_CREDENTIALS=/gcloud/application_default_credentials.json \
		parcello