# Builder stage to build the Go binary of the application.
FROM golang:1.21-alpine3.18 as builder

# Install make and then build the application using `make build`
RUN apk update && apk add --no-cache make=~4.4
WORKDIR /app
COPY . /app/
RUN make build

# Final stage based on a lightweight alpine image.
FROM alpine:3.18

# Copy the binary from the builder stage and set params to be executed.
COPY --from=builder /app/blink.bin /
# Expose server and prometheus metrics port.
EXPOSE 8080
CMD ["/blink.bin"]