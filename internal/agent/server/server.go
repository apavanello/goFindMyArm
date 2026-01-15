package server

import (
	"context"
	"fmt"
	"net"
	"os/exec"

	"github.com/apavanello/goFindMyArm/internal/shared/crypto"
	"github.com/apavanello/goFindMyArm/pkg/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type commandServer struct {
	protocol.UnimplementedCommandServiceServer
	passwordHash [crypto.KeySize]byte
}

// Serve starts the gRPC server
func Serve(port int, password string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	var pwdHash [crypto.KeySize]byte
	if password != "" {
		pwdHash = crypto.DeriveKey(password) // Pre-compute hash for comparison
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor(pwdHash, password)),
	)
	protocol.RegisterCommandServiceServer(s, &commandServer{
		passwordHash: pwdHash,
	})

	fmt.Printf("[gRPC] Listening on TCP %d\n", port)
	return s.Serve(lis)
}

// authInterceptor validates the token
func authInterceptor(validHash [crypto.KeySize]byte, rawPassword string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Public methods
		if info.FullMethod == "/protocol.CommandService/Ping" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata missing")
		}

		// Client sends the SHA256 of the password as token?
		// Or sends the password itself? Better send a hash or use TLS.
		// Since we defined NFR encryption everywhere, assuming gRPC over TLS or simple token match.
		// For MVP without TLS, sending the password/token is risky if finding it on wire.
		// But let's assume the client sends the PASSWORD string in "authorization-token".
		// We comparing the hash to avoid storing plain text in memory?

		tokens := md["authorization"]
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "token missing")
		}

		token := tokens[0]
		// Determine if token matches.
		// Naive check: token == rawPassword
		if token != rawPassword {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}

func (s *commandServer) Ping(ctx context.Context, req *protocol.PingRequest) (*protocol.PingResponse, error) {
	return &protocol.PingResponse{Message: "pong"}, nil
}

func (s *commandServer) ExecuteCommand(ctx context.Context, req *protocol.CommandRequest) (*protocol.CommandResponse, error) {
	switch req.Action {
	case protocol.CommandRequest_REBOOT:
		fmt.Println("[Command] Reboot requested")
		go func() {
			// Delay slightly to return response
			cmd := exec.Command("sudo", "reboot") // Assuming sudoer or root
			cmd.Run()
		}()
		return &protocol.CommandResponse{Success: true, Message: "Rebooting now..."}, nil

	case protocol.CommandRequest_UPDATE:
		// TODO: Implement update logic
		return &protocol.CommandResponse{Success: false, Message: "Update not implemented yet"}, nil

	default:
		return &protocol.CommandResponse{Success: false, Message: "Unknown action"}, nil
	}
}
