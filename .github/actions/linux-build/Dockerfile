FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y ca-certificates gcc-multilib g++-multilib software-properties-common wget && \
    \
    add-apt-repository ppa:git-core/ppa && \
    \
    install -m 0755 -d /etc/apt/keyrings && \
    (wget --output-document=- https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg) && \
    chmod a+r /etc/apt/keyrings/docker.gpg && \
    ( \
      echo \
        "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
        "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" \
        > /etc/apt/sources.list.d/docker.list \
    ) && \
    \
    apt-get update && \
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin git gcc-aarch64-linux-gnu g++-aarch64-linux-gnu && \
    \
    wget --output-document=/tmp/go.tar.gz https://go.dev/dl/go1.25.0.linux-amd64.tar.gz && \
    tar --extract --gunzip --file=/tmp/go.tar.gz --directory=/usr/local && \
    ln -s /usr/local/go/bin/go /usr/local/bin/ && \
    \
    wget --output-document=/tmp/goreleaser.deb https://github.com/goreleaser/goreleaser-pro/releases/download/v1.25.0-pro/goreleaser-pro_1.25.0.pro_amd64.deb && \
    dpkg --install /tmp/goreleaser.deb

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
