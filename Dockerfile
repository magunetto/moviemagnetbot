FROM golang:1.10

WORKDIR /go/src/github.com/magunetto/moviemagnetbot
COPY . .

RUN cd cmd/moviemagnetbot && go build

CMD ["./cmd/moviemagnetbot/moviemagnetbot"]
