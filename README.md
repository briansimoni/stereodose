# Stereodose

drug-inspired music app

[![Build Status](https://travis-ci.org/briansimoni/stereodose.svg?branch=gorm)](https://travis-ci.org/briansimoni/stereodose)

## Setup

For local development first create a .env file that looks like this:

```bash
STEREODOSE_CLIENT_ID=someclientid
STEREODOSE_CLIENT_SECRET=someclientsecret
STEREODOSE_REDIRECT_URL=http://localhost:4000/auth/callback
STEREODOSE_AUTH_KEY=somesecretkeythatyoucangenerate
STEREODOSE_ENCRYPTION_KEY=somesecretkeythatyoucangenerate
```

Obtain the client ID and secret by creating an app in the Spotify developer dashboard, and set the callback to http://localhost:3000/auth/callback
and
http://localhost:4000/auth/callback

It is helpful to have both

You can generate a random AUTH_KEY and ENCRYPTION_KEY by switching to the scripts directory

`go run secret_generator.go`

If you don't want to use a .env file, you can use regular system environment variables instead.

Once you have the variables set:
```
docker-compose up
# new terminal
cd app/views/
npm install
npm start
```
and you're all set to start writing code.

It's running a proxy server that comes bundled with React. This enables hot reloading among other nice things. It listens on port 3000 and proxies requests to the golang server on port 4000.

To more closely simulate production builds, run `npm run build` and visit localhost:4000
