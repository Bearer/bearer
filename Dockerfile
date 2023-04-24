FROM ubuntu:23.04

RUN apt-get update && \
    apt-get install -y git && \
    \
    git config --global --add safe.directory '*' && \
    \
    adduser --system --group runuser

COPY bearer /usr/local/bin/

USER runuser

ENTRYPOINT ["bearer"]
