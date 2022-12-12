# syntax = docker/dockerfile-upstream:1.2.0-labs

ARG TOOLCHAIN

# cleaned up specs and compiled versions
FROM scratch AS generate
FROM ghcr.io/sentinelos/certificates:2022-10-11 AS certificates
FROM ghcr.io/sentinelos/fhs:0.0.1 AS fhs

# runs markdownlint
FROM docker.io/node:19.2.0-alpine3.16 AS lint-markdown
WORKDIR /src
RUN npm i -g markdownlint-cli@0.32.2
RUN npm i sentences-per-line@0.2.1
COPY .markdownlint.json .
RUN markdownlint --ignore "CHANGELOG.md" --ignore "**/node_modules/**" --ignore '**/hack/chglog/**' --rules node_modules/sentences-per-line/index.js .

# base toolchain image
FROM ${TOOLCHAIN} AS toolchain
RUN apk --update --no-cache add bash curl build-base protoc protobuf-dev

# build tools
FROM --platform=${BUILDPLATFORM} toolchain AS tools
ENV GO111MODULE on
ARG CGO_ENABLED
ENV CGO_ENABLED ${CGO_ENABLED}
ENV GOPATH /go
ARG GOLANGCILINT_VERSION
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCILINT_VERSION} \
	&& mv /go/bin/golangci-lint /bin/golangci-lint
ARG GOFUMPT_VERSION
RUN go install mvdan.cc/gofumpt@${GOFUMPT_VERSION} \
	&& mv /go/bin/gofumpt /bin/gofumpt
RUN go install golang.org/x/vuln/cmd/govulncheck@latest \
	&& mv /go/bin/govulncheck /bin/govulncheck
ARG GOIMPORTS_VERSION
RUN go install golang.org/x/tools/cmd/goimports@${GOIMPORTS_VERSION} \
	&& mv /go/bin/goimports /bin/goimports

# tools and sources
FROM tools AS base
WORKDIR /src
COPY ./go.mod .
COPY ./go.sum .
RUN --mount=type=cache,target=/go/src go mod download
RUN --mount=type=cache,target=/go/src go mod verify
COPY . .
RUN --mount=type=cache,target=/go/src go list -mod=readonly all >/dev/null

# builds ensurer-linux-amd64
FROM base AS ensurer-linux-amd64-build
COPY --from=generate / /
WORKDIR /src
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="github.com/sentinelos/ensurer/pkg/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/src GOARCH=amd64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=ensurer -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /ensurer-linux-amd64

# builds ensurer-linux-arm64
FROM base AS ensurer-linux-arm64-build
COPY --from=generate / /
WORKDIR /src
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="github.com/sentinelos/ensurer/pkg/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/src GOARCH=arm64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=ensurer -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /ensurer-linux-arm64

# runs gofumpt
FROM base AS lint-gofumpt
RUN FILES="$(gofumpt -l .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'gofumpt -w .':\n${FILES}"; exit 1)

# runs goimports
FROM base AS lint-goimports
RUN FILES="$(goimports -l -local github.com/sentinelos/ensurer .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'goimports -w -local github.com/sentinelos/ensurer .':\n${FILES}"; exit 1)

# runs golangci-lint
FROM base AS lint-golangci-lint
COPY .golangci.yml .
ENV GOGC 50
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/root/.cache/golangci-lint --mount=type=cache,target=/go/src golangci-lint run --config .golangci.yml

# runs govulncheck
FROM base AS lint-govulncheck
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/src govulncheck ./...

# runs unit-tests with race detector
FROM base AS unit-tests-race
ARG TESTPKGS
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/src --mount=type=cache,target=/tmp CGO_ENABLED=1 go test -v -race -count 1 ${TESTPKGS}

# runs unit-tests
FROM base AS unit-tests-run
ARG TESTPKGS
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/src --mount=type=cache,target=/tmp go test -v -covermode=atomic -coverprofile=coverage.txt -coverpkg=${TESTPKGS} -count 1 ${TESTPKGS}

FROM scratch AS ensurer-linux-amd64
COPY --from=ensurer-linux-amd64-build /ensurer-linux-amd64 /ensurer-linux-amd64

FROM scratch AS ensurer-linux-arm64
COPY --from=ensurer-linux-arm64-build /ensurer-linux-arm64 /ensurer-linux-arm64

FROM scratch AS unit-tests
COPY --from=unit-tests-run /src/coverage.txt /coverage.txt

FROM ensurer-${TARGETOS}-${TARGETARCH} AS ensurer

FROM scratch AS ensurer-all
COPY --from=ensurer-linux-amd64 / /
COPY --from=ensurer-linux-arm64 / /

FROM scratch AS ensurer-image
ARG TARGETOS
ARG TARGETARCH
COPY --from=fhs / /
COPY --from=certificates / /
COPY --from=ensurer ensurer-${TARGETOS}-${TARGETARCH} /usr/bin/ensurer

LABEL org.opencontainers.image.title="Ensurer"
LABEL org.opencontainers.image.description="Ensurer is a tool for enforcing policies on your pipelines. "
LABEL org.opencontainers.image.licenses="MPL-2.0"
LABEL org.opencontainers.image.authors="Sentinel OS Authors"
LABEL org.opencontainers.image.documentation="https://github.com/sentinelos/ensurer/blob/main/README.md"
LABEL org.opencontainers.image.source="https://github.com/sentinelos/ensurer"

ENTRYPOINT ["/usr/bin/ensurer"]