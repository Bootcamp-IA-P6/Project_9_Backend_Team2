# ─── ETAPA 1: COMPILACIÓN ───
FROM golang:1.26-alpine AS builder

# Instalar git por si alguna dependencia de Go lo requiere
RUN apk add --no-cache git

WORKDIR /app

# Copiar archivos de dependencias primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar el binario de forma estática eliminando dependencias de C (CGO_ENABLED=0)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# ─── ETAPA 2: EJECUCIÓN (PRODUCCIÓN) ───
FROM alpine:latest  

# Instalar certificados CA (esencial para que Go pueda hablar de forma segura con la API de YouTube y Supabase)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar únicamente el binario compilado desde la etapa anterior
COPY --from=builder /app/main .

# Exponer el puerto por defecto
EXPOSE 8080

# Arrancar la aplicación
CMD ["./main"]