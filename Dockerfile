# multistage dockerfile for lightweight production images

# first, build the go binary
FROM golang:1.10 as go

COPY . /go/src/github.com/briansimoni/stereodose

WORKDIR /go/src/github.com/briansimoni/stereodose

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stereodose .


# next, install node_modules and run a build for react
FROM node:10 as node

WORKDIR /stereodose/

COPY --from=go /go/src/github.com/briansimoni/stereodose/stereodose .
COPY --from=go /go/src/github.com/briansimoni/stereodose/app/views ./app/views/

WORKDIR /stereodose/app/views/
RUN npm install
RUN npm run build
RUN rm -rf node_modules


# Finally, take both artifacts and copy to a small, production ready image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /stereodose/
COPY --from=node /stereodose/ .
CMD ["./stereodose"]


