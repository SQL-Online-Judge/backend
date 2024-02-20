FROM golang:1.21-alpine as builder

WORKDIR /SQL-Online-Judge/backend
COPY . .
RUN go build -o bin/core -v ./cmd/core

FROM alpine:latest
WORKDIR /SQL-Online-Judge/backend
RUN apk --no-cache add ca-certificates
COPY --from=builder /SQL-Online-Judge/backend/bin/core /SQL-Online-Judge/backend/bin/core

CMD ["/SQL-Online-Judge/backend/bin/core"]
