FROM golang:1.12.7 as dev

WORKDIR $GOPATH/src/github.com/sample-full-api
ADD . $GOPATH/src/github.com/sample-full-api

RUN go get -u github.com/kardianos/govendor
RUN go get -u github.com/swaggo/swag/cmd/swag

RUN govendor init
RUN govendor sync
RUN swag init

RUN go build

FROM alpine:3.7 as prod

ENV PORT=8080 \
    LOG_LEVEL=INFO \
    DB_USER=root \
    DB_PASS=root \
    DB_HOST=localhost:3306 \
    PROJECT_DIR=/go/src/github.com/sample-full-api

WORKDIR $PROJECT_DIR

COPY --from=dev $PROJECT_DIR/sample-full-api .

EXPOSE 8080
