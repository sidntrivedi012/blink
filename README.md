# blink
A URL shortener service that produces URLs shorter than the blink of an eye.

## Features

- Shortens the URLs.
- Supports redirects for shortened URLs.
- Shows metrics for top 3 shortened domains.
- Uses Redis for storing the data regarding URLs.
- Reports same URL for multiple shortening sessions of same URL instead of generating new one.

## Setup

### Manual

You need to have Go (>1.13 at least) and Redis installed.

1. Clone this repository.
2. Run `make build` to build the binary named as `blink.bin`.
3. In another terminal session, start a local redis server through the `redis-server` command. It should be running by default on port `6379`.
4. Set environment variables:
    - `REDIS_ADDRESS=localhost:6379`
    - `APP_PORT=8080`
4. Run `make run` to execute the binary.

### Docker

1. Download the `docker-compose.yml` file.
2. To execute Blink, run `docker-compose up -d` command and then the app should be ready on port "8080".
3. The `docker-compose.yml` uses my docker image pushed on Docker Hub.

## API Routes

- POST `/api/shorten` : Returns the shortened URL for the URL sent in the request body in the below form.
```json
{"long_url":"https://google.com/helloworld"}
```
- GET `/api/metrics` : Returns the top 3 most shortened domains along with the count of being shortened.
