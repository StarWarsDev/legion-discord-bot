FROM node:12.9.1-alpine as legion-data

ARG LEGION_DATA_VERSION

RUN mkdir -p /src
WORKDIR /src
COPY . .
RUN apk add make git
RUN make LEGION_DATA_VERSION=$LEGION_DATA_VERSION generate-json

FROM golang:1.12-alpine AS build

RUN mkdir -p /src
WORKDIR /src
COPY . .
RUN apk add make git
RUN make build

FROM alpine
LABEL maintainer="Steve Good (thestarwarsdev@gmail.com)"

ARG LEGION_DATA_VERSION
ENV LEGION_DATA_VERSION=$LEGION_DATA_VERSION

RUN mkdir -p /app
WORKDIR /app
COPY --from=build /src/legion-discord-bot /app/.
COPY --from=legion-data /src/legion-data.json /app/.

# install and update ca-certificates so our app can connect to discord
RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates 2>/dev/null || true

RUN addgroup -g 1000 -S ldb \
    && adduser -u 1000 -S ldb -G ldb \
    && chown -R ldb:ldb /app \
    && chmod -R 777 /app
USER ldb

CMD ["/bin/sh", "-c", "./legion-discord-bot"]
