FROM --platform=${BUILDPLATFORM} golang:1.24.4 AS builder

WORKDIR /src
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/client ./cmd/echopb-client
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/server ./cmd/echopb-server

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /out/client /bin/client
COPY --from=builder /out/server /bin/server

CMD ["/bin/server"]
