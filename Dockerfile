FROM golang:latest

WORKDIR ./src/github.com/xelaj/mcclane
COPY . .
RUN ls -la
RUN GO111MODULE=on go mod download

RUN make build

ENTRYPOINT ["./bin/mcclane"]