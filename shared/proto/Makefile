# Makefile for proto generation

# Directories
PROTO_DIR = protos
GEN_DIR = gen

# Proto files
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

# Default target
.PHONY: all
all: generate

# Generate Go files from proto definitions
.PHONY: generate
generate:
	@echo "Generating Go files from proto definitions..."
	@mkdir -p $(GEN_DIR)
	@protoc --go_out=$(GEN_DIR) --go-grpc_out=$(GEN_DIR) $(PROTO_FILES)
	@echo "✅ Generation complete!"

# Clean generated files
.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	@rm -rf $(GEN_DIR)/*
	@echo "✅ Clean complete!"

# Regenerate (clean + generate)
.PHONY: regenerate
regenerate: clean generate

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all/generate - Generate Go files from proto definitions"
	@echo "  clean       - Remove generated files"
	@echo "  regenerate  - Clean and regenerate files"
	@echo "  help        - Show this help message" 