FROM golang:1.25-alpine3.22 AS build

ENV GOPROXY=https://proxy.golang.org,https://proxy.golang.org/cached-only,https://proxy.golangproxy.io/,direct

WORKDIR /app

COPY . .

RUN go build -o ./tracker ./cmd/main.go

FROM scratch

COPY --from=build /app/tracker /bin/tracker
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["bin/tracker"]