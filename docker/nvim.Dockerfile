# Use the official golang latest Debian base
FROM golang:latest AS build-env

# Set up neovim for dev environment - So you can check out LSP/Highlight
RUN apt-get update
RUN apt-get install -y git cmake nodejs npm tar

RUN npm install -g tree-sitter-cli

# Make main dir
WORKDIR /root
RUN mkdir .config
RUN mkdir workspace

# Make workspace
WORKDIR /root/workspace
RUN mkdir bin
ENV PATH="/root/workspace/bin:$PATH"
ENV TERM=xterm-256color

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest
# No need, golang base already setup bin
# RUN export PATH="$PATH:~/go/bin"

# Install tailwind-cli
RUN wget https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
RUN chmod +x tailwindcss-linux-x64
RUN mv tailwindcss-linux-x64 /root/workspace/bin

# Install neovim
RUN curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux64.tar.gz
RUN tar -C /opt -xzf nvim-linux64.tar.gz
ENV PATH="$PATH:/opt/nvim-linux64/bin"

# Download source
RUN git clone https://github.com/nghiango1/InterinGo.git

# Download nvim config
RUN git clone https://github.com/nghiango1/nvim.git
RUN ln -s /root/workspace/nvim /root/.config/
WORKDIR /root/.config/nvim
RUN git switch minimal

# Build step
WORKDIR /root/workspace/InterinGo
RUN make build

WORKDIR /root/workspace/InterinGo/lsp-interingo
RUN go build .
RUN chmod +x interingo-lsp

WORKDIR /root/workspace/InterinGo/tree-sitter-interingo
RUN npm run build

RUN ln -s ~/workspace/InterinGo/interingo ~/workspace/bin
RUN ln -s ~/workspace/InterinGo/lsp-interingo/interingo-lsp ~/workspace/bin

RUN nvim --headless "+Lazy! sync" +qa
RUN ln -s ~/workspace/InterinGo/tree-sitter-interingo/queries ~/.local/share/nvim/lazy/nvim-treesitter/queries/interingo
WORKDIR /root/workspace/InterinGo

ENTRYPOINT nvim test/function-01.iig
