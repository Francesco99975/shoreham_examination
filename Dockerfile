FROM golang:1.21.4-alpine3.18 AS build

RUN apk --no-cache add gcc g++ make git

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN mv .prod.env .env

RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/shorehamex ./cmd/server/*.go

FROM alpine:3.18

RUN apk update && apk upgrade && apk --no-cache add ca-certificates

WORKDIR /go/bin

COPY --from=build /go/src/app/bin /go/bin
COPY --from=build /go/src/app/.env /go/bin/
COPY --from=build /go/src/app/data /go/bin/data
COPY --from=build /go/src/app/static /go/bin/static

EXPOSE 8020

ENTRYPOINT /go/bin/shorehamex --port 8020