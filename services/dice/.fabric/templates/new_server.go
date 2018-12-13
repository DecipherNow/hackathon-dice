package {{.MethodsPackageName}}

import (
	"google.golang.org/grpc"
	
	"github.com/rs/zerolog"
	pb "{{.PBImport}}"
)

type serverData struct{
	zerolog.Logger
}

// New{{.GoServiceName}}Server returns an object that implements the  interface
func CreateAndRegisterServer(
	logger zerolog.Logger,
	grpcServer *grpc.Server,
) {
	var server pb.{{.ServerInterfaceName}}Server = 
		&serverData{
			logger,
		}

	pb.Register{{.ServerInterfaceName}}Server(grpcServer, server)

}