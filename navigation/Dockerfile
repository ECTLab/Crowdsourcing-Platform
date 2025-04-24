FROM golang:1.23.5 AS BUILD_IMAGE

RUN apt update
RUN apt install protobuf-compiler build-essential -y
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
WORKDIR /root/src
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /root/app/ ./...



FROM ubuntu:22.04 AS RUN_IMAGE

RUN apt update && apt install graphviz -y

COPY ./scripts/entrypoint.sh /scripts/
COPY --from=BUILD_IMAGE /root/app/navigation /app/navigation
COPY --from=BUILD_IMAGE /root/app/readiness /app/readiness

WORKDIR /app/

COPY ./config/development.yaml /app/config/

RUN chmod +x /scripts/entrypoint.sh

ENTRYPOINT ["/scripts/entrypoint.sh"]

