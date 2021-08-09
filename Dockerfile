# Dockerfile References: https://docs.docker.com/engine/reference/builder/
# Alpine packages: https://pkgs.alpinelinux.org/packages
FROM alpine:3.14 AS environment
LABEL maintainer="Grigorii Merkushev <brushknight@gmail.com>"
ARG UI_VERSION_ARG=0.0.18
ENV UI_VERSION=$UI_VERSION_ARG
ARG TAG_ARG='dev'
ENV TAG=$TAG_ARG
# install dependencies
RUN apk update
RUN apk upgrade
RUN apk add --update go=1.16.7-r0 gcc=10.3.1_git20210424-r2 g++=10.3.1_git20210424-r2 make=4.3-r0 curl=7.78.0-r0
# download go modules
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# compile and download UI
FROM environment AS build
LABEL maintainer="Grigorii Merkushev <brushknight@gmail.com>"
# required for back compatibility with parent make
COPY envfile.template ./envfile
# compile
COPY cmd ./cmd
COPY internal ./internal
COPY Makefile ./Makefile
RUN make compile
# download UI
RUN mkdir -p ./public
RUN make download/ui

# copy compiled assets & config
FROM alpine:latest AS app
LABEL maintainer="Grigorii Merkushev <brushknight@gmail.com>"
ARG CONFIG_VERSION_ARG="./examples/config/web-sqlite.yml"
ENV CONFIG_VERSION=$CONFIG_VERSION_ARG
WORKDIR /app
# copy executable
COPY --from=build /app/app .
RUN chmod +x ./app
# copy config
COPY $CONFIG_VERSION /app/default-config.yml
# copy UI
RUN mkdir -p ./public
COPY --from=build /app/public/ ./public/
COPY ./static/index.html ./public/index.html
# prepare folder for sqlite db
RUN mkdir /app/db

EXPOSE 80
CMD ["/app/app"]

