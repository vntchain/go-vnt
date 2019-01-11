FROM alpine:3.7

RUN \
  apk add --update go git make gcc musl-dev linux-headers ca-certificates && \
  git clone --depth 1 https://github.com/vntchain/go-vnt && \
  (cd go-vnt && make gvnt) && \
  cp go-vnt/build/bin/gvnt /gvnt && \
  apk del go git make gcc musl-dev linux-headers && \
  rm -rf /go-vnt && rm -rf /var/cache/apk/*

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/gvnt"]
