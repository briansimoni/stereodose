FROM golang:1.10

COPY . /go/src/github.com/briansimoni/stereodose

WORKDIR /go/src/github.com/briansimoni/stereodose

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure

# Only for dev purposes
#RUN go get github.com/codegangsta/gin
RUN go get github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

CMD ./stereodose