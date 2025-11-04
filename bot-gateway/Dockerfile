FROM golang:1.25-alpine3.22 AS build

#ENV GOPROXY=https://proxy.golang.org,https://proxy.golang.org/cached-only,https://proxy.golangproxy.io/,direct

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# RUN go get golang.org/x/text
 
RUN go build -v -o ./gateway ./cmd/main.go

FROM scratch

COPY --from=build /app/gateway /bin/bot
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["bin/bot"]