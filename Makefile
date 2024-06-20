.PHONY: mocks
mocks:
	sleep 1 && rm -rfd mocks && mockery

swag:
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/main.go

run:  swag
	go run cmd/main.go serve

build:  swag
	go build -o tmp/main cmd/main.go

build-win:  swag
	go build -o tmp/main.exe cmd/main.go
		
docker-run:
	docker run -d -p 8080:8080 \
	--name vocagame \
	vocagame:latest

air:
	air -c .air.toml serve

air-win:
	air -c .air.win.toml serve

test-run:
	export DB_HOST=localhost &&  export DB_PORT=5432 && export DB_USERNAME=postgres && export DB_PASSWORD=admin && export DB_NAME=digihope && go run cmd/main.go serve