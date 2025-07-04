module github.com/cuenobi/mcp-platform/mcp-server-jira

go 1.23.9

require (
	github.com/cuenobi/mcp-platform/shared/proto v0.0.0-20250704071713-180a9cb99ca0
	github.com/cuenobi/mcp-platform/shared/proto/gen v0.0.0
	github.com/spf13/cobra v1.9.1
	google.golang.org/grpc v1.73.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/cuenobi/mcp-platform/shared/proto/gen => ../shared/proto/gen
