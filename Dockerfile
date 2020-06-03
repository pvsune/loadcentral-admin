FROM golang:1.14

COPY go.mod go.sum /app/
WORKDIR /app

RUN go mod download && go mod verify

COPY . /app
CMD go run main.go  # TODO: Build, run the binary instead.
