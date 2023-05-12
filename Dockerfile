FROM golang:1.19-alpine

RUN apk update && apk upgrade && apk add --update alpine-sdk && \
    apk add --no-cache bash git openssh make cmake ca-certificates

ENV PROJECT_DIR $GOPATH/bin/app
WORKDIR $PROJECT_DIR

RUN go install github.com/cosmtrek/air@latest

COPY ./go.* ./
RUN go mod download

COPY ./ ./

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOBIN=$GOPATH/bin

ENTRYPOINT ["air", "-c", "./config/air.toml"]