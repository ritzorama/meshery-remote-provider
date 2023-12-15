FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /out/provider ./cmd/provider

FROM alpine:3.20

WORKDIR /app
COPY --from=build /out/provider /usr/local/bin/provider

EXPOSE 8080

ENTRYPOINT ["provider"]
