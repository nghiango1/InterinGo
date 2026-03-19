# Use the official golang latest Debian base
FROM golang:1.26.1-alpine3.23 AS build-env

RUN apk add make

# Make main dir
WORKDIR /root
RUN mkdir workspace

# Make workspace
WORKDIR /root/workspace
RUN mkdir bin
ENV PATH="/root/workspace/bin:$PATH"

# Download source dependancy and pull dependancy
COPY go.mod go.sum /root/workspace
RUN go mod download

# Download source
# Only move over nesseary file needed to build go-exec
COPY --parents ./pkg ./cmd Makefile /root/workspace

# Build step
RUN make go-build

FROM alpine:latest
WORKDIR /app
COPY --from=build-env /root/workspace/dist/ /app
ENTRYPOINT /app/interingo -s
