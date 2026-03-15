.PHONY: protos

protos:
	protoc -I protos/ \
		--go_out=paths=source_relative:protos/currency \
		--go-grpc_out=paths=source_relative:protos/currency \
		protos/currency.proto