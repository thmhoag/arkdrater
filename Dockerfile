FROM golang:1.15.5-alpine as builder

WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags='-w -s -extldflags "-static"' -o ./bin/arkdrater ./cmd/arkdrater


FROM scratch as final

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/arkdrater /
COPY default.config.yaml /config/config.yaml

ENTRYPOINT [ "/arkdrater" ]
