PROTO_PATH     = gen
PROTO_COMPILE ?= protoc

gen:
	mkdir -p $(PROTO_PATH) && \
		$(PROTO_COMPILE) --proto_path=../protocol/ ../protocol/hstream.proto \
			--go_out=$(PROTO_PATH) \
			--go-grpc_out=$(PROTO_PATH)
