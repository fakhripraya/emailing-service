.PHONY: protos

protos:
	 protoc -I protos/ --go-grpc_out=protos/email --go_out=protos/email protos/email.proto