# Build stage

FROM golang:1.19-alpine AS builder
WORKDIR /app
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
COPY . .
RUN go build -o main main.go
# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/main .

COPY app.env .

COPY start.sh .

COPY db/migration ./migration



EXPOSE 8080
CMD [ "/app/main" ] 
# Thực thi file main ở thử mục app


ENTRYPOINT [ "/app/start.sh" ]