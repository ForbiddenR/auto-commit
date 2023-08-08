FROM golang:1.20-alpine3.16


WORKDIR /app
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV TAG=nixx.xx.com/xxxx/test:0.0.5_230117_beta


COPY . .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o test /app/cmd/main.go

CMD ["./test"]
