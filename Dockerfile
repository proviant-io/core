# Dockerfile References: https://docs.docker.com/engine/reference/builder/
# Alpine packages: https://pkgs.alpinelinux.org/packages
FROM alpine:3.14 AS environment
LABEL maintainer="Grigorii Merkushev <brushknight@gmail.com>"
ARG UI_VERSION_ARG=pre-alpha.1
ENV UI_VERSION=$UI_VERSION_ARG
RUN apk update
RUN apk upgrade
RUN apk add --update go=1.16.5-r0 gcc=10.3.1_git20210424-r2 g++=10.3.1_git20210424-r2 make=4.3-r0 curl=7.77.0-r1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM environment AS build
COPY cmd ./cmd
COPY internal ./internal
COPY Makefile ./Makefile
RUN make docker/compile
# download UI
RUN mkdir -p ./public
RUN make download/ui

FROM alpine:latest AS api
LABEL maintainer="Grigorii Merkushev <brushknight@gmail.com>"
WORKDIR /app
RUN mkdir /app/db
VOLUME /app/config
COPY --from=build /app/app .
RUN chmod +x ./app
EXPOSE 80
CMD ["./app"]

### Create executable image
FROM api AS web
LABEL maintainer="Grigorii Merkushev <brushknight@gmail.com>"
WORKDIR /app
RUN mkdir -p ./public
COPY --from=build /app/public/ ./public/
COPY ./static/index.html ./public/index.html
CMD ["./app"]
