package {{.MethodsPackageName}}

import (
    "github.com/pkg/errors"

	{{.ProtobufImport}}
	{{.PBImport}}
)

{{.Comments}}
func (s *serverData) {{.MethodDeclaration}} {
	return errors.New("not implemented")
}

