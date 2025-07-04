package jira

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/cuenobi/mcp-platform/shared/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	Sync(project string) error
	CreateCard(project, prompt string) (string, error)
	Message(prompt string) (string, error)
}

type grpcClient struct {
	conn   *grpc.ClientConn
	client pb.JiraServiceClient
}

func NewGRPCClient(addr string) Client {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewJiraServiceClient(conn)
	return &grpcClient{conn: conn, client: c}
}

func (g *grpcClient) Sync(project string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := g.client.SyncIssues(ctx, &pb.SyncRequest{
		ProjectKey: project,
	})
	return err
}

func (g *grpcClient) CreateCard(project, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if project == "" {
		project = os.Getenv("JIRA_PROJECT_KEY")
	}

	resp, err := g.client.CreateCard(ctx, &pb.CreateCardRequest{
		ProjectKey: project,
		Prompt:     prompt,
	})
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("received nil response from server")
	}
	return resp.IssueKey, nil
}

func (g *grpcClient) Message(prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	resp, err := g.client.Message(ctx, &pb.MessageRequest{
		Prompt: prompt,
	})
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("received nil response from server")
	}
	return resp.Message, nil
}
