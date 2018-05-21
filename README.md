# Stereodose

drug-inspired music app

## Setup

For local development first create a .env file that looks like this:

```bash
STEREODOSE_CLIENT_ID=someclientid
STEREODOSE_CLIENT_SECRET=someclientsecret
STEREODOSE_REDIRECT_URL=http://localhost:3000/auth/callback
STEREODOSE_AUTH_KEY=somesecretkeythatyoucangenerate
```

Obtain the client ID and secret by creating an app in the Spotify developer dashboard, and set the callback to http://localhost:3000/auth/callback

You can generate a random AUTH_KEY by switching to the scripts directory
`go run secret_generator.go`

Once you have the .env file set:
`docker-compose up`
and you're all set to start writing code.