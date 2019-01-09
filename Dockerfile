FROM golang:1.10.1
WORKDIR /go/src/github.com/emicklei/parzello/
COPY . .
ARG version
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$version" .

FROM scratch
COPY --from=0 /go/src/github.com/emicklei/parzello .

# HTTP port
EXPOSE 8080

ENTRYPOINT ["/parzello"]