FROM golang:1.22.6-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy
CMD ["tail", "-f", "/dev/null"]