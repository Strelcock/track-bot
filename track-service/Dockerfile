FROM golang:1.25-alpine3.21 AS build

WORKDIR /app

COPY . .

RUN go build -o ./tracker ./cmd/main.go

FROM scratch

COPY --from=build /app/tracker /bin/tracker
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["bin/tracker"]