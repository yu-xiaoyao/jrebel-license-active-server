FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go build .

FROM golang:alpine
LABEL org.opencontainers.image.author="yu-xiaoyao" \
    org.opencontainers.image.description="JRebel and XRebel active server" \
    org.opencontainers.image.source="https://github.com/yu-xiaoyao/jrebel-license-active-server" \
    org.opencontainers.image.url="https://github.com/yu-xiaoyao/jrebel-license-active-server" \
    org.opencontainers.image.title="jrebel-license-active-server"
WORKDIR /app
COPY --from=builder /build/jrebel-license-active-server .
ENV TZ=Asia/Shanghai
EXPOSE 12345
CMD ["./jrebel-license-active-server"]