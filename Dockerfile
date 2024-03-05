FROM golang:1.22-bullseye

ENV PATH=$PATH:/

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o /bin/app ./cmd/main.go

EXPOSE 80 8080

CMD ["./bin/app"]
