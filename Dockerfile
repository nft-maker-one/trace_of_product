FROM golang:latest as builder
WORKDIR /app
COPY ./trace_of_product /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main
FROM scratch
COPY --from=builder /app/main /
ENTRYPOINT [ "/main -chain","-ip","127.0.0.1","-port","8081" ]
