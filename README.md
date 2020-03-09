# Stereodose

drug-inspired music app

[![Build Status](https://travis-ci.org/briansimoni/stereodose.svg?branch=master)](https://travis-ci.org/briansimoni/stereodose)

## Setup

For local development first create a .env file that looks like this:

```bash
STEREODOSE_CLIENT_ID=someclientid
STEREODOSE_CLIENT_SECRET=someclientsecret
STEREODOSE_REDIRECT_URL=http://localhost:4000/auth/callback
STEREODOSE_AUTH_KEY=somesecretkeythatyoucangenerate
STEREODOSE_ENCRYPTION_KEY=somesecretkeythatyoucangenerate
AWS_ACCESS_KEY_ID=someaccesskey
AWS_SECRET_ACCESS_KEY=somesecret
```

You can optionally specify STEREODOSE_DB_STRING and PORT variables

Obtain the client ID and secret by creating an app in the Spotify developer dashboard, and set the callback to http://localhost:3000/auth/callback
and
http://localhost:4000/auth/callback

It is helpful to have both

You can generate a random AUTH_KEY and ENCRYPTION_KEY by switching to the scripts directory

`go run secret_generator.go`

If you don't want to use a .env file, you can use regular system environment variables instead.


```
docker-compose up
# new terminal
cd app/views/
npm start
```
and you're all set to start writing code.

It's running a proxy server that comes bundled with React. This enables hot reloading among other nice things. It listens on port 3000 and proxies requests to the golang server on port 4000.

To more closely simulate production builds, run `npm run build` and visit localhost:4000

### Debugging
There is a vscode launch configuration for debugging inside the docker container. Simply open vscode, set breakpoints, and run the debugger. It works out of the box with dlv listening on port 40000. **However, if you change golang files or restart the server, you must also restart the debugging session or else you'll get nil pointer errors**


For React debugging, you can also use the "React" debug launch configuration. It requires the vscode chrome debugger extension. You can also read more and find references with the auto-generated documentation: https://github.com/briansimoni/stereodose/tree/master/app/views#debugging-in-the-editor

### Windows Users
In my experience, Docker isn't that great on Windows but sometimes a quick restart of the daemon/VM gets my containers to work. The golang file watcher does not work at all from inside a container. I recommend running the db from it's own container and then using native golang.
