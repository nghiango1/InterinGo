# Use the official golang latest Debian base
FROM golang:1.26.1-alpine3.23 AS build-env

RUN apk add make

# Make main dir
WORKDIR /root/workspace

# Download source dependancy and pull dependancy
COPY go.mod go.sum /root/workspace
RUN go mod download

# Download source
# Only move over nesseary file needed to build go-exec
COPY --parents ./pkg ./cmd Makefile /root/workspace
# Along with builded website dist

# Build step
RUN make embed-content
RUN make go-build

# ??? This can be bad, glibc doesn't supported by alpine by default
FROM frolvlad/alpine-glibc:alpine-3.22_glibc-2.42 AS website-build-env

WORKDIR /root/workspace

RUN apk add --update nodejs npm
RUN npm install tree-sitter-cli

# Copy package.json and package-lock.json (if available)
# COPY --parents tree-sitter-interingo/package*.json .
# WORKDIR /root/workspace/tree-sitter-interingo
# RUN npm install

# Copy the application source code
WORKDIR /
COPY --parents tree-sitter-interingo/ /root/workspace
WORKDIR /root/workspace/tree-sitter-interingo
RUN npx tree-sitter generate

# FROM alpine:latest
# WORKDIR /app
# COPY --from=build-env /root/workspace/dist/ /app
# RUN apk add neovim
