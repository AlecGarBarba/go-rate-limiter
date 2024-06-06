# Go Rate Limiter

This project is a rate limiter implemented in Go.

## Assumptions

This assumes you have Go 1.21 or higher as thats what I used while developing this rate limiter.

It also assumes some basic docker / redis / Go knowledge.

## Setup

1. Clone the repository:

```bash
git clone https://github.com/yourusername/go-rate-limiter.git
cd go-rate-limiter
```

2. Build the Docker image:

```bash
make docker-build
```

3. Create and start a redis container from redis image

```bash
docker pull redis
make start-redis-local
```

When running dev/prod environments, make sure to have the right configuration files.

## Running the Project

You can run the project locally or in a development environment.

### Running Locally

Start a local Redis server:

```bash
	make start-redis-local
```

2. Run the application:

```bash
	make run-local
```

## Other commands

To stop the local Redis server, use make stop-redis-local.

To debug the Redis server, use make debug-redis.

## Testing

You can send a GET request to the "/" route of your application using curl. To send the request 10 times, use

```bash
make test-curl COUNT=10.
```

You can replace 10 with the number of requests you want to send. To change the configuration of the rate limiter, you can use the config files // TODO, not yet implemented
