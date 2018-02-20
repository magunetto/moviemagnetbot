FROM golang:1.10

WORKDIR /go/src/app
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep && dep ensure -v
RUN go-wrapper install

CMD ["go-wrapper", "run"]
