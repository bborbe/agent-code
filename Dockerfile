# Default to Docker Hub. The Makefile always passes DOCKER_REGISTRY explicitly,
# so this default only applies to a build that passes no build-arg — which is
# exactly the remote arm: hack/build-job.yaml forwards SOURCE_DATE_EPOCH and
# nothing else. Leaving the authenticated nuke mirror here made the cluster pull
# 401 Unauthorized, since the Job mounts the docker.io push credential and not a
# mirror login.
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