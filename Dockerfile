# Gunakan image Golang resmi sebagai base
FROM golang:alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod dan sum
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi menjadi binary bernama 'main'
RUN go build -o main .

# Gunakan image yang lebih ringan untuk menjalankan aplikasi
FROM alpine:latest
WORKDIR /root/

# Copy binary dari builder
COPY --from=builder /app/main .
# Copy folder templates karena dibutuhkan untuk tampilan HTML
COPY --from=builder /app/templates ./templates

# Jalankan aplikasi
CMD ["./main"]