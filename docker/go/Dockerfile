# build
FROM golang:1.24.5-alpine AS buildrdb

WORKDIR /app

COPY ../../go.mod ./
RUN go mod download

COPY ../.. ./

WORKDIR /app/cmd
RUN go build -o /rdb-app

# Final image
FROM alpine:latest
COPY --from=buildrdb /rdb-app /rdb-app

EXPOSE 8080

CMD ["./rdb-app"]
