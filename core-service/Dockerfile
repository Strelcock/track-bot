FROM golang:1.25-alpine3.21 AS build

WORKDIR /app

COPY . .

RUN go build -o ./core ./cmd/main.go

FROM scratch


COPY --from=build /app/core /bin/bot


CMD ["bin/bot"]