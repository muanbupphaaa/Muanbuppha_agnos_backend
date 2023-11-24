# Dockerfile
FROM golang:latest

WORKDIR /usr/src/app

COPY ./main.go ./

EXPOSE 8080

CMD ["go", "run", "main.go"]