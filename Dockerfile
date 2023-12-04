FROM ubuntu:23.04

RUN apt-get update && \
    apt-get install -y git && \
    useradd --system --user-group --create-home runuser

COPY bearer /usr/local/bin/

USER runuser

RUN git config --global --add safe.directory '*'

ENTRYPOINT ["bearer"]
