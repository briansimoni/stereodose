FROM golang:1.10

ADD . /go/src/github.com/briansimoni/stereodose

WORKDIR /go/src/github.com/briansimoni/stereodose

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure

RUN go install

# Only for dev purposes
RUN go get github.com/codegangsta/gin

CMD stereodose