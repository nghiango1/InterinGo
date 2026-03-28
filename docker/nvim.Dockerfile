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

# Build step - which doesn't embed any document into the server

# Built the InterinGo program language
RUN make interingo-build

# Built the InterinGo - LSP for program language supported tool
RUN make lsp-build

FROM alpine:latest
WORKDIR /app/bin

RUN apk add neovim

COPY --from=build-env /root/workspace/dist/ /app/bin
ENV PATH=/app/bin:$PATH
COPY test/input/ /app

WORKDIR /root/.config/nvim/
COPY assets/init.lua .

WORKDIR /app
ENTRYPOINT ["nvim", "."]
