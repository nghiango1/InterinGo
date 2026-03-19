# Use the official golang latest Debian base
FROM node:25-alpine3.22 AS website-build-env

WORKDIR /root/workspace

# Copy package.json and package-lock.json (if available)
COPY --parents website/package*.json .
WORKDIR /root/workspace/website
RUN npm install

# Copy the application source code
COPY --parents website/assets/ website/src/ website/*.ts website/*.js /root/workspace
RUN npm run build

# Use the official golang latest Debian base
FROM golang:1.26.1-alpine3.23 AS build-env

RUN apk add make

# Make main dir
WORKDIR /root/workspace
RUN mkdir bin
ENV PATH="/root/workspace/bin:$PATH"

# Download source dependancy and pull dependancy
COPY go.mod go.sum /root/workspace
RUN go mod download

# Download source
# Only move over nesseary file needed to build go-exec
COPY --parents ./pkg ./cmd Makefile /root/workspace
# Along with builded website dist
WORKDIR /
COPY --from=website-build-env --parents /root/workspace/website/dist .
WORKDIR /root/workspace

# Build step
RUN make embed-content
RUN make go-build

FROM alpine:latest
WORKDIR /app
COPY --from=build-env /root/workspace/dist/ /app
EXPOSE 8080
ENTRYPOINT /app/interingo -s
