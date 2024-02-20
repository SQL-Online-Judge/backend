FROM golang:1.21-alpine as builder

WORKDIR /SQL-Online-Judge
COPY . .
RUN go build -o core .

FROM alpine:latest
WORKDIR /SQL-Online-Judge
RUN apk --no-cache add ca-certificates
COPY --from=builder /SQL-Online-Judge/core /SQL-Online-Judge/core

CMD ["/SQL-Online-Judge/core"]
