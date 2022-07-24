# [Choice] Java version: 11, 15
ARG VARIANT=8
# Image is alpine based
FROM docker:20.10.11-dind-rootless
# Go into root
USER root

# Update package reg
RUN apk update
RUN apk upgrade

RUN apk add wget

# GO
RUN wget https://golang.org/dl/go1.18.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.18.4.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin

# C
RUN apk add \
    autoconf \
    gcc \
    g++ \
    gdb \
    git \
    libdrm-dev \
    libepoxy-dev \
    make \
    mesa-dev \
    strace

# TERMINAL
RUN apk add bash
# Git
RUN apk add git
RUN apk add py-pip
# For Docker in Docker
RUN apk add docker-compose

RUN apk add curl
RUN apk add npm

RUN apk add doxygen
RUN apk add graphviz

RUN npm config set unsafe-perm true

ENV PATH "$PATH:/usr/local/go/bin"

RUN go install github.com/micro/micro/v3@latest

ENV DOCKER_HOST=unix:///var/run/docker.sock
