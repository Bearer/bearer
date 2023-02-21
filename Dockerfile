FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add --no-cache git

COPY bearer /usr/local/bin/

RUN addgroup -S rungroup && adduser -S runuser -G rungroup
USER runuser

RUN git config --global --add safe.directory '*'

ENTRYPOINT ["bearer"]