# This docker files does a multi stage docker build and creates
# a statically cross platform compiled binary shipped in a harded
# docker image from scratch

ARG GO_VERSION="1.18"
ARG ALPINE_VERSION="3.15"
ARG APP_VERSION

ARG TARGETOS
ARG TARGETARCH

# Golang cross compiling statical binary
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as build

ENV CGO_ENABLED=0
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

WORKDIR /build

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ["*.go", "./"]

# ldflags -w disables debug, letting the file be smaller.
# netgo makes sure we use built-in net package and not the system’s one.
RUN go build -tags netgo -ldflags '-w' -o /visits

# Final harded image from scratch
FROM scratch
COPY --from=build /visits /visits

EXPOSE 8080
CMD [ "/visits" ]