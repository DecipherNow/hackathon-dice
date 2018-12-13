package methods

import (
	"google.golang.org/grpc"

	pb "github.com/lucasmoten/project-2502/services/dice/protobuf"
	"github.com/rs/zerolog"
)

type serverData struct {
	zerolog.Logger
}

// NewDiceServer returns an object that implements the  interface
func CreateAndRegisterServer(
	logger zerolog.Logger,
	grpcServer *grpc.Server,
) {
	var server pb.DiceServer = &serverData{
		logger,
	}

	pb.RegisterDiceServer(grpcServer, server)

}
