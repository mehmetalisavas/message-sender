FROM golang:1.23

WORKDIR /app

### Install golang-migrate tool for migrations
# RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Required for golang-migrate to work
# RUN go install github.com/go-sql-driver/mysql@latest


# Install nc (netcat) to check MySQL connection
# RUN apt-get update && apt-get install -y netcat-openbsd


COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN make build

CMD ["./bin/message-sender"]
