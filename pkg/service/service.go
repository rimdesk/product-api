package service

import (
	"context"
	//"errors"

	"connectrpc.com/connect"
	productv1 "github.com/rimdesk/product-api/gen/protos/rimdesk/product/v1"

	//"github.com/rimdesk/product-api/pkg/clients"
	//"github.com/rimdesk/product-api/pkg/data/domains"
	//"github.com/rimdesk/product-api/pkg/data/dtos"
	"github.com/rimdesk/product-api/pkg/data/entities"
	"github.com/rimdesk/product-api/pkg/data/repository"
)

type ProductService interface {
	ListProducts(ctx context.Context, request *connect.Request[productv1.ListProductsRequest]) (*productv1.ListProductsResponse, error)
	CreateProduct(ctx context.Context, request *connect.Request[productv1.CreateProductRequest]) (*productv1.CreateProductResponse, error)
	GetProduct(ctx context.Context, request *connect.Request[productv1.GetProductRequest]) (*productv1.GetProductResponse, error)
	UpdateProduct(ctx context.Context, request *connect.Request[productv1.UpdateProductRequest]) (*productv1.UpdateProductResponse, error)
	DeleteProduct(ctx context.Context, request *connect.Request[productv1.DeleteProductRequest]) (*productv1.DeleteProductResponse, error)
	SearchWarehouse(ctx context.Context, request *connect.Request[productv1.SearchWarehouseRequest]) (*productv1.SearchWarehouseResponse, error)
}

type productService struct {
	productRepository repository.ProductRepository
	//warehouseClient   clients.WarehouseClient
}


func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}
//TO DO
func (service *productService) ListProducts(ctx context.Context, request *connect.Request[productv1.ListProductsRequest]) (*productv1.ListProductsResponse, error) {
	products, err := service.productRepository.FindAll(request.Msg.String())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	productProtos := make([]*productv1.Product, len(products))
	for _, product := range products {
		productProtos = append(productProtos, product.ToProto())
	}

	return &productv1.ListProductsResponse{Products: productProtos}, nil
}


func (service *productService) GetProduct(ctx context.Context, request *connect.Request[productv1.GetProductRequest]) (*productv1.GetProductResponse, error) {
	if _,err := service.productRepository.FindById(request.Msg.GetId()); err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return &productv1.GetProductResponse{}, nil
}

func (service *productService) CreateProduct(ctx context.Context, request *connect.Request[productv1.CreateProductRequest]) (*productv1.CreateProductResponse, error) {
	inventory, err := entities.NewProductFromRequest(request.Msg.GetProduct())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := service.productRepository.Create(inventory); err != nil {
		return nil, err
	}

	return &productv1.CreateProductResponse{}, nil
}


func (service *productService) UpdateProduct(ctx context.Context, request *connect.Request[productv1.UpdateProductRequest]) (*productv1.UpdateProductResponse, error) {

	companyID := request.Msg.Id
	product, err := service.productRepository.FindById(companyID)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	if err := service.productRepository.Update(product); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &productv1.UpdateProductResponse{}, nil
}

func (service *productService) DeleteProduct(ctx context.Context, request *connect.Request[productv1.DeleteProductRequest]) (*productv1.DeleteProductResponse, error) {
	productID := request.Msg.GetId()

	product, err := service.productRepository.FindById(productID)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	if err := service.productRepository.Delete(product); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return nil, err
}

func (service *productService) SearchWarehouse(ctx context.Context, request *connect.Request[productv1.SearchWarehouseRequest]) (*productv1.SearchWarehouseResponse, error) {
	_, err := service.productRepository.FindById(request.Msg.CompanyId)	
	if err != nil {
		return nil, err
	}

	return nil, nil
}
