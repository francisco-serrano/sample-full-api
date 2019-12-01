FROM golang:1.12.7 as dev

WORKDIR $GOPATH/src/github.com/sample-full-api
ADD . $GOPATH/src/github.com/sample-full-api

RUN go get -u github.com/kardianos/govendor
RUN go get -u github.com/swaggo/swag/cmd/swag

RUN govendor init
RUN govendor sync
RUN swag init

RUN CGO_ENABLED=0 go build

FROM alpine:3.7 as prod

ENV PROJECT_DIR=/go/src/github.com/sample-full-api

WORKDIR $PROJECT_DIR

COPY --from=dev $PROJECT_DIR/sample-full-api .

ENTRYPOINT ./sample-full-api
