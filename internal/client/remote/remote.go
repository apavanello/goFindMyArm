package remote

import (
	"context"
	"fmt"
	"time"

	"github.com/apavanello/goFindMyArm/pkg/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Reboot connects to the agent and requests a reboot
func Reboot(ip string, port int, password string) error {
	target := fmt.Sprintf("%s:%d", ip, port)

	// Dial
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := protocol.NewCommandServiceClient(conn)

	// Auth Metadata
	md := metadata.Pairs("authorization", password)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Execute
	req := &protocol.CommandRequest{
		Action: protocol.CommandRequest_REBOOT,
	}

	resp, err := client.ExecuteCommand(ctx, req)
	if err != nil {
		return fmt.Errorf("reboot command failed: %v", err)
	}

	if !resp.Success {
		return fmt.Errorf("reboot refused: %s", resp.Message)
	}

	return nil
}
