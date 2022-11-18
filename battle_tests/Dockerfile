# FROM golang:alpine AS builder
FROM golang:latest AS builder

ARG CURIO_VERSION
ARG BATTLE_TEST_SHA

RUN mkdir /app
ADD . /app/
WORKDIR /app

# RUN apk update && apk upgrade && apk add --no-cache build-base git

RUN go build -ldflags="-X 'github.com/bearer/curio/battle_tests/build.CurioVersion=${CURIO_VERSION}' -X 'github.com/bearer/curio/battle_tests/build.BattleTestSHA=${BATTLE_TEST_SHA}'" \
  -a -o dist/battle_tests ./battle_tests/battle_tests.go

# FROM alpine
FROM golang:latest

COPY --from=builder /app/dist/ /app/

# RUN apk update && apk upgrade && apk add --no-cache git libgcc libstdc++ musl-utils file
# required to run curio on alpine
# https://github.com/gliderlabs/docker-alpine/issues/219
# /lib64/ld-linux-x86-64.so.2 is provided by gcompat
# To get this info, I ran file and ldd on /app/curio
# https://askubuntu.com/questions/1388695/another-no-such-file-or-directory-but-the-file-exists/1388699#1388699
# RUN apk add --no-cache gcompat procps perf libc6-compat

WORKDIR /app

RUN mkdir battle-test-tmp
RUN chmod +x /app/curio

CMD ["/app/battle_tests"]