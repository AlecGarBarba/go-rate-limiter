docker-build:
	docker build -t go-rate-limiter .

run-local: 
	ENV=local API_URL=http://localhost:3000 go run main.go


start-redis-local:
	docker run -d -p 6379:6379 --name redis redis

stop-redis-local:
	docker stop redis

debug-redis:
	docker exec -it redis redis-cli

run-dev:  # TBD depending the environment. 
    API_URL=http://dev.api.com go run main.go

run-prod: 
    API_URL=http://prod.api.com go run main.go

test-curl:
	for i in $$(seq 1 $(COUNT)); do curl http://localhost:8080/; done

test-curl-info:
	for i in $$(seq 1 $(COUNT)); do curl -i http://localhost:8080/; done