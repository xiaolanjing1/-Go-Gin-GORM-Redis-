FROM golang:1.26

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./

RUN go mod download

COPY  . .

RUN go build -o main .

EXPOSE 9090

CMD ["./main"]