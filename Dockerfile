FROM golang:1.22.3-alpine AS builder


ENV GOPATH=/

COPY ./ ./
RUN go mod download && go mod verify
RUN go build -o stock-service ./cmd/main.go

CMD [ "./stock-service" ]
