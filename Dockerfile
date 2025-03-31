FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -a -o main ./cmd/api

FROM scratch
COPY --from=builder /app/main /main
EXPOSE 8080
CMD ["/main"]
