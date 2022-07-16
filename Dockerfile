FROM golang:1.18-alpine3.16 AS build

WORKDIR /build

COPY ./ ./

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o boggle -ldflags="-s -w" cmd/main.go

FROM alpine:3.16.0

WORKDIR /app

RUN mkdir -p /app/configs && \
    mkdir -p /app/data

COPY --from=build /build/boggle /app/boggle
COPY ./configs/config.hcl ./configs
COPY ./data/words.zip ./data

USER 1000

ENTRYPOINT [ "/app/boggle" ]