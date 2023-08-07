create:
	protoc --proto_path=uproto \
	--go_out=uproto \
	--go_opt=paths=source_relative \
	--go-grpc_out=uproto \
	--go-grpc_opt=paths=source_relative \
	uproto/*.proto