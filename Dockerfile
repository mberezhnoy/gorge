########################################################
# Ubuntu with go (to match github actions) environment #
########################################################
FROM ubuntu:20.04 as gobuntu

# gcc for cgo
# git for vscode
# ca-certificates and wget to download go
RUN apt-get update && apt-get install -y --no-install-recommends \
    g++ \
    gcc \
    libc6-dev \
    make \
    pkg-config \
    wget \ 
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

ENV GOPATH=$HOME/go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

ENV GOLANG_VERSION 1.14.6

RUN wget -c https://dl.google.com/go/go${GOLANG_VERSION}.linux-amd64.tar.gz -O - | tar -xz -C /usr/local

RUN go version
##############################
# Base dev image             #
##############################

FROM gobuntu as development

ENV GO111MODULE=on

RUN apt-get update && apt-get install -y --no-install-recommends \
    # Install Proj - C library for coordinate system conversion and its requirements 
    libproj-dev \
    # Graphviz is needed for pprof
    graphviz 

# Symlink this, so it's available under same path both here and on Mac when installed via brew
RUN ln -s /usr/lib/x86_64-linux-gnu/libproj.a /usr/local/lib/libproj.a

WORKDIR /workspace

################################
# Test/lint production   image #
################################

FROM development as tester

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make test && \
    make lint && \
    make typescript

################################
# Builder for production image #
################################

FROM tester as builder

ARG VERSION=0.0.0

RUN make build

################################
# Production image             #
################################
FROM gcr.io/distroless/cc-debian10 as production

COPY --from=builder /go/bin/gorge-server /go/bin/gorge-cli /usr/local/bin/

EXPOSE 7080

ENTRYPOINT ["gorge-server"]