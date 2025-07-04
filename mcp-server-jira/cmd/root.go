package cmd

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	pb "github.com/cuenobi/mcp-platform/shared/proto/gen"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	jira "github.com/cuenobi/mcp-platform/mcp-server-jira/internal"
)

type server struct {
	pb.UnimplementedJiraServiceServer
}

func (s *server) SyncIssues(ctx context.Context, req *pb.SyncRequest) (*pb.SyncResponse, error) {
	log.Printf("Syncing project: %s", req.ProjectKey)
	time.Sleep(time.Second)
	return &pb.SyncResponse{Status: "synced"}, nil
}

func (s *server) CreateCard(ctx context.Context, req *pb.CreateCardRequest) (*pb.CreateCardResponse, error) {
	log.Printf("CreateCard called with prompt: %s", req.Prompt)

	issueIdea, err := jira.GenerateIssueIdea(req.Prompt)
	if err != nil {
		return nil, err
	}

	issueIdea.Title = strings.ReplaceAll(issueIdea.Title, "\n", " ")
	if len(issueIdea.Title) > 255 {
		issueIdea.Title = issueIdea.Title[:255]
	}
	if len(issueIdea.Description) > 1000 {
		issueIdea.Description = issueIdea.Description[:1000]
	}

	issueKey, err := jira.CreateIssue(req.ProjectKey, issueIdea.Title, issueIdea.Description)
	if err != nil {
		return nil, err
	}

	return &pb.CreateCardResponse{
		IssueKey: issueKey,
		Status:   "created",
	}, nil
}

func (s *server) Message(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Printf("Message called with prompt: %s", req.Prompt)

	response, err := jira.ReceivePrompt(req.Prompt)
	if err != nil {
		return nil, err
	}

	return &pb.MessageResponse{
		Message: response,
	}, nil
}

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
