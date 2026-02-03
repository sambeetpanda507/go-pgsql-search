FROM ubuntu:latest

RUN apt-get update && apt-get install -y

# Install git 
RUN apt install git -y

# Install wget
RUN apt-get update && apt-get install wget -y

# Install golang
ENV GO_VERSION=1.25.0

RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# Persist Go in PATH
ENV PATH="/usr/local/go/bin:${PATH}"

# Clone the repository
RUN git clone https://github.com/sambeetpanda507/go-pgsql-search.git
