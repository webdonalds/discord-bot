### Builder
FROM golang:1.24-alpine AS builder
RUN apk update && apk add git && apk add ca-certificates && apk add tzdata

WORKDIR /usr/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o dist/bot ./cmd/bot


### Make executable image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/src/app/dist /

CMD [ "/bot" ]
