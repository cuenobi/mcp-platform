package cmd

import (
	"log"
	"net"

	pb "github.com/cuenobi/mcp-platform/shared/proto/gen"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var rootCmd = &cobra.Command{
	Use:   "mcp-server-jira",
	Short: "Run Jira gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Jira gRPC server...")

		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterJiraServiceServer(grpcServer, &server{})

		log.Println("Listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
