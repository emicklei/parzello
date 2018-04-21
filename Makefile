.PHONY: gen fmt run dock drun

gen:
	rm -f v1/parzello.pb.go
	protoc -I. --go_out=plugins=grpc:${GOPATH}/src v1/parzello.proto

fmt:
	protofmt -w v1/parzello.proto

run:
	go run *.go -v

dock:
	docker build -t parzello .

drun:
	docker run -it \
		-v ~/.config/gcloud/:/gcloud \
		-e GOOGLE_APPLICATION_CREDENTIALS=/gcloud/application_default_credentials.json \
		parzello