# Stage 1: Build ứng dụng Go
FROM golang:1.23 AS builder

WORKDIR /app

# Copy module files và tải dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ mã nguồn
COPY . .

# Build binary tương thích Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o darius .

# Stage 2: Tạo container nhẹ để chạy
FROM alpine:latest

WORKDIR /root/

# Cài đặt chứng chỉ SSL nếu cần
RUN apk --no-cache add ca-certificates

# Copy binary từ stage builder
COPY --from=builder /app/darius .
COPY config.yaml /root/config.yaml

# Cấp quyền thực thi cho file binary
RUN chmod +x /root/darius

# Expose các cổng nếu cần
EXPOSE 50051 8080

# Chạy ứng dụng
ENTRYPOINT ["./darius"]
