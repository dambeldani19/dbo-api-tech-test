# Gunakan image Go resmi sebagai image dasar
FROM golang:1.22 AS builder

# Setel direktori kerja
WORKDIR /app

# Salin go.mod dan go.sum dan unduh dependensi
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode aplikasi
COPY . .

# Bangun aplikasi Go
RUN go build -o myapp .

# Gunakan image minimal untuk menjalankan aplikasi
# FROM debian:bullseye-slim
FROM ubuntu:22.04

# Salin binary dari tahap builder
COPY --from=builder /app/myapp /usr/local/bin/myapp

# Salin file .env
COPY --from=builder /app/.env /app/.env

# Jalankan aplikasi saat container dijalankan
ENTRYPOINT ["/usr/local/bin/myapp"]
