package exceptions

import (
	"connectrpc.com/connect"
	"errors"
	"fmt"
)

var ErrMissingTenantHeader = connect.NewError(connect.CodeInvalidArgument, errors.New("x-tenant-id is required in the header"))
var ErrFailedParsingTokenClaims = connect.NewError(connect.CodeInvalidArgument, errors.New("token claims could not be parsed"))
var ErrInvalidToken = connect.NewError(connect.CodeUnauthenticated, errors.New(fmt.Sprintf("invalid token")))
var ErrMissingOrInvalidToken = connect.NewError(connect.CodeUnauthenticated, errors.New(fmt.Sprintf("missing or invalid token")))