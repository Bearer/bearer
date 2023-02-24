---
title: Installation
layout: layouts/doc.njk
---

# Installing Bearer

Installing Bearer can be done though multiple methods referenced below.

## Install Script

The most common way to install Bearer is with the install script. It will auto-select the best build for your architecture. Defaults installation to ./bin and to the latest release version:

```text
curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh
```

The default installation script works well for most use cases, but if you need more control you can customize the options by passing additional parameters.

- `-b`: sets the installation directory (defaults to `./bin`)
- `-d`: enables debug logging
- `[tag]`: specifies a tag release (defaults to the latest release)

```text
curl -sfL https://raw.githubusercontent.com/Bearer/curio/main/contrib/install.sh | sh -s -- -b /usr/local/bin
```

## Homebrew

Using Bearer's official [Homebrew tap](https://github.com/Bearer/homebrew-tap):

```text
brew install bearer/tap/bearer
```

## Debian/Ubuntu

```text
$ sudo apt-get install apt-transport-https
$ echo "deb [trusted=yes] https://apt.fury.io/bearer/ /" | sudo tee -a /etc/apt/sources.list.d/fury.list
$ sudo apt-get update
$ sudo apt-get install bearer
```

## RHEL/CentOS

Add repository setting:

```text
$ sudo vim /etc/yum.repos.d/fury.repo
[fury]
name=Gemfury Private Repo
baseurl=https://yum.fury.io/bearer/
enabled=1
gpgcheck=0
```

Then install with yum:

```text
$ sudo yum -y update
$ sudo yum -y install bearer
```

## Docker

Bearer is also available as a Docker image on [Docker Hub](https://hub.docker.com/r/bearer/bearer) and [ghcr.io](https://github.com/Bearer/bearer/pkgs/container/bearer).

With docker installed, you can run the following command with the appropriate paths in place of the examples.

```text
docker run --rm -v /path/to/repo:/tmp/scan bearer/bearer:latest-amd64 scan /tmp/scan
```

Additionally, you can use docker compose. Add the following to your docker-compose.yml file and replace the volumes with the appropriate paths for your project:

```text
version: "3"
services:
  bearer:
    platform: linux/amd64
    image: bearer/bearer:latest-amd64
    volumes:
      - /path/to/repo:/tmp/scan
```

Then, run the docker compose run command to run Bearer with any specified flags:

```text
docker compose run bearer scan /tmp/scan --debug
```

## Binary

Download the archive file for your operating system/architecture from here.

Unpack the archive, and put the binary somewhere in your $PATH (on UNIX-y systems, /usr/local/bin or the like). Make sure it has permission to execute.

