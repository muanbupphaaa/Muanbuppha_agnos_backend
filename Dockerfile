
FROM golang:latest

WORKDIR /usr/src/app

RUN go get github.com/gorilla/mux

COPY ./main.go ./

EXPOSE 8080

RUN go build -o main .

CMD ["./main"]
