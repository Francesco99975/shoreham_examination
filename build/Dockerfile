FROM golang:1.21.2-alpine3.18 AS build

RUN apk --no-cache add gcc g++ make git

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/easypix ./cmd/main.go

FROM alpine:3.18

RUN apk update && apk upgrade && apk --no-cache add ca-certificates

WORKDIR /go/bin

COPY --from=build /go/src/app/bin /go/bin

EXPOSE 8020

ENTRYPOINT /go/bin/easypix --port 8020