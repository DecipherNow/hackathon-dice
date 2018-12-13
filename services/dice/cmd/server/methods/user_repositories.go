package methods

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"

	pb "github.com/lucasmoten/project-2502/services/dice/protobuf"
)

func (s *serverData) UserRepositories(ctx context.Context, request *pb.UserRequest) (*pb.UserRepositoriesResponse, error) {
	return nil, errors.New("not implemented")
}
