create:
	protoc --proto_path=proto \
	--go_out=proto \
	--go_opt=paths=source_relative \
	--go-grpc_out=proto \
	--go-grpc_opt=paths=source_relative \
	proto/*.proto

start: 
	docker-compose up

stop:
	docker-compose down

restart:
	docker-compose down && docker-compose up

lint:
	golangci-lint run ./... --config=./.golangci.yml
