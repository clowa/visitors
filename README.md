[![Build docker images](https://github.com/clowa/visitors/actions/workflows/docker-test.yml/badge.svg)](https://github.com/clowa/visitors/actions/workflows/docker-test.yml)

# Overview

Basic golang application displaying the total number of visitors of the website. A redis database is used

Supported platforms:

- `linux/386`
- `linux/s390x`
- `linux/amd64`
- `linux/arm/v6`
- `linux/arm/v7`
- `linux/arm64`
- `linux/mips64le`

# Getting started

To run the application locally you can simply do this via docker.

1. Clone this repository by running 'git clone ...'
2. Start the docker containers by running `docker-compose up --build`
3. Visit http://localhost

# Configuration

The following environment variables can be used to configure the application:
| Variable | Description | Example |
|-----------------------|--------------------------------------|--------------|
| `VISITORS_PORT` | Listening port of web app | `8080` |
| `VISITORS_REDIS_HOST` | IP or DNS of redis backend with port | `redis:6379` |
| `VISITORS_REDIS_DB` | Redis database ID to use | `0` |
