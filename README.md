![Build docker images](https://github.com/clowa/visitors/actions/workflows/docker-buildx.yml/badge.svg)

# Overview

Basic golang application displaying the total number of visitors of the website. A redis database is used

Supported platforms:

- `linux/riscv64`
- `linux/ppc64le`
- `linux/s390x`
- `linux/386`
- `linux/mips64le`
- `linux/mips64`
- `linux/amd64`
- `linux/amd64/v2`
- `linux/arm/v7`
- `linux/arm/v6`
- `linux/arm64`

# Getting started

To run the application locally you can simply do this via docker.

1. Clone this repository by running 'git clone ...'
2. Start the docker containers by running `docker-compose up --build`
3. Visit http://localhost
