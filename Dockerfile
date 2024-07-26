FROM node:latest AS googl-bye-frontend
WORKDIR /app/react
COPY react/package*.json ./
RUN npm install
COPY react/ ./
RUN npm run build

FROM golang:latest AS googl-bye-backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add \
    ca-certificates \
    postgresql-client \
    git
WORKDIR /root/
COPY --from=googl-bye-backend /app/main .
COPY --from=googl-bye-frontend /app/react/dist ./react/dist
# Wait for database to start
COPY wait-for-it.sh . 
RUN chmod +x wait-for-it.sh
EXPOSE 8080
CMD ["./wait-for-it.sh", "googl-bye-db", "./main"]