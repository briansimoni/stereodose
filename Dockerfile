# Inspired by https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324

# first, build the go binary
FROM 502859415194.dkr.ecr.us-east-1.amazonaws.com/golang:1.16.2-alpine as go

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

COPY . /go/src/github.com/briansimoni/stereodose

WORKDIR /go/src/github.com/briansimoni/stereodose

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o stereodose .


# next, install node_modules and run a build for react
FROM node:14-alpine as node

WORKDIR /stereodose/

COPY --from=go /go/src/github.com/briansimoni/stereodose/stereodose .
COPY --from=go /go/src/github.com/briansimoni/stereodose/app/views ./app/views/

WORKDIR /stereodose/app/views/
RUN npm install
RUN npm run build
RUN rm -rf node_modules

FROM scratch

# Import the user and group files from the builder.
COPY --from=go /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go /etc/passwd /etc/passwd
COPY --from=go /etc/group /etc/group
WORKDIR /stereodose/
# Copy our static executable.
COPY --from=node /stereodose/ .
# Use an unprivileged user.
USER appuser:appuser

CMD ["./stereodose"]




