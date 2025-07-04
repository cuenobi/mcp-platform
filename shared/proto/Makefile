PROTO_DIR = protos
GEN_DIR = gen

PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

.PHONY: all
all: generate

.PHONY: generate
generate:
	@echo "Generating Go files from proto definitions..."
	@mkdir -p $(GEN_DIR)
	@protoc --go_out=$(GEN_DIR) --go-grpc_out=$(GEN_DIR) $(PROTO_FILES)
	@echo "✅ Generation complete!"

.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	@rm -rf $(GEN_DIR)/*
	@echo "✅ Clean complete!"

.PHONY: regenerate
regenerate: clean generate

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all/generate - Generate Go files from proto definitions"
	@echo "  clean       - Remove generated files"
	@echo "  regenerate  - Clean and regenerate files"
	@echo "  help        - Show this help message" 