package helpers

import (
	"context"
	"errors"
	"log"

	"connectrpc.com/connect"
	"github.com/rimdesk/product-api/pkg/types"
	"google.golang.org/grpc/metadata"
)

type grpcRequestHelper struct {
	authenticator types.ServiceAuthenticator
}

func (helper grpcRequestHelper) GetAccessToken(request connect.AnyRequest) (string, error) {
	return helper.authenticator.ExtractHeaderToken(request)
}

func (helper grpcRequestHelper) GetUserClaims(ctx context.Context) *types.UserAuthClaims {
	userClaims := ctx.Value(types.ContextKeyUser).(*types.UserAuthClaims)
	return userClaims
}

func (helper grpcRequestHelper) GetTenant(ctx context.Context) (string, error) {
	// Extract metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("could not extract metadata")
	}

	// Check if the X-Company-Id header is present
	companyID := md[types.XTenantKey]
	if len(companyID) == 0 {
		return "", errors.New("could not extract company id")
	}

	log.Printf("Received X-Company-Id: %s", companyID[0])

	return companyID[0], nil
}

func NewContextHelper(authenticator types.ServiceAuthenticator) types.ContextHelper {
	return &grpcRequestHelper{
		authenticator: authenticator,
	}
}
