# Default to Docker Hub. Both arms pass DOCKER_REGISTRY explicitly — the local
# `docker build` from Makefile.docker, and the remote Job through
# --opt build-arg:DOCKER_REGISTRY — so this default applies only to a bare
# `docker build`. It is docker.io rather than the authenticated nuke mirror
# because the Job mounts a docker.io push credential and no mirror login: with
# the mirror here, a remote build that passed no build-arg died 401 Unauthorized
# pulling golang.
ARG DOCKER_REGISTRY=docker.io
FROM ${DOCKER_REGISTRY}/golang:1.27.3 AS build
ARG BUILD_GIT_VERSION=dev
ARG BUILD_GIT_COMMIT=none
ARG BUILD_DATE=unknown
COPY . /workspace
WORKDIR /workspace
# vendor/ is not committed. The local arm creates it in `check-go-mod` before
# the build runs, but a remote Job clones a clean tree and reaches this line
# with no vendor/ present — `go build -mod=vendor` then fails. Create it here
# when it is absent, so both arms build from the same Dockerfile.
RUN [ -d vendor ] || go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -mod=vendor -ldflags "-s" -a -installsuffix cgo -o /main
CMD ["/bin/bash"]

FROM ${DOCKER_REGISTRY}/alpine:3.23 AS alpine
RUN apk --no-cache add ca-certificates

FROM alpine
ARG BUILD_GIT_VERSION=dev
ARG BUILD_GIT_COMMIT=none
ARG BUILD_DATE=unknown
LABEL org.opencontainers.image.version="${BUILD_GIT_VERSION}"
COPY --from=build /main /main
ENV ZONEINFO=/zoneinfo.zip
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
ENV BUILD_GIT_VERSION=${BUILD_GIT_VERSION}
ENV BUILD_GIT_COMMIT=${BUILD_GIT_COMMIT}
ENV BUILD_DATE=${BUILD_DATE}
ENTRYPOINT ["/main", "-v=2"]