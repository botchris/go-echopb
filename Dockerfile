FROM golang:1.24.4 AS builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/client ./cmd/client

FROM scratch

COPY --from=builder /out/server /bin/server
COPY --from=builder /out/client /bin/client

ENTRYPOINT ["/bin/server"]
