FROM golang:alpine AS build
WORKDIR /go/src/github.com/D1abloRUS/thingspeak-client
COPY / .
RUN apk add --no-cache git \
    && go get -u github.com/kelseyhightower/envconfig \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o thingspeak-client -ldflags '-s' .

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/D1abloRUS/thingspeak-client/thingspeak-client .
EXPOSE 3000

CMD ["/thingspeak-client"]