FROM golang:latest

WORKDIR ./src/github.com/xelaj/mcclane
COPY . .

RUN GO111MODULE=on go mod download

RUN go build -o ./bin/mcclane ./cmd/mcclane

ENTRYPOINT ["./bin/mcclane"]