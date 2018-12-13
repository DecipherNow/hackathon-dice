package methods

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"

	pb "github.com/lucasmoten/project-2502/services/dice/protobuf"
)

func (s *serverData) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	return nil, errors.New("not implemented")
}
