FROM golang:alpine AS build_base

RUN apk add bash curl git make

ARG MIGRATE_VERSION=4.7.1

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | tar xz &&  \
	mv migrate.linux-amd64 /usr/local/bin/migrate && \
	migrate -version

RUN mkdir -p /src/app
WORKDIR /src/app
VOLUME /src/app

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod tidy

# Development stage
FROM build_base AS dev

COPY --from=build_base /usr/local/bin/migrate /usr/local/bin/migrate

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download
COPY entrypoint.sh /usr/local/bin/entrypoint

RUN chmod +x /usr/local/bin/entrypoint

ENTRYPOINT ["/usr/local/bin/entrypoint"]

# This image builds the server for production usage
FROM build_base AS builder

# copy source files and build the binary
COPY . .
# And compile the project
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server cmd/main.go
# Expose port 8000 to the outside world
EXPOSE 8000

#In this last stage, we start from a fresh Alpine image, to reduce the image size and not ship the Go compiler in our production artifacts.
FROM alpine:3.11 AS production

ENV APP_ENV=prod
RUN mkdir -p /src/app

WORKDIR /src/app

# Finally we copy the statically compiled Go binary.
COPY --from=builder /bin/server /bin/app
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /src/app/migrations migrations/

ENTRYPOINT ["/bin/app"]
