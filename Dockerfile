FROM golang:1.18-alpine3.16 AS build

WORKDIR /build

COPY ./ ./

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o boggle -ldflags="-s -w" cmd/main.go

FROM alpine:3.16.0

WORKDIR /app

COPY --from=build /build/boggle /app/boggle

USER 1000

ENTRYPOINT [ "/app/boggle" ]