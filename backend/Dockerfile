FROM golang:1.20-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
ARG NGINX_PORT
EXPOSE $NGINX_PORT
CMD ["go", "run", "cmd/main.go"]

