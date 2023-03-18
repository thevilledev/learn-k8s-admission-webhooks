FROM golang:1.20-alpine AS builder

WORKDIR /app

RUN set -x && \
    apk add --no-cache make

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY pkg pkg/
COPY Makefile ./

RUN set -x && \
    make build

FROM scratch

COPY --from=builder /app/cmd/webhook/webhook /webhook

ENTRYPOINT [ "/webhook" ]
