# Proto Package

This package contains Protocol Buffer definitions and generated Go code for the MCP Platform services.

## Directory Structure

```
shared/proto/
├── protos/           # Protocol Buffer definition files (.proto)
│   └── jira.proto   # Jira service definitions
├── gen/             # Generated Go code (DO NOT EDIT)
│   ├── jira.pb.go   # Generated message types
│   └── jira_grpc.pb.go # Generated gRPC service code
├── go.mod           # Go module definition
├── go.sum           # Go module dependencies
├── Makefile         # Build automation
└── README.md        # This file
```

## Services

### Jira Service

The Jira service provides gRPC endpoints for interacting with Jira:

- **SyncIssues**: Synchronizes issues from a Jira project
- **CreateCard**: Creates a new Jira card based on a natural language prompt

## Usage

### Importing in Go

```go
import pb "github.com/cuenobi/mcp-platform/shared/proto/gen"
```

### Generating Code

Use the Makefile to generate Go code from proto definitions:

```bash
# Generate all files
make generate

# Clean and regenerate
make regenerate

# Clean generated files
make clean

# Show help
make help
```

### Manual Generation

If you prefer to run protoc directly:

```bash
protoc --go_out=gen --go-grpc_out=gen protos/*.proto
```

## Development Guidelines

1. **Proto Files**: Edit only files in the `protos/` directory
2. **Generated Files**: Never edit files in the `gen/` directory - they are auto-generated
3. **Regeneration**: Always regenerate after modifying proto files
4. **Comments**: Add comprehensive comments to proto definitions
5. **Versioning**: Consider backward compatibility when making changes

## Dependencies

- `protoc` - Protocol Buffer compiler
- `protoc-gen-go` - Go plugin for protoc
- `protoc-gen-go-grpc` - Go gRPC plugin for protoc

## Installation

```bash
# Install Protocol Buffer compiler
brew install protobuf

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
``` 