FROM golang:alpine AS builder

RUN mkdir /app
ADD . /app/
WORKDIR /app

ADD .gitconfig /root/

RUN apk update && apk upgrade && apk add --no-cache build-base git

RUN go build -a -o dist/battle_tests ./battle_tests/battle_tests.go

FROM alpine

COPY --from=builder /app/dist/ /app/

RUN apk update && apk upgrade && apk add --no-cache git libgcc libstdc++

WORKDIR /app

RUN mkdir battle-test-tmp

CMD ["/app/battle_tests"]