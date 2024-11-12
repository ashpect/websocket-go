FROM golang:1.23-alpine AS builder

# Set environment variables
ENV GO111MODULE=on 
ENV CGO_ENABLED=0
# ENV GOOS=linux
# ENV GOARCH=amd64
# cause i am building on macos
# pls add you os and arch in above, if you want specific

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o websocket-server .

# multi-stage build
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# copy the binary from builder
COPY --from=builder /app/websocket-server .
EXPOSE 8080

CMD ["./websocket-server"]
