

FROM golang:1.26-alpine AS builder
 
WORKDIR /app
 

COPY go.mod go.sum ./
RUN go mod download
 

COPY main.go init.go utils.go benchmark.go ./
 

RUN CGO_ENABLED=0 GOOS=linux go build -o app .
 

FROM alpine:latest
 
WORKDIR /app
 
COPY --from=builder /app/app .
 

COPY ../../../../mnt/large_storage/2026_nndescent/ ../../../../mnt/large_storage/2026_nndescent/
 
CMD ["./app"]
 
