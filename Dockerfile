FROM golang:1.24.4 AS builder

WORKDIR /src
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/client ./cmd/client

FROM scratch

COPY --from=builder /out/server /bin/server
COPY --from=builder /out/client /bin/client

CMD ["/bin/server"]

