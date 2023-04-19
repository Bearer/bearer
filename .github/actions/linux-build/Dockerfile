FROM ubuntu:14.04

RUN apt-get update && \
    apt-get install -y ca-certificates gcc-multilib g++-multilib software-properties-common wget && \
    \
    add-apt-repository ppa:git-core/ppa && \
    apt-get update && \
    apt-get install -y git && \
    \
    wget --output-document=/tmp/go.tar.gz https://go.dev/dl/go1.19.8.linux-amd64.tar.gz && \
    tar --extract --gunzip --file=/tmp/go.tar.gz --directory=/usr/local && \
    ln -s /usr/local/go/bin/go /usr/local/bin/ && \
    \
    wget --output-document=/tmp/goreleaser.deb https://github.com/goreleaser/goreleaser-pro/releases/download/v1.17.2-pro/goreleaser-pro_1.17.2.pro_amd64.deb && \
    dpkg --install /tmp/goreleaser.deb

ENTRYPOINT ["/usr/bin/goreleaser"]
