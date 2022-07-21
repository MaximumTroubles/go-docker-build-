FROM golang:1.18.4-alpine3.16 as builder
COPY go.mod go.sum /go/src/github.com/maximumtroubles/go-docker-build/
WORKDIR /go/src/github.com/maximumtroubles/go-docker-build
RUN go mod download
COPY . /go/src/github.com/maximumtroubles/go-docker-build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/bucketeer github.com/MaximumTroubles/go-docker-build-

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/maximumtroubles/go-docker-build/build/bucketeer /usr/bin/bucketeer
EXPOSE 8080 8080
ENTRYPOINT [ "/usr/bin/bucketeer" ]