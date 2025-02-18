---
title: Installation
layout: layouts/doc.njk
---

# Installing Bearer CLI

Installing Bearer CLI can be done though multiple methods referenced below. To update Bearer CLI, follow the update instructions for your install method.

## Installation options

### Install Script

The most common way to install Bearer CLI is with the install script. It will auto-select the best build for your architecture. Defaults installation to ./bin and to the latest release version:

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh
```

The default installation script works well for most use cases, but if you need more control you can customize the options by passing additional parameters.

- `-b`: sets the installation directory (defaults to `./bin`)
- `-d`: enables debug logging
- `[tag]`: specifies a tag release (defaults to the latest release)

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh -s -- -b /usr/local/bin
```

Note; Bearer requires `git` (2.34+) to be available on your `$PATH`. For MacOS we only support 12+ (Monterey) and later. Where these dependencies cannot be met we recommend using Docker.

### Homebrew

Using Bearer CLI's official [Homebrew tap](https://github.com/Bearer/homebrew-tap):

```bash
brew install bearer/tap/bearer
```

### Debian/Ubuntu

```bash
sudo apt-get update && sudo apt-get install ca-certificates -y && sudo update-ca-certificates
sudo apt-get install apt-transport-https
echo "deb [trusted=yes] https://apt.fury.io/bearer/ /" | sudo tee -a /etc/apt/sources.list.d/fury.list
sudo apt-get update
sudo apt-get install bearer
```

### RHEL/CentOS

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

```bash
sudo yum -y update
sudo yum -y install bearer
```

### Docker

Bearer CLI is also available as a Docker image on [Docker Hub](https://hub.docker.com/r/bearer/bearer) and [ghcr.io](https://github.com/bearer/bearer/pkgs/container/bearer).

With docker installed, you can run the following command with the appropriate paths in place of the examples.

```text
docker run --rm -v /path/to/repo:/tmp/scan bearer/bearer:latest-amd64 scan /tmp/scan
```

Additionally, you can use docker compose. Add the following to your docker-compose.yml file and replace the volumes with the appropriate paths for your project:

```yml
version: "3"
services:
  bearer:
    platform: linux/amd64
    image: bearer/bearer:latest-amd64
    volumes:
      - /path/to/repo:/tmp/scan
```

Then, run the docker compose run command to run Bearer CLI with any specified flags:

```bash
docker compose run bearer scan /tmp/scan --debug
```

### Binary

Download the archive file for your operating system/architecture from the [release page](https://github.com/bearer/bearer/releases).

Unpack the archive, and put the binary somewhere in your `$PATH` (on UNIX-y systems, `/usr/local/bin` or the like). Make sure it has permission to execute. Note; Bearer requires `git` to be available on your `$PATH` also.

## Updating Bearer CLI

Updating Bearer CLI provides you with the latest fixes and features. You can see an overview of the latest release on the [releases page](https://github.com/Bearer/bearer/releases/latest/).

### Install Script

To update to the latest version with the install script, run the install command again to override the existing installation.

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/bearer/main/contrib/install.sh | sh
```

### Homebrew

Update brew and upgrade Bearer CLI

```bash
brew update && brew upgrade bearer/tap/bearer
```

### Debian/Ubuntu

```bash
sudo apt-get update
sudo apt-get install bearer
```

### RHEL/CentOS

```bash
sudo yum -y update bearer
```

### Docker

The Docker configurations above uses the latest release, but you may need to tell Docker to explicitly retrieve the latest image:

```bash
docker image pull bearer/bearer:latest-amd64
```

### Binary

To update Bearer CLI when using the binary, download the latest release and overwrite your existing installation location.
