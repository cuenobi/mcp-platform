package cmd

import (
	"context"
	"log"
	"time"

	    pb "github.com/cuenobi/mcp-platform/shared/proto/gen"

	jira "github.com/cuenobi/mcp-platform/mcp-server-jira/internal"
	llm "github.com/cuenobi/mcp-platform/mcp-server-jira/internal"
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

	idea, err := llm.GenerateIssueIdea(req.Prompt)
	if err != nil {
		return nil, err
	}

	issueKey, err := jira.CreateIssue(req.ProjectKey, idea.Title, idea.Description)
	if err != nil {
		return nil, err
	}

	return &pb.CreateCardResponse{
		IssueKey: issueKey,
		Status:   "created",
	}, nil
}
