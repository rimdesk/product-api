package types

import (
	"context"
	"encoding/json"
	"net/http"

	"connectrpc.com/connect"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	productv1 "github.com/rimdesk/product-api/gen/rimdesk/product/v1"
	"github.com/rimdesk/product-api/pkg/data/entities"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

const (
	// ContextKeyUser is used to store the authenticated user's claims in context.
	ContextKeyUser = "UserClaimsKey"
	// XTenantKey is the metadata key for the company Id header
	XTenantKey = "x-tenant-id"
)

// UserAuthClaims represents the JWT claims structure
type UserAuthClaims struct {
	Exp               int64          `json:"exp"`
	Iat               int64          `json:"iat"`
	Jti               string         `json:"jti"`
	Iss               string         `json:"iss"`
	Aud               []string       `json:"aud"`
	Id                string         `json:"sub"`
	Typ               string         `json:"typ"`
	Azp               string         `json:"azp"`
	Sid               string         `json:"sid"`
	Acr               string         `json:"acr"`
	AllowedOrigins    []string       `json:"allowed-origins"`
	RealmAccess       RealmAccess    `json:"realm_access"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
	Scope             string         `json:"scope"`
	EmailVerified     bool           `json:"email_verified"`
	Organization      []string       `json:"organization"`
	Name              string         `json:"name"`
	PreferredUsername string         `json:"preferred_username"`
	GivenName         string         `json:"given_name"`
	FamilyName        string         `json:"family_name"`
	Email             string         `json:"email"`
	jwt.RegisteredClaims
}

// RealmAccess defines roles at the realm level
type RealmAccess struct {
	Roles []string `json:"roles"`
}

// ResourceAccess defines roles at the resource level
type ResourceAccess struct {
	Account AccountRoles `json:"account"`
}

// AccountRoles defines roles within the "account" resource
type AccountRoles struct {
	Roles []string `json:"roles"`
}

func (u UserAuthClaims) String() string {
	jb, _ := json.Marshal(u)
	return string(jb)
}

type GlobalConfig interface {
	LoadEnv()
	DatabaseConfig() *gorm.Config
	GetServerAddress() string
	Http2() *http2.Server
}

type ProductService interface {
	ListProducts(ctx context.Context, request *connect.Request[productv1.ListProductsRequest]) ([]*productv1.Product, error)
	CreateProduct(ctx context.Context, request *connect.Request[productv1.CreateProductRequest]) (*productv1.Product, error)
	GetProduct(ctx context.Context, request *connect.Request[productv1.GetProductRequest]) (*productv1.Product, error)
	UpdateProduct(ctx context.Context, request *connect.Request[productv1.UpdateProductRequest]) (*productv1.Product, error)
	DeleteProduct(ctx context.Context, request *connect.Request[productv1.DeleteProductRequest]) error
}

type ProductRepository interface {
	FindAll(string) ([]*entities.Product, error)
	FindById(id string) (*entities.Product, error)
	FindByCompanyIdAndId(string, string) (*entities.Product, error)
	Create(*entities.Product) error
	Update(*entities.Product) error
	Delete(*entities.Product) error
}

type ContextHelper interface {
	GetTenant(context.Context) (string, error)
	GetUserClaims(context.Context) *UserAuthClaims
	GetAccessToken(request connect.AnyRequest) (string, error)
}

type ServiceAuthenticator interface {
	ExtractHeaderToken(connect.AnyRequest) (string, error)
	ExtractToken(ctx context.Context) (string, error)
	GetVerifier() *oidc.IDTokenVerifier
	ValidateTokenMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

type GRPCAuthMiddleware interface {
	CorsMiddleware(http.Handler) http.Handler
	LoggingUnaryInterceptor() connect.UnaryInterceptorFunc
	CheckTenantIdPresenceInHeader() connect.UnaryInterceptorFunc
	UnaryTokenInterceptor(ServiceAuthenticator) connect.UnaryInterceptorFunc
}
