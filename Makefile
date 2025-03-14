PROTO_DIR=./internal/api/proto
OUT_DIR=generated

PROTOC_GEN_GO=$(shell which protoc-gen-go)

all: generate

generate:
	@if [ -z "$(PROTOC_GEN_GO)" ]; then \
		echo "protoc-gen-go не найден. Установите его с помощью 'go install google.golang.org/protobuf/cmd/protoc-gen-go'"; \
		exit 1; \
	fi
	mkdir -p $(OUT_DIR)
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       -I $(PROTO_DIR) $(PROTO_DIR)/*.proto

clean:
	rm -rf $(OUT_DIR)

.PHONY: all generate clean
