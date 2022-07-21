For this build I used guidance from here:
https://blog.logrocket.com/how-to-build-a-restful-api-with-docker-postgresql-and-go-chi/
Source code:
https://gitlab.com/idoko/bucketeer

Steps:

For API routing using lightweight router for building GO HTTP services:
https://github.com/go-chi/chi/
https://pkg.go.dev/github.com/go-chi/chi

go get -u github.com/go-chi/chi/v5

The render package helps manage HTTP request/response payloads.
https://github.com/go-chi/render
https://pkg.go.dev/github.com/go-chi/render


To interact with our PostgreSQL databse we need this driver:
https://github.com/lib/pq

go get github.com/lib/pq

For migratation use:
https://github.com/golang-migrate/migrate
To create migrate files:
migrate create -ext sql -dir db/migrations -seq create_items_table

to add new permanent environment variables in latest Ubuntu versions

sudo -H gedit /etc/environment

don't forget to logout and login again to enable the environment variables.

export POSTGRESQL_URL="postgres://<user>:<password>@localhost:5432/<db_name>?sslmode=disable"
export POSTGRESQL_URL="postgres://postgres:postgres@localhost:5432/postgres_db?sslmode=disable"
migrate -database ${POSTGRESQL_URL} -path db/migrations up

Original docker build from the guide in case something wrong with mine.

# Here we choice docker image to pull and provide alias to it: 
FROM golang:1.14.6-alpine3.12 as builder
# On this step we copy dependencies files to go/src/ folder as a requirment 
COPY go.mod go.sum /go/src/gitlab.com/idoko/bucketeer/
# Estaiblish work directory
WORKDIR /go/src/gitlab.com/idoko/bucketeer
# Execute go mod command
RUN go mod download
# Copy root directory 
COPY . /go/src/gitlab.com/idoko/bucketeer
# ??
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/bucketeer gitlab.com/idoko/bucketeer

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/gitlab.com/idoko/bucketeer/build/bucketeer /usr/bin/bucketeer
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/bucketeer"]