package types

import (
	"context"

	"connectrpc.com/connect"
	productv1 "github.com/rimdesk/product-api/gen/rimdesk/product/v1"
	"github.com/rimdesk/product-api/pkg/data/entities"
	"golang.org/x/net/http2"
	"gorm.io/gorm"
)

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
	DeleteProduct(ctx context.Context, request *connect.Request[productv1.DeleteProductRequest])  error
}

type ProductRepository interface {
	FindAll(string) ([]*entities.Product, error)
	FindById(id string) (*entities.Product, error)
	FindByCompanyIdAndId(string, string) (*entities.Product, error)
	Create(*entities.Product) error
	Update(*entities.Product) error
	Delete(*entities.Product) error
}