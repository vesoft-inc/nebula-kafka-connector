FROM alpine:3.18.2

RUN apk update \
  && apk upgrade \
  && apk add --no-cache \
    curl jq util-linux bash gcompat \
  && rm -rf /var/cache/apk/*

ADD bin/br-ent /usr/local/bin/br-ent