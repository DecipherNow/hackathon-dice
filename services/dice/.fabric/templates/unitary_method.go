package {{.MethodsPackageName}}

import (
    "golang.org/x/net/context"

    "github.com/pkg/errors"
    
	{{.ProtobufImport}}
	{{.PBImport}}
)

{{.Comments}}
func (s *serverData) {{.MethodDeclaration}} {
	return nil, errors.New("not implemented")
}


