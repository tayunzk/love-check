FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY backend/go.mod backend/go.sum* ./
RUN GOPROXY=https://goproxy.cn,direct go mod tidy && GOPROXY=https://goproxy.cn,direct go mod download

COPY backend/ .
RUN GOPROXY=https://goproxy.cn,direct go get github.com/gin-gonic/gin@v1.9.1 && CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY ./frontend/ ./frontend/
RUN mkdir -p data

EXPOSE 8080
CMD ["./server"]