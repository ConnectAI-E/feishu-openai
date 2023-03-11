FROM golang as golang

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build
ADD /code /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o feishu_chatgpt

FROM alpine:latest

WORKDIR /dist

RUN apk add --no-cache bash
COPY --from=golang /build/config.example.yaml /dist/config.yaml
COPY --from=golang /build/feishu_chatgpt /dist
ADD entrypoint.sh /dist/entrypoint.sh

RUN chmod +x /dist/entrypoint.sh
EXPOSE 9000
CMD /dist/entrypoint.sh
