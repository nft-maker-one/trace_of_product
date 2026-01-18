FROM golang:1.24
WORKDIR /app
COPY ./ /app/
RUN go mod download
RUN go build -o main
EXPOSE 8081
CMD ["./main", "-chain", "-ip", "0.0.0.0", "-port", "8081"]