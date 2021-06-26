# Dockerfile References: https://docs.docker.com/engine/reference/builder/
# Alpine packages: https://pkgs.alpinelinux.org/packages
FROM alpine:3.14 AS build
LABEL maintainer="Grigorii Merkushev <brushknight@gmail.com>"
RUN apk update
RUN apk upgrade
RUN apk add --update go=1.16.5-r0 gcc=10.3.1_git20210424-r2 g++=10.3.1_git20210424-r2 make=4.3-r0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY Makefile ./Makefile
RUN make docker/compile

### Create executable image
FROM alpine:latest AS app
LABEL maintainer="Grigorii Merkushev <brushknight@gmail.com>"
WORKDIR /app
RUN mkdir /app/db
VOLUME /app/db
COPY --from=build /app/app .
RUN chmod +x ./app
COPY public ./public
EXPOSE 80
CMD ["./app"]
