FROM golang as golang

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build
ADD /code /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o feishu_chatgpt

FROM scratch

WORKDIR /dist

COPY --from=golang /build/config.yaml /dist/config.yaml
COPY --from=golang /build/feishu_chatgpt /dist

EXPOSE 9000
ENTRYPOINT ["/dist/feishu_chatgpt"]
