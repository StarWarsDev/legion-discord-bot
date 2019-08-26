FROM golang:1.12-alpine AS build

RUN mkdir -p /src
WORKDIR /src
COPY . .
RUN apk add make git
RUN make build

FROM alpine

LABEL maintainer="Steve Good (thestarwarsdev@gmail.com)"

ENV DISCORD_TOKEN USE_THE_REAL_TOKEN

RUN mkdir -p /app
WORKDIR /app
COPY --from=build /src/legion-discord-bot /app/.
COPY --from=build /src/legion-data.json /app/.

# install and update ca-certificates so our app can connect to discord
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

RUN useradd ldb
USER ldb

CMD ["/bin/sh", "-c", "./legion-discord-bot"]
