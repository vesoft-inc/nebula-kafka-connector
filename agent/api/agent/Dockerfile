FROM reg.vesoft-inc.com/ci/golang:1.17-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build/zero

ADD go.mod .
ADD go.sum .
COPY pkg pkg
COPY pkgbuz pkgbuz
COPY api api
RUN go mod download

COPY ./api/agent/etc /app/etc
RUN go build -ldflags="-s -w" -o /app/agent ./api/agent/agent.go

FROM reg.vesoft-inc.com/ci/alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/agent /app/agent
COPY --from=builder /app/etc /app/etc

EXPOSE 8888

CMD ["./agent", "-f", "etc/agent-api.yaml"]
