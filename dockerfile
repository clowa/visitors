# This docker files does a multi stage docker build and creates
# a statically cross platform compiled binary shipped in a harded
# docker image from scratch

ARG GO_VERSION="1.20"
ARG ALPINE_VERSION="3.17"
ARG APP_VERSION

ARG TARGETOS
ARG TARGETARCH

# Golang cross compiling statical binary
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as build

WORKDIR /build

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ["src/*.go", "./"]

# ldflags -w disables debug, letting the file be smaller.
# netgo makes sure we use built-in net package and not the system’s one.
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS}  \
    GOARCH=${TARGETARCH}  \
    GOARM=$(${TARGETVARIANT} | tr -d 'v')  \
    go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o /visitors

# Final harded image from scratch
# Build from google distroless projekt image
# FROM gcr.io/distroless/static AS final 
FROM scratch
COPY --from=build /visitors /visitors

EXPOSE 8080
CMD [ "/visitors" ]