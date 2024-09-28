FROM golang:1.20
WORKDIR /app
COPY ./ /app/
RUN go mod download
RUN go build -o main
CMD [ "./main","-ip","127.0.0.1""-port","8081","-chain"]