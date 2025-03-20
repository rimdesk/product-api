package middlewares

import (
	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"context"
	"github.com/rimdesk/product-api/pkg/exceptions"
	"github.com/rimdesk/product-api/pkg/routes"
	"github.com/rimdesk/product-api/pkg/types"
	"github.com/rs/cors"
	"log"
	"net/http"
	"slices"
	"time"
)

type grpcAuthMiddleware struct{}

func (middleware *grpcAuthMiddleware) CorsMiddleware(h http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})

	return corsMiddleware.Handler(h)
}

func (middleware *grpcAuthMiddleware) UnaryTokenInterceptor(authenticator types.ServiceAuthenticator) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// Extract the full method name from the request
			fullMethod := req.Spec().Procedure
			// If the method is public, skip authentication
			if slices.Contains(routes.PublicRoutes, fullMethod) {
				return next(ctx, req)
			}

			// Otherwise, apply the authentication middleware
			token, err := authenticator.ExtractHeaderToken(req)
			if err != nil {
				return nil, exceptions.ErrMissingOrInvalidToken
			}

			// Validate the token
			idToken, err := authenticator.GetVerifier().Verify(ctx, token)
			if err != nil {
				return nil, exceptions.ErrInvalidToken
			}

			// Add user info to context
			claims := new(types.UserAuthClaims)
			if err := idToken.Claims(claims); err != nil {
				return nil, connect.NewError(connect.CodeInternal, exceptions.ErrFailedParsingTokenClaims)
			}

			// Add the claims to the context
			newCtx := context.WithValue(ctx, types.ContextKeyUser, claims)

			// Proceed with the handler
			return next(newCtx, req)
		}
	}
}

func (middleware *grpcAuthMiddleware) CheckTenantIdPresenceInHeader() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			tenantId := req.Header().Get(types.XTenantKey)
			if tenantId == "" {
				return nil, connect.NewError(connect.CodeInvalidArgument, exceptions.ErrMissingTenantHeader)
			}

			return next(ctx, req)
		}
	}
}

func (middleware *grpcAuthMiddleware) LoggingUnaryInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			start := time.Now()
			procedure := req.Spec().Procedure

			log.Printf("gRPC Method: %s, Request: %+v", procedure, req)

			resp, err := next(ctx, req)

			duration := time.Since(start)

			if err != nil {
				log.Printf("gRPC Method: %s, Error: %v, Duration: %s", procedure, err, duration)
			} else {
				log.Printf("gRPC Method: %s, Response: %+v, Duration: %s", procedure, resp, duration)
			}

			return resp, err
		}
	}
}

func New() types.GRPCAuthMiddleware {
	return &grpcAuthMiddleware{}
}