FROM golang:1.25-alpine

COPY ./ ./

RUN go mod download
RUN go build -o subs_service ./cmd

CMD ["./subs_service"]