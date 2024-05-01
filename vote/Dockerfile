FROM golang:1.20-alpine

COPY . /
RUN go mod download

EXPOSE 8080

ENTRYPOINT ["go", "run", "/main.go"]