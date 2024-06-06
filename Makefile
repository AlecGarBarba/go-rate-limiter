docker-build:
	docker build -t go-rate-limiter .

run-local: 
	ENV=local API_URL=http://localhost:3000 go run main.go

run-dev:  # TBD depending the environment. 
    API_URL=http://dev.api.com go run main.go

run-prod: 
    API_URL=http://prod.api.com go run main.go