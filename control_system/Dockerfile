FROM golang:1.24 AS builder
WORKDIR /app
COPY ./ /app/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
EXPOSE 8080
CMD ["./main"]
