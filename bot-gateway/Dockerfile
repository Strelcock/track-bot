FROM golang:1.25-alpine3.21 AS build

ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /app

COPY . .

RUN go build -o ./gateway ./cmd/main.go

FROM scratch

COPY --from=build /app/gateway /bin/bot
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["bin/bot"]