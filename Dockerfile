FROM golang:1.10.1
WORKDIR /go/src/github.com/emicklei/parcello/
COPY . .
ARG version
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$version" .

FROM scratch
COPY --from=0 /go/src/github.com/emicklei/parcello .

# gRPC port
EXPOSE 9090

ENTRYPOINT ["/parcello"]