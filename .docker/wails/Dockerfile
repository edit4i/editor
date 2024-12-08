FROM ubuntu:22.04

# Prevent tzdata questions
ENV DEBIAN_FRONTEND=noninteractive

# Install system dependencies
RUN apt-get update && apt-get install -y \
    wget \
    git \
    build-essential \
    libgtk-3-0 \
    libwebkit2gtk-4.0-dev \
    gcc-aarch64-linux-gnu \
    musl \
    musl-dev \
    musl-tools \
    upx \
    curl \
    pkg-config \
    software-properties-common \
    linux-headers-generic \
    npm

# Install Go 1.22
RUN wget https://go.dev/dl/go1.22.7.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.7.linux-amd64.tar.gz && \
    rm go1.22.7.linux-amd64.tar.gz

# Add Go to PATH
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

# Set musl as the C compiler
ENV CC=/usr/bin/musl-gcc

# Install Node.js 22
RUN apt-get install -y curl && \
    apt-get remove -y nodejs nodejs-doc libnode-dev && \
    apt-get autoremove -y && \
    curl -fsSL https://deb.nodesource.com/setup_22.x -o nodesource_setup.sh && \
    bash nodesource_setup.sh && \
    apt-get update && \
    apt-get install -y nodejs && \
    rm nodesource_setup.sh

# Install Wails
RUN go install github.com/wailsapp/wails/v2/cmd/wails@v2.9.2

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of the application
COPY . .

# Build frontend dependencies
RUN cd frontend && npm install

# Install additional Linux headers and development files
RUN apt-get update && apt-get install -y \
    linux-libc-dev \
    linux-headers-$(uname -r) \
    musl-dev \
    build-essential \
    gcc-multilib \
    libc6-dev

# Expose any necessary ports
EXPOSE 3000
